package service

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

// WatchConfig 监控配置文件变化并热加载程序
func (c *Config) WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}

func (c *Config) Init() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath(".") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")         // 设置配置文件格式为YAML
	viper.AutomaticEnv()                // 读取匹配的环境变量
	viper.SetEnvPrefix("BROWSER_AGENT") // 读取环境变量的前缀为BROWSER_AGENT
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件

	}
	viper.SetDefault(ConfigLogLevel, "Info")
	viper.SetDefault(ConfigDataPath, exePath+string(os.PathSeparator)+"data"+string(os.PathSeparator))
	viper.SetDefault(ConfigPythonPath, exePath+string(os.PathSeparator)+"python"+string(os.PathSeparator)+"python.exe")
	viper.SetDefault(ConfigBrowserPath, "")
	viper.SetDefault(LLMApiBase, "")
	viper.SetDefault(LLMModelName, "")
	viper.SetDefault(LLMApiKey, "")
	_ = viper.WriteConfig()

	return nil
}

func InitConfig(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.Init(); err != nil {
		return err
	}

	// 监控配置文件变化并热加载程序
	// c.WatchConfig()

	return nil
}
