deploy:
	go build 
	zip package.zip -r todolist config_json
	aws lambda update-function-code  --function-name ctotodo --zip-file fileb://package.zip
