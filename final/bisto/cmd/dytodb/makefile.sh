#!/bin/zsh

case $1 in
deploy)
    GOOS=linux GOARCH=amd64 go build -o main main.go
    zip main.zip main
    aws lambda update-function-code --function-name DynamoToRDS --zip-file fileb://main.zip
    ;;
clean)
    rm main main.zip
    ;;
test)
    go run main.go
    ;;
*)
    echo -n "unknown action\n"
    ;;
esac
