package mqttclient

import (
	"encoding/json"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	message "github.com/1337rahulraj/beep/message"
)

func Connect(address string) mqtt.Client {
	opts := mqtt.NewClientOptions() // set options on client
	opts.AddBroker(address)         // set address
	opts.SetCleanSession(true)      // only receive current data from broker
	client := mqtt.NewClient(opts)  // create a client
	token := client.Connect()       // connect client to broker

	// wait for a max of 1 minute for the flow associated with the token to complete
	if !token.WaitTimeout(1*time.Minute) && token.Error() != nil {
		log.Fatal(token.Error())
	}

        log.Printf("Beep: Connected to Mqtt Broker Address: %s", address)
	return client
}

func RcvStream(client mqtt.Client, topics map[string]byte) chan message.Message {
	msgChannel := make(chan message.Message)
	go func() {
		var messageHandler mqtt.MessageHandler = func(c mqtt.Client, msg mqtt.Message) {
			data := make(map[string]string)
			err := json.Unmarshal(msg.Payload(), &data)
			if err != nil {
				//fmt.Println("Msg deserializaiton error", err)
			}
			//fmt.Println(reflect.TypeOf(data["RH"]))

			msgChannel <- data
		}

		token := client.SubscribeMultiple(topics, messageHandler)
		if !token.WaitTimeout(30*time.Second) && token.Error() != nil {
			log.Print("subscribe Error", token.Error())
			return
		}
	}()

	return msgChannel

}
