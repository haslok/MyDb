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
    "github.com/haslok/MyDb"
    "fmt"
)

func main() {
    db := MyDb.NewDatabase("exampleDB")
    err := db.CreateTable("users", []string{"id", "name", "email"})
    if err != nil {
        fmt.Println(err)
        return
    }

    err = db.InsertInto("users", []string{"1", "John Doe", "john@example.com"})
    if err != nil {
        fmt.Println(err)
        return
    }

    err = db.Save()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Database saved successfully")
}
```
