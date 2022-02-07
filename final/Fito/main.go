package testing

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	_ "github.com/lib/pq"
	"github.com/urakozz/go-dynamodb-stream-subscriber/stream"
)

func main() {
  cfg := aws.NewConfig().WithRegion("us-east-1")
  sess := session.New()
  streamSvc := dynamodbstreams.New(sess, cfg)
  dynamoSvc := dynamodb.New(sess, cfg)
  table := "LeMemeCoins"
  
  result, err := streamSvc.ListStreams(&dynamodbstreams.ListStreamsInput{ })
  fmt.Printf(result.String())
  

  return 
  streamSubscriber := stream.NewStreamSubscriber(dynamoSvc, streamSvc, table)
  streamSubscriber.SetLimit(1)
  ch,_ := streamSubscriber.GetStreamDataAsync()

  fmt.Printf("starting to hear channel\n")
  record :=  <-ch
  re := regexp.MustCompile(`N: "(.*)"`)
  btcValueString := re.FindAllStringSubmatch(record.String(), -1)[0][1]
  mxValue,_ := strconv.ParseFloat(btcValueString,64)   
  
  resp,err :=http.Get("https://api.bitso.com/v3/ticker/?book=btc_usd")
  if err != nil{
    panic(err)
  }
  decoder := json.NewDecoder(resp.Body) 
  var datum map[string] map[string]string
  decoder.Decode(&datum)
  usdValue ,err:= strconv.ParseFloat(datum["payload"]["last"],64) 
  db,err:= sql.Open("postgres", 
                  "host=ctolearn.cluster-ro-cjincqaxcmb8.us-east-1.rds.amazonaws.com "+
                  "port=5432 user=postgres password=hdAXa4yVe7HWRXb dbname=postgres sslmode=disable")
  if err != nil{
    panic(err)
  }
  defer db.Close();

  usd2mx := usdValue/mxValue
  _,err=db.Exec("INSERT INTO btc_loco.btc_loco (mexican_peso,hk_dollar,usd_dollar) VALUES($1,$2,$3)",
          &mxValue,&usd2mx,&usd2mx)
 
  if err != nil{
    panic(err)
  }

}
