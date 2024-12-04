package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/1337rahulraj/beep/message"
	"github.com/Knetic/govaluate"
)

type Rule struct {
	GeartypeSubgear string `json:"Geartype-Subgear"`
	Condition       Condition
	Id              string
}

type Condition struct {
	AlarmName   string   `json:"AlarmName"`
	Expression  string   `json:"Expression"`
	Duration    int      `json:"Duration"`
	Stn         []string `json:"Stn"`
	ConditionId string   `json:"ConditionId"`
	Category    string   `json:"Category"`
	Severity    string   `json:"Severity"`
}

//type RuleItem struct {
//	Rule Rule
//}
//
//type RulePacket struct {
//	//Action string
//	Rule Rule
//}

type RulesMap struct {
	Rules []Rule
}

func NewRulesMap() *RulesMap {
	return &RulesMap{Rules: []Rule{}}
}

func (r *RulesMap) DeserializeRuleFile(filename string) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonData, r)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RulesMap) GetRuleIdsFromRuleFile() []string {
	var SliceOfRuleIds []string

	for _, rule := range r.Rules {
		SliceOfRuleIds = append(SliceOfRuleIds, rule.GeartypeSubgear)
	}

	return SliceOfRuleIds
}

func (r *RulesMap) GetAllConditionsOfGeartypeSubgear(GeartypeSubgear string) []Condition {
	var SliceOfCondition []Condition
	for _, rule := range r.Rules {
		if rule.GeartypeSubgear == GeartypeSubgear {
			SliceOfCondition = append(SliceOfCondition, rule.Condition)
		}
	}

	return SliceOfCondition
}

func (c Condition) BuildParametersMap(msg message.Message) map[string]interface{} {

	// Regex to extract keys
	keyRegex := regexp.MustCompile(`\b([A-Z_][A-Z0-9_]*)\b`)

	keysInExpression := keyRegex.FindAllString(c.Expression, -1)

	ParametersMap := make(map[string]interface{})
	for _, key := range keysInExpression {
		value := msg.GetValue(key)
		convertedValue := convertStringToAppropriateType(value)
		ParametersMap[key] = convertedValue
	}
	return ParametersMap
}

func convertStringToAppropriateType(value interface{}) interface{} {
	// Convert the value to a string
	valueString, ok := value.(string)
	if !ok {
		return value
	}
	// Try converting to int
	if intValue, err := strconv.Atoi(valueString); err == nil {
		return intValue
	}

	// Try converting to float64
	if floatValue, err := strconv.ParseFloat(valueString, 64); err == nil {
		return floatValue
	}

	// Try converting to bool
	if boolValue, err := strconv.ParseBool(valueString); err == nil {
		return boolValue
	}

	// If none of the above conversions work, return the original string
	return value
}

func (c Condition) EvaluateExpression(parametersMap map[string]interface{}) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(c.Expression)
	if err != nil {
		fmt.Println("err1", err)
		return false, err
	}
	result, err := expression.Evaluate(parametersMap)
	if err != nil {
		fmt.Println("err2", err)
		return false, err
	}
	castedBool, ok := result.(bool)
	if !ok {
		return false, errors.New("unable to find bool expr")
	}
	return castedBool, nil
}
