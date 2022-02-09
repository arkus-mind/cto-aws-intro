package main

import (
  "context"
  "fmt"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "database/sql"
_ "github.com/lib/pq"

)

func handleRequest(ctx context.Context, e events.DynamoDBEvent) {
  db,err:= sql.Open("postgres", 
  "host=ctolearn.cluster-cjincqaxcmb8.us-east-1.rds.amazonaws.com "+
  "port=5432 user=postgres password=hdAXa4yVe7HWRXb dbname=postgres sslmode=disable")
  if err != nil{
    fmt.Printf("something wrong "+err.Error())
    panic(err)
  }
  defer db.Close()

  for _, record := range e.Records {
    fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
    key:=""
    mxValue:=0.0
    for _, value := range record.Change.NewImage {
      if value.DataType() == events.DataTypeString{
        key = value.String()
      }
      if value.DataType() == events.DataTypeNumber{
        mxValue,_ = value.Float()
      }
    }
    if (key!= "lastValue"){
      usdValue  := mxValue/20.68
      hkValue   := mxValue/2.66
      _,err=db.Exec("INSERT INTO btc_loco.btc_loco (mexican_peso,hk_dollar,usd_dollar) VALUES($1,$2,$3)",
      &mxValue,&usdValue,&hkValue)
      if err != nil{
        panic(err)
      }
    }
  }
}
func main(){
  lambda.Start(handleRequest)
}
