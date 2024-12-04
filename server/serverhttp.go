//package server
//
//import (
//	"encoding/json"
//	"fmt"
//	"log"
//	"net/http"
//	"os"
//	"sync"
//
//	"github.com/1337rahulraj/beep/eye"
//	"github.com/1337rahulraj/beep/rules"
//	"github.com/gin-gonic/gin"
//)
//
//var (
//	fileMutex   sync.Mutex
//	rulesMap    *rules.RulesMap
//	alarms      []string
//	alarmsMutex sync.Mutex
//)
//
//func StartServer(address string) {
//	router := gin.Default()
//	router.Any("/rules", handleRules)
//    router.Run(address)
//}
//
//func handleRules(c *gin.Context) {
//    var rule rules.Rule
//
//	switch c.Request.Method {
//	case http.MethodGet:
//		file, err := os.ReadFile("rules.json")
//		if err != nil {
//			log.Println(err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		err = json.Unmarshal(file, &rulesMap)
//		if err != nil {
//			log.Println(err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(http.StatusOK, rulesMap)
//	case http.MethodPost:
//		err := c.BindJSON(&rule)
//        fmt.Println(rule)
//		if err != nil {
//			log.Println(err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		err = addRule("rules.json", rule)
//		if err != nil {
//			log.Println("Error adding rule", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(http.StatusCreated, gin.H{"message": "Rule added successfully"})
//    case http.MethodPut:
//        err := c.BindJSON(&rule)
//        if err != nil {
//            log.Println(err)
//            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//            return
//        }
//        err = updateRule("rules.json", rule)
//        if err != nil {
//            log.Println("Error updating rule", err)
//            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//            return
//        }
//        c.JSON(http.StatusOK, gin.H{"message": "Rule updated successfully"})
//    case http.MethodDelete:
//        err := c.BindJSON(&rule)
//        if err != nil {
//            log.Println(err)
//            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//            return
//        }
//        err = deleteRule("rules.json", rule)
//        if err != nil {
//            log.Println("Error deleting rule", err)
//            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//            return
//        }
//        c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})
//    default:
//        c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed"})
//    }
//}
//
////func StartServerHTTP(address string) {
////	fmt.Println("Server: I am Starting...")
////
////	rulesMap = rules.NewRulesMap()
////	rulesMap.DeserializeRuleFile("rules.json")
////
////	http.HandleFunc("/rules", handleRules)
////	http.HandleFunc("/alarm/history", handleHistory)
////	http.HandleFunc("/alarm/active", handleActiveAlarms)
////	http.HandleFunc("/alarm/closed", handleClosedAlarms)
////
////	fmt.Printf("Server: Listening on address %s\n", address)
////	err := http.ListenAndServe(address, nil)
////	if err != nil {
////		fmt.Println("Error starting server:", err)
////		os.Exit(1)
////	}
////
////}
//
//func handleActiveAlarms(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	file, err := os.ReadFile("alarms.json")
//	if err != nil {
//		http.Error(w, "Error reading closed alarms file", http.StatusInternalServerError)
//		return
//	}
//
//	var alarms struct {
//		Alarms []eye.Alarm `json:"Alarms"`
//	}
//	err = json.Unmarshal(file, &alarms)
//	if err != nil {
//		fmt.Println("Error Unmarshalling in handleClosedAlarms", err)
//	}
//
//	var activeAlarms []eye.Alarm
//	for _, alarm := range alarms.Alarms {
//		if alarm.IsActive == true {
//			activeAlarms = append(activeAlarms, alarm)
//		}
//	}
//	response, err := json.Marshal(activeAlarms)
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(response)
//	if err != nil {
//		http.Error(w, "Error writing response", http.StatusInternalServerError)
//		return
//	}
//}
//
//func handleClosedAlarms(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	file, err := os.ReadFile("alarms.json")
//	if err != nil {
//		http.Error(w, "Error reading closed alarms file", http.StatusInternalServerError)
//		return
//	}
//
//	var alarms struct {
//		Alarms []eye.Alarm `json:"Alarms"`
//	}
//	err = json.Unmarshal(file, &alarms)
//	if err != nil {
//		fmt.Println("Error Unmarshalling in handleClosedAlarms", err)
//	}
//
//	var closedAlarms []eye.Alarm
//	for _, alarm := range alarms.Alarms {
//		if alarm.IsActive == false {
//			closedAlarms = append(closedAlarms, alarm)
//		}
//	}
//	response, err := json.Marshal(closedAlarms)
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(response)
//	if err != nil {
//		http.Error(w, "Error writing response", http.StatusInternalServerError)
//		return
//	}
//}
//
//func handleHistory(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	// Read the alarms.json file
//	file, err := os.ReadFile("alarms.json")
//	if err != nil {
//		http.Error(w, "Error reading alarms file", http.StatusInternalServerError)
//		return
//	}
//
//	// Set content type header
//	w.Header().Set("Content-Type", "application/json")
//
//	// Write the file contents to the response
//	_, err = w.Write(file)
//	if err != nil {
//		http.Error(w, "Error writing response", http.StatusInternalServerError)
//		return
//	}
//}
//
////func handleRules(w http.ResponseWriter, r *http.Request) {
////	switch r.Method {
////	case http.MethodGet:
////		json.NewEncoder(w).Encode(rulesMap)
////	case http.MethodPost:
////		var rule rules.Rule
////		json.NewDecoder(r.Body).Decode(&rule)
////		err := addRule("rules.json", rule)
////		if err != nil {
////			w.WriteHeader(http.StatusInternalServerError)
////			return
////		}
////		w.WriteHeader(http.StatusCreated)
////	case http.MethodPut:
////		var rule rules.Rule
////		json.NewDecoder(r.Body).Decode(&rule)
////		err := updateRule("rules.json", rule)
////		if err != nil {
////			fmt.Println("Error updating rule:", err)
////			w.WriteHeader(http.StatusInternalServerError)
////			return
////		}
////		w.WriteHeader(http.StatusOK)
////	case http.MethodDelete:
////		var rule rules.Rule
////		json.NewDecoder(r.Body).Decode(&rule)
////		err := deleteRule("rules.json", rule)
////		if err != nil {
////			fmt.Println("Error deleting rule:", err)
////			w.WriteHeader(http.StatusInternalServerError)
////			return
////		}
////		w.WriteHeader(http.StatusOK)
////	default:
////		w.WriteHeader(http.StatusMethodNotAllowed)
////	}
////}
//
//func writeToFile(filename string, rulesMap *rules.RulesMap) error {
//	//fileMutex.Lock()
//	//defer fileMutex.Unlock()
//
//	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0644)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	jsonEncoder := json.NewEncoder(file)
//	jsonEncoder.SetEscapeHTML(false)
//	jsonEncoder.SetIndent("", "    ")
//
//	err = jsonEncoder.Encode(rulesMap)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func addRule(filename string, rule rules.Rule) error {
//	fileMutex.Lock()
//	defer fileMutex.Unlock()
//
//	rulesMap := rules.NewRulesMap()
//	rulesMap.DeserializeRuleFile(filename)
//
//	rulesMap.Rules = append(rulesMap.Rules, rule)
//
//	return writeToFile(filename, rulesMap)
//}
//
//func updateRule(filename string, rule rules.Rule) error {
//	fileMutex.Lock()
//	defer fileMutex.Unlock()
//
//	//rulesMap := rules.NewRulesMap()
//	rulesMap.DeserializeRuleFile(filename)
//
//	for i, ruleItem := range rulesMap.Rules {
//		if ruleItem.Id == rule.Id {
//			rulesMap.Rules[i] = rule
//			break
//		}
//	}
//
//	return writeToFile(filename, rulesMap)
//}
//
//func deleteRule(filename string, rule rules.Rule) error {
//	//fileMutex.Lock()
//	//defer fileMutex.Unlock()
//
//	//rulesMap := rules.NewRulesMap()
//	rulesMap.DeserializeRuleFile(filename)
//
//	for i, ruleItem := range rulesMap.Rules {
//		if ruleItem.Id == rule.Id {
//			rulesMap.Rules = append(rulesMap.Rules[:i], rulesMap.Rules[i+1:]...)
//			break
//		}
//	}
//
//	return writeToFile(filename, rulesMap)
//
//}

package server

import (
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    "github.com/1337rahulraj/beep/rules"
    "log"
    "net/http"
)

var db *sqlx.DB

func StartServer(address string) {
    var err error
    db, err = InitDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    router := gin.Default()
    router.Any("/rules", handleRules)
    router.Run(address)
}

func handleRules(c *gin.Context) {
    var rule rules.Rule

    switch c.Request.Method {
    case http.MethodGet:
        rulesMap, err := GetAllRules(db)
        if err != nil {
            log.Println(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, rulesMap)

    case http.MethodPost:
        if err := c.BindJSON(&rule); err != nil {
            log.Println(err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        if err := AddRule(db, rule); err != nil {
            log.Println("Error adding rule:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "Rule added successfully"})

    case http.MethodPut:
        if err := c.BindJSON(&rule); err != nil {
            log.Println(err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        if err := UpdateRule(db, rule); err != nil {
            log.Println("Error updating rule:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Rule updated successfully"})

    case http.MethodDelete:
        if err := c.BindJSON(&rule); err != nil {
            log.Println(err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        if err := DeleteRule(db, rule.Id); err != nil {
            log.Println("Error deleting rule:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})

    default:
        c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed"})
    }
}
