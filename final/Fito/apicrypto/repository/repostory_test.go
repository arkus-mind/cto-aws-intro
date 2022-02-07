package repository

import(
  "testing"
)

func TestCryptoFilter(t * testing.T){
  repo,err := NewFilterRepository()
  if err != nil{
    t.Fatal(err.Error())
  }
  l,err:=repo.List("2022-01-01 00:00:00", "2022-03-01 00:00:00", true,true,true) 
  if err != nil{
    t.Fatal(err.Error())
  }
  t.Log(l)
  
  l,err =repo.List("2022-01-01 00:00:00", "2022-03-01 00:00:00", true,false,true) 
  if err != nil{
    t.Fatal(err.Error())
  }
  t.Log(l)

}
