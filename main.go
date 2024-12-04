package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/1337rahulraj/beep/eye"
	mqtt "github.com/1337rahulraj/beep/mqttclient"
	"github.com/1337rahulraj/beep/rules"
	"github.com/1337rahulraj/beep/server"
)

var AlarmChannel = make(chan map[string]string, 100)

func main() {
	go server.StartServer(":1337")
	host := address
	port := "1883"
	topics := map[string]byte{"sensor_NOQ": 1, "sensor_SUA": 1, "sensor_JGM": 1} //"sensor_data5": 1, "sensor_LKR": 1}

	mqttClient := mqtt.Connect(host + ":" + port)
	msgChannel := mqtt.RcvStream(mqttClient, topics)

	idTimerMap := eye.NewIdTimerMap()

	rulesMap := rules.NewRulesMap()
	rulesFilename := os.Args[1]

	err := rulesMap.DeserializeRuleFile(rulesFilename)
	if err != nil {
		fmt.Println("Error deserializing rule file:", err)
	}

	doneChan := make(chan bool)

	go func(doneChan chan bool) {
		defer func() {
			doneChan <- true
		}()

		err := watchFile("rules.json")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("File has been changed")
	}(doneChan)

	go func(doneChan chan bool) {
		if <-doneChan {
			err := rulesMap.DeserializeRuleFile(rulesFilename)
			//fmt.Println(rulesMap)
			if err != nil {
				fmt.Println(err)
			}
		}
	}(doneChan)

	for msg := range msgChannel {

		ruleId := msg.GetRuleId()
		stn := msg.GetStn()

		if slices.Contains(rulesMap.GetRuleIdsFromRuleFile(), ruleId) {
			for _, condition := range rulesMap.GetAllConditionsOfGeartypeSubgear(ruleId) {
				if slices.Contains(condition.Stn, stn) {

					parametersMap := condition.BuildParametersMap(msg)

					result, err := condition.EvaluateExpression(parametersMap)
					if err != nil {
						fmt.Println(err)
					}

					timerId := eye.TimerId{Id: msg.GetId(), ConditionId: condition.ConditionId}

					if result {
						idTimerMap.AddStartTimer(timerId, condition.Duration, parametersMap, condition)
					} else {
						idTimerMap.StopDeleteTimer(timerId, parametersMap, condition, "delete")
					}

				}
			}
		}
	}
}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
