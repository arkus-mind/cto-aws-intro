package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todolist/config"
	"todolist/repository"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"errors"
)

func genericError(errlegend string) events.APIGatewayProxyResponse{
  return events.APIGatewayProxyResponse{
    StatusCode : http.StatusInternalServerError,
    Body : errlegend,
  }
}

func successResponse(legend string) events.APIGatewayProxyResponse{
 return events.APIGatewayProxyResponse{
    StatusCode : http.StatusOK,
    Body : legend,
  }
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  
  var data map[string]string
  var err error
  if req.RequestContext.HTTPMethod == "POST"{
    err = json.Unmarshal([]byte(req.Body),&data)
  }else{
    data = req.QueryStringParameters
  }

  ret := "OK"

  if err != nil{
    log.Print(err.Error())
    return genericError(err.Error()+ "\n"+req.Body),nil
  }
  config,err := config.Parse("config_json/development_config.json")
  
  if err != nil{
    return genericError(err.Error()),nil
  }

  repo,err := repository.NewTodoRepository(config) 
  
  if err != nil {
    return genericError(err.Error()),nil
  }
  switch data["command"]{
    case "new":
      uid := ""
      uid,err = repo.NewTask(data["title"],data["description"])
      ret = uid
    case "update-title":
      err = repo.UpdateTitle(data["uuid"],data["title"])
    case "update-description":
      err = repo.UpdateDescription(data["uuid"],data["description"])
    case "update-status":
      i,_:=strconv.Atoi(data["status"])
      err = repo.UpdateStatus(data["uuid"], repository.TaskStatus(i))
    case "delete":
      err = repo.Delete(data["uuid"])
    case "search":
      terms,err := repo.List(data["term"])
      if err == nil{
        json,jerr := json.Marshal(terms)
        ret = string(json)
        if jerr != nil{
          err = jerr
        }
      }
    default:
      err = errors.New("Invalid command")
  }
  if err != nil{
    return genericError(err.Error()),nil
  }
  return successResponse(ret),nil
}
func main() {
  lambda.Start(router)
}
