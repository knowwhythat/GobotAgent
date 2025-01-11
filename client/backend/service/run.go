package service

import (
	"browser-agent/backend/service/sys_exec"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var agentCommand *exec.Cmd
var agentServerPort = 8888
var agentConn *websocket.Conn

func ExecuteUserTask(ctx context.Context, userInput string) error {
	apiBase := viper.GetString(LLMApiBase)
	modelName := viper.GetString(LLMModelName)
	apiKey := viper.GetString(LLMApiKey)
	configBrowserPath := viper.GetString(ConfigBrowserPath)
	if len(apiBase) == 0 || len(modelName) == 0 || len(apiKey) == 0 {
		return errors.New("请先配置模型参数")
	}
	err := checkAgentConn(ctx)
	if err != nil {
		return err
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]interface{})
	sendMessage["message_id"] = messageId
	data := make(map[string]interface{})
	sendMessage["method"] = "execute_task"

	data["task"] = userInput
	data["apiBase"] = apiBase
	data["modelName"] = modelName
	data["apiKey"] = apiKey
	data["configBrowserPath"] = configBrowserPath
	sendMessage["data"] = data
	request, err := json.Marshal(sendMessage)
	Logger.Logger.Info().Msg(string(request))
	if err != nil {
		return err
	}
	if err := agentConn.WriteMessage(1, request); err != nil {
		_ = agentConn.Close()
		agentConn = nil
		err := checkAgentConn(ctx)
		if err != nil {
			return err
		}
		if err := agentConn.WriteMessage(1, request); err != nil {
			_ = agentConn.Close()
			agentConn = nil
			return err
		}
		return nil
	}
	return nil
}

func startAgentServerCommand() {
	for ; agentServerPort < 9999; agentServerPort++ {
		if PortCheck(agentServerPort) {
			break
		}
	}

	agentCommand = sys_exec.BuildCmd(viper.GetString(ConfigPythonPath), "-u", "-B", "-m",
		"gobot_agent.main", strconv.Itoa(agentServerPort),
		filepath.Join(viper.GetString(ConfigDataPath), "log", "gobot_agent_log"+time.Now().Format("2006-01-02")+".txt"))
	var stderr bytes.Buffer
	agentCommand.Stderr = &stderr
	err := agentCommand.Run()
	if err != nil {
		agentCommand = nil
		errStr := stderr.String()
		if len(errStr) == 0 {
		}
		Logger.Logger.Error().Msg(errStr)
		Logger.Logger.Error().Err(err)
	}
}

func checkAgentConn(ctx context.Context) error {
	if agentCommand == nil {
		go startAgentServerCommand()
		time.Sleep(time.Second)
	}
	for i := 0; i < 30; i++ {
		if agentConn == nil {
			Logger.Logger.Info().Msgf("agentServerPort: %d", agentServerPort)
			conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%d/ws", agentServerPort), nil)
			if err == nil {
				agentConn = conn
				go listenInspectResponse(ctx)
				return nil
			} else {
				time.Sleep(time.Second)
				continue
			}
		} else {
			return nil
		}
	}
	return errors.New("连接BrowserAgent失败")
}

func listenInspectResponse(ctx context.Context) {
	for agentConn != nil {
		_, message, err := agentConn.ReadMessage()
		if err != nil {
			Logger.Logger.Error().Err(err)
			return
		}
		Logger.Logger.Info().Msgf("BrowserAgent响应: %s", string(message))
		runtime.EventsEmit(ctx, AgentResponseEvent, string(message))
	}
}
