package tcpserver

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	tcpserver tcpserverconf
}

type tcpserverconf struct {
	host        string
	port        uint32
	maxpacksize uint32
}

// globals
var config *Config

func init() {
	getConfig()
	initDefault()
	initGlobals()
}

func getConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Fatal error config file: %v", err)
	}
}

func initGlobals() {
	config = &Config{
		tcpserver: tcpserverconf{
			host:        viper.GetString("tcpserver.host"),
			port:        viper.GetUint32("tcpserver.port"),
			maxpacksize: viper.GetUint32("tcpserver.maxpacksize"),
		},
	}

	log.Printf("configuration is %#v", config)
}

func initDefault() {
	viper.SetDefault("tcpserver.host", "tcpserver")
	viper.SetDefault("tcpserver.port", 8080)
	viper.SetDefault("tcpserver.maxpacksize", MaxPackSize)
}
