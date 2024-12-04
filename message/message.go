package message

import (
	"fmt"
	"strconv"
	// "github.com/1337rahulraj/beep/rules"
	// "slices"
)

type Message map[string]string

//type Key string

type Id struct {
	Stn      string
	Name     string
	Loc      string
	Geartype string
	Subgear  string
}

func (m Message) GetId() Id {
	id := Id{Stn: m["STN"], Name: m["NAME"], Loc: m["LOC"], Geartype: m["GEARTYPE"], Subgear: m["SUBGEAR"]}
	return id
}

func (i Id) IsEmpty() bool {
	return i.Stn == "" &&
		i.Name == "" &&
		i.Loc == "" &&
		i.Geartype == "" &&
		i.Subgear == ""
}

func (msg Message) GetKeys() []string {
	var keys []string
	for key, value := range msg {
		if isInt(value) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (msg Message) GetStn() string {
	return msg["STN"]

}

// Rule Id is "msg["GEARTYPE"]-msg["SUBGEAR"]"
func (msg Message) GetRuleId() string {
	return fmt.Sprintf("%s-%s", msg["GEARTYPE"], msg["SUBGEAR"])

}

func isInt(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil

}

func (msg Message) GetValue(key string) interface{} {
	return msg[key]
}

