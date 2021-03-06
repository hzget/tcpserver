package tcpserver

import (
	"github.com/spf13/viper"
	"log"
)

const (
	WorkerPoolSize = 10
	TaskQueueSize  = 20
)

const (
	TcpSizeAdjust = 0
)

type Config struct {
	tcpserver tcpserverconf
	app       appconf
	protocol  protocolconf
}

type tcpserverconf struct {
	host        string
	port        uint32
	maxconn     uint32
	maxpacksize uint32
}

type appconf struct {
	workerpoolsize uint32
	taskqueuesize  uint32
}

type protocolconf struct {
	tcpsizeadjust uint32
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
			maxconn:     viper.GetUint32("tcpserver.maxconn"),
			maxpacksize: viper.GetUint32("tcpserver.maxpacksize"),
		},
		app: appconf{
			workerpoolsize: viper.GetUint32("app.workerpoolsize"),
			taskqueuesize:  viper.GetUint32("app.taskqueuesize"),
		},
		protocol: protocolconf{
			tcpsizeadjust: viper.GetUint32("protocol.tcpsizeadjust"),
		},
	}

	log.Printf("configuration is %#v", config)
}

func initDefault() {
	viper.SetDefault("tcpserver.host", "0.0.0.0")
	viper.SetDefault("tcpserver.port", 8080)
	viper.SetDefault("tcpserver.maxconn", MaxConn)
	viper.SetDefault("tcpserver.maxpacksize", MaxPackSize)
	viper.SetDefault("app.workerpoolsize", WorkerPoolSize)
	viper.SetDefault("app.taskqueuesize", TaskQueueSize)
	viper.SetDefault("protocol.tcpsizeadjust", TcpSizeAdjust)
}
