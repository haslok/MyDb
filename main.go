// version: 1.0.2
package MyDb

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"sync"
	"strings"
)

// Table represents a table in the database
type Table struct {
	Columns []string               // Column names
	Rows    []map[string]string    // Rows of data as a map of column names to values
	mu      sync.Mutex             // Mutex for concurrent access
}

// Database represents a database with a collection of tables
type Database struct {
	Name   string             // Name of the database
	Tables map[string]*Table  // Map of table names to tables
	mu     sync.Mutex         // Mutex for concurrent access
}

// NewDatabase creates a new database with the given name
func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Tables: make(map[string]*Table),
	}
}

// CreateTable creates a new table in the database
func (db *Database) CreateTable(name string, columns []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Validate table and column names
	if !isValidName(name) {
		return fmt.Errorf("invalid table name: %s", name)
	}
	for _, col := range columns {
		if !isValidName(col) {
			return fmt.Errorf("invalid column name: %s", col)
		}
	}

	// Check if the table already exists
	if _, exists := db.Tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}

	// Create the table
	db.Tables[name] = &Table{Columns: columns}
	return nil
}

// InsertInto inserts a row of data into the specified table
func (db *Database) InsertInto(tableName string, data map[string]string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the table exists
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	// Validate the data columns
	for key := range data {
		if !contains(table.Columns, key) {
			return fmt.Errorf("column %s does not exist in table %s", key, tableName)
		}
	}

	// Lock the table and insert the row
	table.mu.Lock()
	defer table.mu.Unlock()
	table.Rows = append(table.Rows, data)
	return nil
}


// Delete removes rows from the specified table that match all the given conditions
func (db *Database) Delete(tableName string, conditions map[string]string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the table exists
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	// Lock the table to ensure thread safety
	table.mu.Lock()
	defer table.mu.Unlock()

	// Filter rows that do not match the conditions
	var remainingRows []map[string]string
	for _, row := range table.Rows {
		match := true
		for col, val := range conditions {
			if row[col] != val {
				match = false
				break
			}
		}
		if !match {
			remainingRows = append(remainingRows, row)
		}
	}

	// Update the table with remaining rows
	table.Rows = remainingRows
	return nil
}



// UpdateData updates rows in the specified table based on a condition
func (db *Database) UpdateData(tableName string, condition func(row map[string]string) bool, data map[string]string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the table exists
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	// Validate that the data map matches the table columns
	for key := range data {
		if !contains(table.Columns, key) {
			return fmt.Errorf("column %s does not exist in table %s", key, tableName)
		}
	}

	// Lock the table and update matching rows
	table.mu.Lock()
	defer table.mu.Unlock()
	for i, row := range table.Rows {
		if condition(row) {
			// Update the row with the new data
			for key, value := range data {
				row[key] = value
			}
			table.Rows[i] = row
		}
	}
	return nil
}

// SearchRows searches for rows in the specified table based on a condition
func (db *Database) SearchRows(tableName string, condition func(row map[string]string) bool) ([]map[string]string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the table exists
	table, exists := db.Tables[tableName]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	// Lock the table and search for rows matching the condition
	table.mu.Lock()
	defer table.mu.Unlock()

	var results []map[string]string
	for _, row := range table.Rows {
		if condition(row) {
			results = append(results, row)
		}
	}
	return results, nil
}

