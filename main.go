package MyDb

import (
    "fmt"
    "os"
    "encoding/csv"
    "sync"
    "regexp"
)

type Table struct {
    Columns []string
    Rows    [][]string
    mu      sync.Mutex
}

type Database struct {
    Name   string
    Tables map[string]*Table
    mu     sync.Mutex
}

func NewDatabase(name string) *Database {
    return &Database{Name: name, Tables: make(map[string]*Table)}
}

// CreateTable creates a new table in the database if it doesn't already exist
func (db *Database) CreateTable(name string, columns []string) error {
    db.mu.Lock()
    defer db.mu.Unlock()

    if !isValidName(name) {
        return fmt.Errorf("Invalid table name: %s", name)
    }

    for _, col := range columns {
        if !isValidName(col) {
            return fmt.Errorf("Invalid column name: %s", col)
        }
    }

    if _, exists := db.Tables[name]; exists {
        return fmt.Errorf("Table %s already exists", name)
    }

    db.Tables[name] = &Table{Columns: columns}
    return nil
}

// InsertInto inserts a row of data into the specified table
func (db *Database) InsertInto(tableName string, data []string) error {
    db.mu.Lock()
    defer db.mu.Unlock()

    if len(data) == 0 {
        return fmt.Errorf("No data to insert")
    }

    table, exists := db.Tables[tableName]
    if !exists {
        return fmt.Errorf("table %s does not exist", tableName)
    }

    if len(table.Columns) != len(data) {
        return fmt.Errorf("data length does not match number of columns")
    }

    table.mu.Lock()
    defer table.mu.Unlock()

    table.Rows = append(table.Rows, data)
    return nil
}

// Save saves the database to a directory and creates a CSV file for each table
func (db *Database) Save() error {
    db.mu.Lock()
    defer db.mu.Unlock()

    // Create the directory for the database if it doesn't exist
    err := os.MkdirAll(db.Name, os.ModePerm)
    if err != nil {
        return fmt.Errorf("failed to create directory for database: %v", err)
    }

    // Iterate over each table and create a CSV file for it
    for tableName, table := range db.Tables {
        // Create a CSV file for the table
        file, err := os.Create(fmt.Sprintf("%s/%s.csv", db.Name, tableName))
        if err != nil {
            return fmt.Errorf("failed to create CSV file for table %s: %v", tableName, err)
        }

        // Write the data to the CSV file
        writer := csv.NewWriter(file)

        // Write the column headers
        if err := writer.Write(table.Columns); err != nil {
            file.Close()
            return fmt.Errorf("failed to write column headers for table %s: %v", tableName, err)
        }

        // Write the rows
        for _, row := range table.Rows {
            if err := writer.Write(row); err != nil {
                file.Close()
                return fmt.Errorf("failed to write data to table %s: %v", tableName, err)
            }
        }

        writer.Flush()
        file.Close()

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
