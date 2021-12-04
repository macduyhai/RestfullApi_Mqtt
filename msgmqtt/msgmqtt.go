package msgmqtt

import (
	"beetai_cloud_config_json/box"
	"beetai_cloud_config_json/file"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttCmsBi : Khoi tao doi tuong mqtt
var MqttCmsBi mqtt.Client

// CmsHostBi : host MQTT
// const CmsHostBi string = "tcp://broker.hivemq.com:1883"
const CmsHostBi string = "tcp://vuaop.com:1883"

// CmsAccessTokenBi : User
const CmsAccessTokenBi = ""

// CmsPassBi : Password
const CmsPassBi = ""

// idBox : Phân định device
var idBox = ""

// CmsTopicIn : Server to Box
var CmsTopicIn = "/v1/devices/NTQ/" + idBox + "/telemetry"

// CmsTopicOut : Box to Server
var CmsTopicOut = "/v1/devices/NTQ/" + idBox + "/request/"

//PublishData : Function
func PublishData(idBox string, payload string) { // idBox : Mac of device
	CmsTopicIn = "/v1/devices/NTQ/" + idBox + "/telemetry"
	CmsTopicOut = "/v1/devices/NTQ/" + idBox + "/request/"
	// fmt.Println("idBox: " + idBox)
	fmt.Println("TOPIC IN :" + CmsTopicIn)
	// fmt.Println("TOPIC OUT :" + CmsTopicOut)
	// If Test device with static Topic
	//CmsTopicIn = "TNQ_MQTT"
	// fmt.Println("Test device with static Topic")
	// fmt.Println("TOPIC IN :" + CmsTopicIn)
	//var payload string = "{" + "\"ip_private\":" + "\"" + ip + "\"" +"," + "\"box_id\":" + "\"" + id_cam + "\"" + "}"

	fmt.Printf("Payload = %v\n", payload)
	Token1 := MqttCmsBi.Publish(CmsTopicIn, 1, false, payload)
	if Token1.Wait() && Token1.Error() != nil {
		fmt.Printf("Error Publish message : %v\n", Token1.Error())
	} else {
		fmt.Println("Send message")
	}
	fmt.Println("-------------------------------------------------------------")
}

// MqttBegin : Khoi tao MQTT
func MqttBegin() {

	OptsCmsBI := mqtt.NewClientOptions()
	OptsCmsBI.AddBroker(CmsHostBi)
	OptsCmsBI.SetUsername(CmsAccessTokenBi)
	OptsCmsBI.SetPassword(CmsPassBi)
	OptsCmsBI.SetCleanSession(true)
	OptsCmsBI.SetConnectionLostHandler(MQTTLostConnectHandler)
	OptsCmsBI.SetOnConnectHandler(MQTTOnConnectHandler)

	MqttCmsBi = mqtt.NewClient(OptsCmsBI)
	if Token1 := MqttCmsBi.Connect(); Token1.Wait() && Token1.Error() == nil {
		fmt.Println("MQTT CMS  Connected")
		MqttCmsBi.Subscribe(CmsTopicOut, 0, MqttMessageHandler)
	} else {
		fmt.Println("MQTT CMS  cant not Connected 1234")
		fmt.Printf("Loi CMS  : %v \n", Token1.Error())
		fmt.Println("-------------------")
	}
}

// number_repub : check number push msg faild
var number_repub = 0

//MQTTLostConnectHandler: Check lost connect mqtt server - can't publish msg
func MQTTLostConnectHandler(c mqtt.Client, err error) {
	// c.Disconnect(10)
	// MqttStt = false
	number_repub = number_repub + 1
	fmt.Println("MQTT CMS  Lost Connect")
	fmt.Println("Number reconnect publish msg: " + strconv.Itoa(number_repub))
	fmt.Println(err)
	fmt.Println("------------------------------------------------------------------")
}

var number_resub = 0

// Check lost connect mqtt server - can't subcriber msg from cms
func MQTTOnConnectHandler(client mqtt.Client) {
	number_resub = number_resub + 1
	fmt.Println("Lostconnect chanel subcriber: MQTT_OnConnectHandler")
	fmt.Println("Number reconnect subcriber msg: " + strconv.Itoa(number_resub))
	// client.Unsubscribe(CmsTopicOut)
	// time.Sleep(10)
	//
	if token := client.Subscribe(CmsTopicOut, 0, MqttMessageHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		fmt.Println("Subscriber Error MQTT")

	} else {
		fmt.Println("Subcriber is MQTTOnConnectHandler () OKIE ")
	}
}
func MqttMessageHandler(MqttBI mqtt.Client, message mqtt.Message) {
	fmt.Println("=================== MqttMessageHandler ====================")
	fmt.Printf("Message %s\n", message)
	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MSG:\n %s\n", message.Payload())
	fmt.Println("=========================== + = + ==========================")
	dec := json.NewDecoder(bytes.NewReader(message.Payload()))
	var list map[string]interface{}
	if err := dec.Decode(&list); err != nil {
		fmt.Printf("Error:%v\n", err)
		fmt.Println("Message:Loi form message\n")
	} else {
		//***********************************************//
		if list["method"] == "sub_begin" { // Test subcriber
			fmt.Println("================>> Subcriber is OKIE <========== \n")

		} else if list["method"] == "upgrade_engine" {
			file.Println(`list["method"] == "upgrade_engine"`)
			s, err := json.Marshal(list["params"])
			if err != nil {
				//fmt.Println(err)
			}
			topic := strings.Replace(message.Topic(), "request", "response", 1)
			msg := `{"method":"upgrade_engine","status":` + strconv.Itoa(1) + `}`
			//file.Println(topic)
			file.Write_log("Message:"+msg+" \n", box.Path_log_luncher)
			CmsResponse(MqttBI, topic, msg)
			params := string(s)
			fmt.Println(params)
			// UpgradeEngine(params)
		} else {
			//fmt.Println(list)
		}
	}

}

// CmsResponse : Phan hoi message tu box to Server
func CmsResponse(c mqtt.Client, topic string, msg string) {
	c.Publish(topic, 0, false, msg)
	file.Println("Publish message done")
}
