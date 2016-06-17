package controllers

import "github.com/revel/revel"
import "github.com/ZenClark/reveltest/app"
import "fmt"

type App struct {
	*revel.Controller
}

type Message struct {
  ID int64
  Msg string
}

func (c App) Index() revel.Result {
  var IsError bool = false
  sql := "FindMessages 0;"
  
  stmt, statementError := app.DB.Prepare(sql)
  
  defer stmt.Close()
  
  if statementError != nil {
    Error:=statementError.Error()
    IsError = true
    fmt.Printf("The error is: %s\n", Error)
    return c.Render(Error, IsError)
  }
  
  rows, queryError := stmt.Query()
  
  if queryError != nil {
    Error:=queryError.Error()
    IsError = true
    return c.Render(Error, IsError)
  }
  
  var ID int64
  var message string
  
  rows.Next()
  scanError := rows.Scan(&ID, &message)
  if scanError != nil {
    Error:=scanError.Error()
    IsError = true
    return c.Render(Error, IsError)
  }
  
  var Messages = make([]Message, 0, 10)
  
  Messages = append(Messages, Message{0, "This is the first message!"})
  Messages = append(Messages, Message{ID, message})
  for rows.Next() == true {
    rows.Scan(&ID, &message)
    Messages = append(Messages, Message{ID, message})
  }
  
  
	return c.Render(Messages, IsError)
}

func (c App) NewMessage(message string) revel.Result {

  _, err := app.DB.Exec("CreateMessage ?;", message)
  
  if err != nil {
   fmt.Printf("Error: %s\n", err.Error()) 
  }
  
  return c.Redirect("/")
  
}
