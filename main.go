// version : 1.0.2
package MyDb

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"sync"
)

// Table represents a table in the database
type Table struct {
	Columns []string   // Column names
	Rows    [][]string // Rows of data
	mu      sync.Mutex // Mutex for concurrent access
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
func (db *Database) InsertInto(tableName string, data []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the table exists
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	// Validate data length
	if len(data) != len(table.Columns) {
		return fmt.Errorf("data length does not match number of columns")
	}

	// Lock the table and insert the row
	table.mu.Lock()
	defer table.mu.Unlock()
	table.Rows = append(table.Rows, data)
	return nil
}

// UpdateData updates rows in the specified table based on a condition
func (db *Database) UpdateData(tableName string, condition func(row []string) bool, data []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the table exists
	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	// Validate data length
	if len(data) != len(table.Columns) {
		return fmt.Errorf("invalid data for table %s: expected %d columns, got %d", tableName, len(table.Columns), len(data))
	}

	// Lock the table and update matching rows
	table.mu.Lock()
	defer table.mu.Unlock()
	for i, row := range table.Rows {
		if condition(row) {
			table.Rows[i] = data
		}
	}
	return nil
}


func (db *Database) SelectTable(tableName string) (*Table, error){
	// select from tablename.csv file
	file, err := os.Open(fmt.Sprintf("%s/%s.csv",db.Name, tableName))

	if err != nil {
		return nil, err
	}
	defer file.Close()
	// read csv file
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

	table.Rows = rows

	return table, nil

}



// Save saves the database to a directory and creates a CSV file for each table
func (db *Database) Save() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Create the directory for the database
	if err := os.MkdirAll(db.Name, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for database: %v", err)
	}

	// Save each table as a CSV file
	for tableName, table := range db.Tables {
		file, err := os.Create(fmt.Sprintf("%s/%s.csv", db.Name, tableName))
		if err != nil {
			return fmt.Errorf("failed to create CSV file for table %s: %v", tableName, err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		// Write column headers
		if err := writer.Write(table.Columns); err != nil {
			return fmt.Errorf("failed to write column headers for table %s: %v", tableName, err)
		}

		// Write rows
		for _, row := range table.Rows {
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write data to table %s: %v", tableName, err)
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("failed to flush writer for table %s: %v", tableName, err)
		}
	}

	return nil
}

// isValidName validates table and column names
func isValidName(name string) bool {
	validName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validName.MatchString(name)
}
