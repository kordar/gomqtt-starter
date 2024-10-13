package gomqtt_starter

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	goframeworkmqtt "github.com/kordar/goframework-mqtt"
	logger "github.com/kordar/gologger"
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
	load func(moduleName string, itemId string, item map[string]string)
}

func NewMqttModule(name string, load func(moduleName string, itemId string, item map[string]string)) *MqttModule {
	return &MqttModule{name, load}
}

func (m MqttModule) Name() string {
	return m.name
}

func (m MqttModule) _load(id string, cfg map[string]string) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	err := goframeworkmqtt.AddMqttInstanceArgs(id, cfg, CredentialsProvider, DefaultPublishHandler, OnConnect, OnConnectionLost, OnReconnecting, OnConnectAttempt)
	if err != nil {
		logger.Errorf("[gomqtt-starter] 初始化mqtt异常，err=%v", err)
		return
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		logger.Debugf("[%s] triggering custom loader completion", m.Name())
	}

	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
}

func (m MqttModule) Load(value interface{}) {

	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, cast.ToStringMapString(value))
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMapString(item))
	}

}

func (m MqttModule) Close() {
}
