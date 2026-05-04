package main

import (
  "fmt"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
  "github.com/homemusic/backend/internal/model"
)

func main() {
  db, err := gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{})
  if err != nil { panic(err) }
  err = db.AutoMigrate(&model.User{})
  if err != nil { panic(err) }
  fmt.Println("done")
}