// SelectTable selects a table from a CSV file
func (db *Database) SelectTable(tableName string) (*Table, error) {
	// Open the table's CSV file
	file, err := os.Open(fmt.Sprintf("%s/%s.csv", db.Name, tableName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	columns, err := reader.Read()
	if err != nil {
		return nil, err
	}

	table := &Table{
		Columns: columns,
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Convert rows to map[string]string
	var mappedRows []map[string]string
	for _, row := range rows {
		mappedRow := make(map[string]string)
		for i, col := range columns {
			mappedRow[col] = row[i]
		}
		mappedRows = append(mappedRows, mappedRow)
	}

	table.Rows = mappedRows

	return table, nil
}

// Save saves the database to a directory and creates a CSV file for each table
func (db *Database) Save() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Ensure the database directory exists
	if err := os.MkdirAll(db.Name, os.ModePerm); err != nil {
		return err
	}

	// Save each table as a CSV file
	for tableName, table := range db.Tables {
		file, err := os.Create(fmt.Sprintf("%s/%s.csv", db.Name, tableName))
		if err != nil {
			return err
		}

		writer := csv.NewWriter(file)
		// Write column headers
		if err := writer.Write(table.Columns); err != nil {
			file.Close()
			return err
		}

		// Write rows
		for _, row := range table.Rows {
			var rowData []string
			for _, col := range table.Columns {
				rowData = append(rowData, row[col])
			}
			if err := writer.Write(rowData); err != nil {
				file.Close()
				return err
			}
		}

		writer.Flush()
		file.Close()
	}

	return nil
}

// isValidName checks if a name is valid (alphanumeric with underscores)
func isValidName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	return matched
}

// contains checks if a string is present in a slice of strings
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

//execute command
// MyDb executes SQL-like commands for the database
func (db *Database) Command(command string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Remove unnecessary spaces
	command = regexp.MustCompile(`\s+`).ReplaceAllString(command, " ")
	command = strings.TrimSpace(command)

	// Command parsing
	parts := strings.SplitN(command, " ", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid command: %s", command)
	}

	action := strings.ToLower(parts[0])
	switch action {
	case "delete":
		// Example: DELETE FROM users WHERE name = ahmad
		matches := regexp.MustCompile(`delete from (\w+) where (.+)`).FindStringSubmatch(strings.ToLower(command))
		if len(matches) != 3 {
			return fmt.Errorf("invalid DELETE command: %s", command)
		}
		tableName := matches[1]
		conditions := parseConditions(matches[2])
		return db.Delete(tableName, conditions)

	case "update":
		// Example: UPDATE users SET name = ahmad WHERE id = 1
		matches := regexp.MustCompile(`update (\w+) set (.+) where (.+)`).FindStringSubmatch(strings.ToLower(command))
		if len(matches) != 4 {
			return fmt.Errorf("invalid UPDATE command: %s", command)
		}
		tableName := matches[1]
		data := parseConditions(matches[2])
		conditions := parseConditions(matches[3])
		return db.UpdateData(tableName, func(row map[string]string) bool {
			return matchConditions(row, conditions)
		}, data)

	case "get", "select":
		// Example: GET FROM users WHERE name = ahmad
		matches := regexp.MustCompile(`get from (\w+) where (.+)`).FindStringSubmatch(strings.ToLower(command))
		if len(matches) != 3 {
			return fmt.Errorf("invalid GET command: %s", command)
		}
		tableName := matches[1]
		conditions := parseConditions(matches[2])
		rows, err := db.SearchRows(tableName, func(row map[string]string) bool {
			return matchConditions(row, conditions)
		})
		if err != nil {
			return err
		}
		fmt.Println("Results:", rows)
		return nil

	case "insert":
		// Example: INSERT INTO users 1, ahmad, 55
		matches := regexp.MustCompile(`insert to (\w+) (.+)`).FindStringSubmatch(strings.ToLower(command))
		if len(matches) != 3 {
			return fmt.Errorf("invalid INSERT command: %s", command)
		}
		tableName := matches[1]
		values := strings.Split(matches[2], ",")
		valuesMap := make(map[string]string)
		if table, exists := db.Tables[tableName]; exists {
			if len(values) != len(table.Columns) {
				return fmt.Errorf("number of values does not match columns in table %s", tableName)
			}
			for i, col := range table.Columns {
				valuesMap[col] = strings.TrimSpace(values[i])
			}
		} else {
			return fmt.Errorf("table %s does not exist", tableName)
		}
		return db.InsertInto(tableName, valuesMap)

	case "create":
		// Example: CREATE TABLE users has id, name, age
		matches := regexp.MustCompile(`create table (\w+) has (.+)`).FindStringSubmatch(strings.ToLower(command))
		if len(matches) != 3 {
			return fmt.Errorf("invalid CREATE command: %s", command)
		}
		tableName := matches[1]
		columns := strings.Split(matches[2], ",")
		for i := range columns {
			columns[i] = strings.TrimSpace(columns[i])
		}
		return db.CreateTable(tableName, columns)

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// Helper function to parse conditions
func parseConditions(conditionStr string) map[string]string {
	conditions := make(map[string]string)
	pairs := strings.Split(conditionStr, " and ")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			conditions[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return conditions
}

// Helper function to match conditions against a row
func matchConditions(row map[string]string, conditions map[string]string) bool {
	for key, value := range conditions {
		if row[key] != value {
			return false
		}
	}
	return true
}
