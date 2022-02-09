package main

import (
  "log"
  "net/http"
  "encoding/json"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/google/uuid"
  "strconv"
)

type Item struct {
  Id string `json:"id"`
  Memecoins float64
}

func main(){
  resp,err :=http.Get("https://api.bitso.com/v3/ticker/?book=btc_mxn")
  if err != nil{
    panic(err)
  }

  decoder := json.NewDecoder(resp.Body) 
  var datum map[string] map[string]string
  decoder.Decode(&datum)
  
  sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
      }))
  
  svc := dynamodb.New(sess)
  btc ,err:= strconv.ParseFloat(datum["payload"]["last"],64)
  
  var item Item  

  result, err := svc.GetItem(&dynamodb.GetItemInput{
    TableName: aws.String("LeMemeCoins"),
    Key: map[string]*dynamodb.AttributeValue{
        "id": {
            S: aws.String("lastValue"),
        },
    },
  })
  if err != nil{
    panic(err.Error())
  }
  
  err = dynamodbattribute.UnmarshalMap(result.Item,&item)
   
  if err != nil{
    panic(err.Error())
  }
  log.Printf("%f VS %f", item.Memecoins, btc) 
  if result.Item !=nil && item.Memecoins == btc{
    log.Printf("value already exists")
    return 
  }

  item = Item{uuid.New().String(),btc}
  av, err := dynamodbattribute.MarshalMap(item)
  if err != nil {
    log.Fatalf("Got error marshalling%s", err)
  }
  
  input := &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String("LeMemeCoins"),
  }
  _, err = svc.PutItem(input)
  if err != nil {
    log.Fatalf("Got error calling PutItem: %s", err)
  }
  
  item.Id = "lastValue"
  av, err = dynamodbattribute.MarshalMap(item)
  if err != nil {
    log.Fatalf("Got error marshalling%s", err)
  }
  input = &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String("LeMemeCoins"),
  }
  _, err = svc.PutItem(input)

  if err != nil {
    log.Fatalf("Got error calling PutItem: %s", err)
  }
  
  log.Printf("done")
}
