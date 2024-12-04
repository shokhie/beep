package server

import (
	"encoding/json"
	"fmt"
	"github.com/1337rahulraj/beep/rules"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Rule represents the database model for rules
type RuleDB struct {
    ID            string `db:"id"`
    GearTypeSubGear string `db:"gear_type_subgear"`
    Condition     json.RawMessage `db:"condition"`
}

// Database connection configuration
const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    //password = "your_password"
    dbname   = "mydb"
)

// InitDB initializes database connection
func InitDB() (*sqlx.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
        host, port, user, dbname)
        
    db, err := sqlx.Connect("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    // Create the rules table if it doesn't exist
    schema := `
    CREATE TABLE IF NOT EXISTS rules (
        id VARCHAR(36) PRIMARY KEY,
        gear_type_subgear VARCHAR(50) NOT NULL,
        condition JSONB NOT NULL
    );`

    _, err = db.Exec(schema)
    if err != nil {
        return nil, fmt.Errorf("error creating table: %w", err)
    }

    return db, nil
}

// AddRule adds a new rule to the database
func AddRule(db *sqlx.DB, rule rules.Rule) error {
    conditionJSON, err := json.Marshal(rule.Condition)
    if err != nil {
        return fmt.Errorf("error marshaling condition: %w", err)
    }

    query := `
        INSERT INTO rules (id, gear_type_subgear, condition)
        VALUES ($1, $2, $3)`

    _, err = db.Exec(query, rule.Id, rule.GeartypeSubgear, conditionJSON)
    if err != nil {
        return fmt.Errorf("error inserting rule: %w", err)
    }

    return nil
}

// UpdateRule updates an existing rule
func UpdateRule(db *sqlx.DB, rule rules.Rule) error {
    conditionJSON, err := json.Marshal(rule.Condition)
    if err != nil {
        return fmt.Errorf("error marshaling condition: %w", err)
    }

    query := `
        UPDATE rules 
        SET gear_type_subgear = $2, condition = $3
        WHERE id = $1`

    result, err := db.Exec(query, rule.Id, rule.GeartypeSubgear, conditionJSON)
    if err != nil {
        return fmt.Errorf("error updating rule: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("rule with ID %s not found", rule.Id)
    }

    return nil
}

// DeleteRule deletes a rule from the database
func DeleteRule(db *sqlx.DB, ruleID string) error {
    query := `DELETE FROM rules WHERE id = $1`

    result, err := db.Exec(query, ruleID)
    if err != nil {
        return fmt.Errorf("error deleting rule: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("rule with ID %s not found", ruleID)
    }

    return nil
}

// GetAllRules retrieves all rules from the database
func GetAllRules(db *sqlx.DB) (*rules.RulesMap, error) {
    var rulesDB []RuleDB
    query := `SELECT * FROM rules`

    err := db.Select(&rulesDB, query)
    if err != nil {
        return nil, fmt.Errorf("error querying rules: %w", err)
    }

    rulesMap := rules.NewRulesMap()
    for _, ruleDB := range rulesDB {
        var rule rules.Rule
        rule.Id = ruleDB.ID
        rule.GeartypeSubgear = ruleDB.GearTypeSubGear
        
        err = json.Unmarshal(ruleDB.Condition, &rule.Condition)
        if err != nil {
            return nil, fmt.Errorf("error unmarshaling condition: %w", err)
        }

        rulesMap.Rules = append(rulesMap.Rules, rule)
    }

    return rulesMap, nil
}
