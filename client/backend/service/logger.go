package service

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WailsLog struct {
	zerolog.Logger
}

var Logger WailsLog

func Init(name string) {
	logLevel := viper.GetString(ConfigLogLevel)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		zerolog.SetGlobalLevel(level)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	timeFormat := "2006-01-02 15:04:05"
	zerolog.TimeFieldFormat = timeFormat

	logDir := filepath.Join(viper.GetString(ConfigDataPath), "log")
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("Mkdir failed, err:", err)
		return
	}

	fileName := filepath.Join(logDir, name+"-"+time.Now().Format("2006-01-02")+".log")
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
	consoleWriter.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s;", i)
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	Logger = WailsLog{
		zerolog.New(multi).With().Timestamp().Logger().With().Caller().Logger(),
	}
}

func (WailsLog) Print(message string) {
	Logger.Logger.Info().Msg(message)
}
func (WailsLog) Trace(message string) {
	Logger.Logger.Trace().Msg(message)
}
func (WailsLog) Debug(message string) {
	Logger.Logger.Debug().Msg(message)
}
func (WailsLog) Info(message string) {
	Logger.Logger.Info().Msg(message)
}
func (WailsLog) Warning(message string) {
	Logger.Logger.Warn().Msg(message)
}
func (WailsLog) Error(message string) {
	Logger.Logger.Error().Msg(message)
}
func (WailsLog) Fatal(message string) {
	Logger.Logger.Fatal().Msg(message)
}
