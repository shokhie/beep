package eye

import (
	"encoding/json"
	"fmt"
	"github.com/1337rahulraj/beep/message"
	"github.com/1337rahulraj/beep/rules"
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

type IdTimerMap struct {
	Timers map[TimerId]*time.Timer
	mu     sync.RWMutex
}

type TimerId struct {
	Id          message.Id
	ConditionId string
}

type Alarm struct {
	AlarmId    string `json:"alarmId"`
	Id         string `json:"id"`
	IsActive   bool   `json:"isActive"`
	EventCount int    `json:"eventCount"`
	CreatedAt  string `json:"created At"`
	UpdatedAt  string `json:"updated At"`
	ClearedAt  string `json:"cleared At"`
	Expression string `json:"expression"`
	AlarmName  string `json:"alarmName"`
	Category   string `json:"category"`
	Severity   string `json:"severity"`
	Stn        string `json:"stn"`
	Duration   int    `json:"duration"`
	Name       string `json:"name"`
	Location   string `json:"location"`
}

var timerIdAlarmIdMap = make(map[string]string)
var logfileMutex sync.Mutex

func (t TimerId) String() string {
	return fmt.Sprintf("%v-%s", t.Id, t.ConditionId)
}

func NewIdTimerMap() *IdTimerMap {
	return &IdTimerMap{Timers: make(map[TimerId]*time.Timer)}
}

func (idTimerMap *IdTimerMap) AddStartTimer(timerId TimerId, duration int, parametersMap map[string]interface{}, condition rules.Condition) {

	idTimerMap.mu.Lock()
	defer idTimerMap.mu.Unlock()

	_, exists := idTimerMap.Timers[timerId]
	if !exists {
		timer := time.AfterFunc(time.Duration(duration)*time.Minute, func() {
			TriggerAlarm(timerId, condition)
			idTimerMap.StopDeleteTimer(timerId, parametersMap, condition, "trigger")
		})
		idTimerMap.Timers[timerId] = timer

		//startTimerInfo := fmt.Sprintf("xxxxxxxxxxxxxxxxxxxx\nId: %v\nExpression: %s ParametersMap: %v\nTimer Created at time: %v\nxxxxxxxxxxxxxxxxxxxx\n\n", timerId.Id, condition.Expression, parametersMap, time.Now().Format("02-01-2006 03:04:05 PM"))
		//fmt.Printf(startTimerInfo)

	} else {
		//startTimerInfo := fmt.Sprintf("====================\nId: %v\nExpression: %s ParametersMap: %v\nTimer already exists\n====================\n\n", timerId.Id, condition.Expression, parametersMap)
		//fmt.Printf(startTimerInfo)
	}
}

func (idTimerMap *IdTimerMap) StopDeleteTimer(timerId TimerId, parametersMap map[string]interface{}, condition rules.Condition, flag string) {
	idTimerMap.mu.Lock()
	timer, exists := idTimerMap.Timers[timerId]
	if exists {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		delete(idTimerMap.Timers, timerId)

		if flag == "trigger" {
			//fmt.Printf("====================\nTriggered So Deleted\n====================\n\n")
		} else if flag == "delete" {
			//stopTimerInfo := fmt.Sprintf("xxxxxxxxxxxxxxxxxxxx\nId: %v\nExpression: %s ParametersMap: %v\nTimer Stopped and Deleted\nxxxxxxxxxxxxxxxxxxxx\n\n", timerId.Id, condition.Expression, parametersMap)
			//fmt.Println(stopTimerInfo)

			alarmId, exist := timerIdAlarmIdMap[timerId.String()]
			if exist {
				index, alarm, err := GetAlarm(alarmId)
				if err != nil {
					fmt.Println(err)
				}
				if alarm.IsActive == true {
					alarm.IsActive = false
					clearedAt := time.Now().Format("02-01-2006 03:04:05 PM")
					alarm.ClearedAt = clearedAt
					//prettyAlarm0, _ := json.MarshalIndent(alarm, "", "    ")
					//fmt.Println(string(prettyAlarm0))
					updateAlarmToFile("alarms.json", alarm, index)
				} else if alarm.IsActive == false {
					alarm.ClearedAt = time.Now().Format("02-01-2006 03:04:05 PM")
					//prettyAlarm1, _ := json.MarshalIndent(alarm, "", "    ")
					//fmt.Println(string(prettyAlarm1))
					updateAlarmToFile("alarms.json", alarm, index)
				}
			} else {
				//fmt.Printf("====================\nNo entry to close alarm\n====================\n\n")
			}
		}
	} else {
        if flag == "delete" {
            alarmId, exist := timerIdAlarmIdMap[timerId.String()]
			if exist {
				index, alarm, err := GetAlarm(alarmId)
				if err != nil {
					fmt.Println(err)
				}
				if alarm.IsActive == true {
					alarm.IsActive = false
					clearedAt := time.Now().Format("02-01-2006 03:04:05 PM")
					alarm.ClearedAt = clearedAt
					//prettyAlarm0, _ := json.MarshalIndent(alarm, "", "    ")
					//fmt.Println(string(prettyAlarm0))
					updateAlarmToFile("alarms.json", alarm, index)
				} else if alarm.IsActive == false {
					alarm.ClearedAt = time.Now().Format("02-01-2006 03:04:05 PM")
					//prettyAlarm1, _ := json.MarshalIndent(alarm, "", "    ")
					//fmt.Println(string(prettyAlarm1))
					updateAlarmToFile("alarms.json", alarm, index)
				}
			} else {
				//fmt.Printf("====================\nNo entry to close alarm\n====================\n\n")
			}
        }
    }
	idTimerMap.mu.Unlock()
	//fmt.Printf("====================\nId: %v\nExpression: %s ParametersMap: %v\nCool No timer to delete\n====================\n\n", timerId.Id, condition.Expression, parametersMap)
}

//func createClosedAlarm(timerId TimerId, condition rules.Condition) {
//currentTime := time.Now()
//formattedTime := currentTime.Format("02-01-2006 03:04:05 PM")

//closedAlarm := map[string]interface{}{
//	"Id":         timerId.String(),
//	"Closed At":  formattedTime,
//	"Expression": condition.Expression,
//	"Alarm Name": condition.AlarmName,
//	"Category":   condition.Category,
//	"Severity":   condition.Severity,
//	"Stns":       condition.Stn,
//	"Duration":   condition.Duration,
//	"Name":       timerId.Id.Name,
//	"Location":   timerId.Id.Loc,
//}
//result, index, alarm := CheckAlarmEntry(timerId)
//fmt.Println(result, index, alarm)
//fmt.Println("Closed Alarm", closedAlarm)

//writeClosedAlarmToFile("closed.json", closedAlarm)
//}

//func (idTimerMap *IdTimerMap) StopDeleteTimerAfterRuleChange(conditionId string) {
//	fmt.Println("inside delte timer...")
//	idTimerMap.mu.Lock()
//	defer idTimerMap.mu.Unlock()
//	for timerId, timer := range idTimerMap.Timers {
//		fmt.Println("hello")
//		fmt.Println("conditionId: ", conditionId, "timerId: ", timerId)
//		if timerId.ConditionId == conditionId {
//			if !timer.Stop() {
//				select {
//				case <-timer.C:
//				default:
//				}
//			}
//			delete(idTimerMap.Timers, timerId)
//		}
//	}
//}

//func (idTimerMap *IdTimerMap) StopDeleteTimerAfterRuleChange(conditionId string) int {
//	fmt.Println("Inside delete timer function...")
//	fmt.Printf("Condition ID to delete: %s\n", conditionId)
//
//	idTimerMap.mu.Lock()
//	defer idTimerMap.mu.Unlock()
//
//	fmt.Printf("Number of timers in the map: %d\n", len(idTimerMap.Timers))
//
//	if len(idTimerMap.Timers) == 0 {
//		fmt.Println("The timer map is empty.")
//		return 0
//	}
//
//	deletedCount := 0
//
//	for timerId, timer := range idTimerMap.Timers {
//		fmt.Printf("Checking timer - conditionId: %s, timerId: %v\n", timerId.ConditionId, timerId)
//
//		if timerId.ConditionId == conditionId {
//			if !timer.Stop() {
//				select {
//				case <-timer.C:
//				default:
//				}
//			}
//			delete(idTimerMap.Timers, timerId)
//			deletedCount++
//			fmt.Printf("Deleted timer for conditionId: %s, timerId: %v\n", conditionId, timerId)
//		}
//	}
//
//	fmt.Printf("Deleted %d timers for conditionId: %s\n", deletedCount, conditionId)
//	return deletedCount
//}

//var AlarmChannel = make(chan map[string]string, 100)

//var triggeredAlarmMap = make(map[string]string)

func TriggerAlarm(timerId TimerId, condition rules.Condition) {
	alarmId, exist := timerIdAlarmIdMap[timerId.String()]
	if exist {
		index, alarm, err := GetAlarm(alarmId)
		if err != nil {
			fmt.Println(err)
		}
		if alarm.IsActive == true {
			EventCount := alarm.EventCount
			EventCount++
			alarm.EventCount = EventCount
			updatedAt := time.Now().Format("02-01-2006 03:04:05 PM")
			alarm.UpdatedAt = updatedAt
			//fmt.Println(alarm)
			updateAlarmToFile("alarms.json", alarm, index)
		} else if alarm.IsActive == false {
			EventCount := 1
			alarmId := uuid.New().String()

			triggeredAlarm := Alarm{
				AlarmId:    alarmId,
				Id:         timerId.String(),
				IsActive:   true,
				EventCount: EventCount,
				CreatedAt:  time.Now().Format("02-01-2006 03:04:05 PM"),
				Expression: condition.Expression,
				AlarmName:  condition.AlarmName,
				Category:   condition.Category,
				Severity:   condition.Severity,
				Stn:        timerId.Id.Stn,
				Duration:   condition.Duration,
				Name:       timerId.Id.Name,
				Location:   timerId.Id.Loc,
			}
			timerIdAlarmIdMap[timerId.String()] = alarmId
			//prettyAlarm, _ := json.MarshalIndent(triggeredAlarm, "", "    ")
			//fmt.Printf("xxxxxxxxxxxxxxxxxxxx\n%v\nxxxxxxxxxxxxxxxxxxxx\n\n", string(prettyAlarm))
			writeAlarmToFile("alarms.json", triggeredAlarm)
		}
	} else {
		EventCount := 1
		alarmId := uuid.New().String()

		triggeredAlarm := Alarm{
			AlarmId:    alarmId,
			Id:         timerId.String(),
			IsActive:   true,
			EventCount: EventCount,
			CreatedAt:  time.Now().Format("02-01-2006 03:04:05 PM"),
			Expression: condition.Expression,
			AlarmName:  condition.AlarmName,
			Category:   condition.Category,
			Severity:   condition.Severity,
			Stn:        timerId.Id.Stn,
			Duration:   condition.Duration,
			Name:       timerId.Id.Name,
			Location:   timerId.Id.Loc,
		}
		timerIdAlarmIdMap[timerId.String()] = alarmId
		//prettyAlarm, _ := json.MarshalIndent(triggeredAlarm, "", "    ")
		//fmt.Printf("xxxxxxxxxxxxxxxxxxxx\n%v\nxxxxxxxxxxxxxxxxxxxx\n\n", string(prettyAlarm))
		writeAlarmToFile("alarms.json", triggeredAlarm)
	}

}

func GetAlarm(alarmId string) (int, Alarm, error) {
	logfileMutex.Lock()
	defer logfileMutex.Unlock()

	var alarms struct {
		Alarms []Alarm `json:"Alarms"`
	}

	file, err := os.OpenFile("alarms.json", os.O_RDONLY, 0644)
	if err != nil {
		return 0, Alarm{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&alarms); err != nil {
		return 0, Alarm{}, err
	}

	for index, alarm := range alarms.Alarms {
		if alarmId == alarm.AlarmId {
			return index, alarm, nil
		}
	}
	return 0, Alarm{}, fmt.Errorf("Alarm does not exist for alarmId %v\n", alarmId)
}

//func IsAlarmEntryExists(alarmId string) (bool, int, Alarm, error) {
//	logfileMutex.Lock()
//	defer logfileMutex.Unlock()
//
//	var alarms struct {
//		Alarms []Alarm `json:"Alarms"`
//	}
//
//	file, err := os.OpenFile("alarms.json", os.O_RDONLY, 0644)
//	if err != nil {
//		return false, 0, Alarm{}, err
//	}
//	defer file.Close()
//
//	decoder := json.NewDecoder(file)
//	if err := decoder.Decode(&alarms); err != nil {
//		return false, 0, Alarm{}, err
//	}
//
//	for index, alarm := range alarms.Alarms {
//		if alarmId == alarm.AlarmId {
//			return true, index, alarm, nil
//		}
//	}
//	return false, 0, Alarm{}, fmt.Errorf("Alarm does not exist for alarmId %v", alarmId)
//}

func updateAlarmToFile(filename string, alarm Alarm, index int) {
	logfileMutex.Lock()
	defer logfileMutex.Unlock()

	// Read existing alarms
	var alarms struct {
		Alarms []Alarm `json:"Alarms"`
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening alarm file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	if fileInfo.Size() > 0 {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&alarms); err != nil {
			fmt.Println("Error decoding existing alarms:", err)
			return
		}
	}

	//Update alarm
	alarms.Alarms[index] = alarm

	// Write updated alarms back to file
	file.Seek(0, 0)
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(alarms); err != nil {
		fmt.Println("Error encoding alarms:", err)
		return
	}
	fmt.Println("Alarm Updated in file")
}

func writeAlarmToFile(filename string, alarm Alarm) {
	logfileMutex.Lock()
	defer logfileMutex.Unlock()

	// Read existing alarms
	var alarms struct {
		Alarms []Alarm `json:"Alarms"`
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening alarm file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	if fileInfo.Size() > 0 {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&alarms); err != nil {
			fmt.Println("Error decoding existing alarms:", err)
			return
		}
	}

	// Append new alarm
	alarms.Alarms = append(alarms.Alarms, alarm)

	// Write updated alarms back to file
	file.Seek(0, 0)
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(alarms); err != nil {
		fmt.Println("Error encoding alarms:", err)
		return
	}

	fmt.Println("Alarm successfully written to file")
}
