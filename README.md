# MyDb

MyDb is a simple database management system written in Go.

## Installation

To install MyDb, use the following command:

```sh
go mod init MyDb
go mod tidy
go get github.com/haslok/MyDb
```
## Usage 
Here the example for using MyDb :
```go
package main

import (
	"fmt"
	"github.com/haslok/MyDb" // Import the MyDb package
)

func main() {
	// Create a new database
	db := MyDb.NewDatabase("MyDB")

	// Create a table
	err := db.CreateTable("Users", []string{"ID", "Name", "Age"})
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// Insert data into the table
	err = db.InsertInto("Users", []string{"1", "Alice", "25"})
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}

	err = db.InsertInto("Users", []string{"2", "Bob", "30"})
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}

	// Update data in the table
	err = db.UpdateData("Users", func(row []string) bool {
		return row[0] == "1" // Condition: Update the row where ID is "1"
	}, []string{"1", "Alice", "26"}) // New data
	if err != nil {
		fmt.Println("Error updating data:", err)
		return
	}

	// Save the database to disk
	err = db.Save()
	if err != nil {
		fmt.Println("Error saving database:", err)
		return
	}

	fmt.Println("Database operations completed successfully!")
}
```
