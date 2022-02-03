package todolist

import (
	"testing"
	"todolist/config"
        "todolist/repository"
)

func TestTodoRepository(t * testing.T){
  conf,err := config.Parse("../config_json/development_config.json")
  if err != nil{
    t.Fatal(err.Error())
  }
  todo,err :=  repository.NewTodoRepository(conf)
  if err != nil{
    t.Fatal(err.Error())
  }
  uuid,err:=todo.NewTask("Create robowaifus", "create anime based robots")

  if err != nil{
    t.Fatal(err.Error())
  }
  err =todo.UpdateDescription(uuid,"Create kawai anime based robots");
  if err !=nil{
    t.Fatal(err.Error())
  }
  err = todo.UpdateStatus(uuid,repository.Done)
  if err != nil{
    t.Fatal(err.Error())
  }
  items,err:=todo.List("robo")
  if err != nil{
    t.Fatal(err.Error())
  }
  t.Logf("got %d items",len(items))
}
