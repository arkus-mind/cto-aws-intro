package repository

import (
	"database/sql"
	"todolist/config"
        "errors"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
) 

type TaskStatus int 

type TodoRepository struct {
  Connection *sql.DB
}
type TodoItem struct {
  Id    string
  Title string
  Description string
  Created_at string
  Updated_at string 
  Status TaskStatus

}

const (
  PendingToWork TaskStatus = iota
  InProgress
  Done
)

func NewTodoRepository (conf * config.WholeConfig) (ret *TodoRepository , err error) {
  ret = &TodoRepository{}
  ret.Connection,err = sql.Open("postgres",conf.ConnectionString())
  return
}

func (repo *TodoRepository)  List (title string) (ret []*TodoItem, err error){
  title = "%"+title+"%"
  rows,err:=repo.Connection.Query("SELECT id::text,title::text,description::text "+
                                  ",Created_at::text,updated_at::text,status "+
                                  "FROM todolist.todolist WHERE title LIKE $1",&title)
  if err != nil{
    return
  }
  for rows.Next(){
    nItem :=new(TodoItem)
    rows.Scan(&nItem.Id,&nItem.Title,&nItem.Description,&nItem.Created_at,
              &nItem.Updated_at,&nItem.Status)
    ret=append(ret,nItem)
  }
  return
}
  
func (repo * TodoRepository) NewTask(title string , description string) (id string , err error) {
  id =uuid.New().String()
  if len(title) ==0 {
    err = errors.New("Invalid title")
    return
  }
  _,err=repo.Connection.Exec("INSERT INTO todolist.todolist (id,title,description)  VALUES($1,$2,$3)",
                       &id,&title,&description); 
  return  
}

func (repo * TodoRepository) UpdateTitle(id string , title string) error{
  if len(title) ==0 {
    return errors.New("Invalid title")
  }
  repo.Connection.Exec("UPDATE todolist.todolist SET updated_at=NOW(),title=$1 WHERE id=$2",&title,&id)  
  return nil
}

func (repo * TodoRepository) UpdateStatus(id string , status TaskStatus) error{

  if status > Done || status < PendingToWork {
    return errors.New("Invalid status")
  }
  _,err:=repo.Connection.Exec("UPDATE todolist.todolist SET updated_at=NOW(),status=$1 WHERE id=$2",&status,&id)  
  return err
}

func (repo * TodoRepository) UpdateDescription(id string , description string) error{
  if len(description) == 0 {
    return errors.New("Invalid description");
  }
  _,err:=repo.Connection.Exec("UPDATE todolist.todolist SET updated_at=NOW(),description=$1 WHERE id=$2",&description,&id)  
  return err
}

func (repo * TodoRepository) Delete(id string ) error{
  _,err:=repo.Connection.Exec("DELETE FROM todolist.todolist WHERE id=$1",id)
  return err
}
