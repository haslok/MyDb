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
	db := NewDatabase("my_database") // Create a new database

	// 1. Create a table
	err := db.Command("CREATE TABLE users HAS name, age, city")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Table 'users' created successfully!")

	// 2. Insert data into the table
	err = db.Command("INSERT TO users ahmad, 23, cairo")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = db.Command("INSERT TO users lila, 30, alexandria")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data inserted successfully!")

	// 3. Get data from the table
	err = db.Command("GET FROM users WHERE age=23")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 4. Update data in the table
	err = db.Command("UPDATE users SET city=giza WHERE name=ahmad")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data updated successfully!")

	// 5. Delete data from the table
	err = db.Command("DELETE FROM users WHERE name=lila")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data deleted successfully!")

	// 6. Save the database to disk
	err = db.Save()
	if err != nil {
		fmt.Println("Error saving the database:", err)
		return
	}
	fmt.Println("Database saved successfully!")
}

```
## Result 
1. the files will be :
   ```
    MyDb (your database name /// folder)
     |_ Users.csv (your table name /// .csv file)
     ...
   ```
2. The test code result will be :
   ```
      Table 'users' created successfully!
       Data inserted successfully!
       Results: [map[age:23 city:cairo name:ahmad]]
       Data updated successfully!
       Data deleted successfully!
       Database saved successfully!
3. the users.csv file will be :
   ```
   id,name,age
   1,Alice,30
   2,Bob,26
