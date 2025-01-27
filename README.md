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
    "github.com/haslok/MyDb"
)

func main() {
    // Create a new database
    db := MyDb.NewDatabase("example_db")

    // Execute a CREATE TABLE command
    err := db.Command("create table users has id, name, email")
    if err != nil {
        fmt.Println("Error creating table:", err)
        return
    }

    // Execute an INSERT command
    err = db.Command("insert to users 1, John Doe, john@example.com")
    if err != nil {
        fmt.Println("Error inserting row:", err)
        return
    }

    // Execute a GET command
    err = db.Command("get from users where id=1")
    if err != nil {
        fmt.Println("Error getting rows:", err)
        return
    }

    // Execute an UPDATE command
    err = db.Command("update users set email=john.doe@example.com where id=1")
    if err != nil {
        fmt.Println("Error updating row:", err)
        return
    }

    // Execute a DELETE command
    err = db.Command("delete from users where id=1")
    if err != nil {
        fmt.Println("Error deleting row:", err)
        return
    }

    db.Save()
    fmt.Println("Commands executed successfully")
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
