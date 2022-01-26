payload='{
  "command"     : "new",
  "title"       : "new task",
  "description" : "this is a new task"
}'
uuid=$(curl --data "$payload" 'https://h53fah3ii9.execute-api.us-east-1.amazonaws.com/default')

payload="{
  \"command\" : \"update-title\",
  \"title\"     : \"this is the new title \",
  \"uuid\"      : \"$uuid\"
}"
curl --data "$payload" 'https://h53fah3ii9.execute-api.us-east-1.amazonaws.com/default'

payload="{
  \"command\" : \"update-status\",
  \"uuid\"    : \"$uuid\",
  \"status\"  : \"2\"
}"
curl --data "$payload" 'https://h53fah3ii9.execute-api.us-east-1.amazonaws.com/default'

payload="{
  \"command\" : \"search\"
}"
curl --data "$payload" 'https://h53fah3ii9.execute-api.us-east-1.amazonaws.com/default'

