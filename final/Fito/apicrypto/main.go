package main 
import (
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "apicrypto/repository"
  "encoding/json"
  "net/http"
)
func router(req events.APIGatewayProxyRequest) (ret events.APIGatewayProxyResponse, err error) {
  query:=req.QueryStringParameters
  
  aux,f:= query["use_mx"]
  use_mx := f  && aux=="true" 
  aux,f = query["use_hk"]
  use_hk := f && aux=="true"
  aux,f = query["use_usd"]
  use_usd := f  && aux=="true"
  
  repo,err:=repository.NewFilterRepository() 
  if err != nil{
    return
  }
  l,err:=repo.List(query["from"],query["to"],use_mx,use_hk,use_usd)
  if err !=nil {
    return events.APIGatewayProxyResponse{
      StatusCode : http.StatusBadRequest,
      Body : err.Error(),
    },nil
  }
  body,_ :=json.Marshal(l) 
  return events.APIGatewayProxyResponse{
    StatusCode : http.StatusOK,
    Body : string(body),
  },nil
}
func main(){
  lambda.Start(router)
}
