package gomqtt_starter

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	goframeworkmqtt "github.com/kordar/goframework-mqtt"
	"github.com/kordar/gologger"
	"github.com/spf13/cast"
)

var (
	CredentialsProvider   mqtt.CredentialsProvider
	DefaultPublishHandler mqtt.MessageHandler
	OnConnect             mqtt.OnConnectHandler
	OnConnectionLost      mqtt.ConnectionLostHandler
	OnReconnecting        mqtt.ReconnectHandler
	OnConnectAttempt      mqtt.ConnectionAttemptHandler
)

type MqttModule struct {
	name string
	load func(name string, value map[string]string)
}

func NewMqttModule(name string, load func(moduleName string, item map[string]string)) *MqttModule {
	return &MqttModule{name, load}
}

func (m MqttModule) Name() string {
	return m.name
}

func (m MqttModule) Load(value interface{}) {
	items := cast.ToStringMap(value)
	for key, val := range items {
		section := cast.ToStringMapString(val)
		err := goframeworkmqtt.AddMqttInstanceArgs(key, section, CredentialsProvider, DefaultPublishHandler, OnConnect, OnConnectionLost, OnReconnecting, OnConnectAttempt)
		if err != nil {
			logger.Errorf("[gomqtt-starter] 初始化mqtt异常，err=%v", err)
			continue
		}
		if m.load != nil {
			m.load(m.Name(), section)
		}
	}

}

func (m MqttModule) Close() {
}
