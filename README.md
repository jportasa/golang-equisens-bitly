# Bitly link clicks pull

The script connects to Bitly API with a token anf gets all link_id's and next for each one gets the clicks done in the current mont.
Afterwards, it connects to a Mysql and inserts the data: link_id, period, total month clicks.

The script can run in an AWS Lambda triggered by Cloudwatch events.

Env vars: BITLYTOKEN, DBHOST, DBUSER, DBPASSWORD, DBNAME


## Upload the binary to Lambda
````
GOOS=linux GOARCH=amd64 go build -o main main.go
zip -9 main.zip main
```