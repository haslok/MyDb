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
 // Import the MyDb package
package main

import (
	"fmt"
	"log"
	"github.com/haslok/MyDb"
)

func main() {
	// Step 1: Create a new database
	db := MyDb.NewDatabase("MyTestDB")

	// Step 2: Create a new table
	err := db.CreateTable("users", []string{"id", "name", "age"})
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	// Step 3: Insert data into the "users" table
	err = db.InsertInto("users", []string{"1", "Alice", "30"})
	if err != nil {
		log.Fatal("Error inserting row:", err)
	}

	err = db.InsertInto("users", []string{"2", "Bob", "25"})
	if err != nil {
		log.Fatal("Error inserting row:", err)
	}

	// Step 4: Update data in the "users" table (for example, change Bob's age)
	err = db.UpdateData("users", func(row []string) bool {
		return row[1] == "Bob"  // condition: row's name is "Bob"
	}, []string{"2", "Bob", "26"})  // new data for Bob
	if err != nil {
		log.Fatal("Error updating data:", err)
	}

	// Step 5: Save the database (it will create a directory and CSV files)
	err = db.Save()
	if err != nil {
		log.Fatal("Error saving database:", err)
	}

	// Step 6: Load a table from the saved CSV file
	table, err := db.SelectTable("users")
	if err != nil {
		log.Fatal("Error selecting table:", err)
	}

	// Print the table's data
	fmt.Println("Table:\n", table.Columns)
	for _, row := range table.Rows {
		fmt.Println(row)
	}
}

```
## Result 
1. the files will be :
   ```
    MyDb (your database name /// folder)
     |_ Users.csv (your table name /// .csv file)
     ...
   ```
2. The Users.csv will be :
   ```
   ID,Name,Age
   1,Alice,26
   2,Bob,30
