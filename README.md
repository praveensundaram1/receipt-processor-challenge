# Receipt Processor

This repository contains my solution to the Receipt Processor Challenge. It is a webserver developed in Go, and manages receipt processing through two primary endpoints:

1. The /receipts/process POST endpoint accepts a Receipt JSON object, computes the points based on predefined rules, and returns a JSON response with the processed receipt's ID.
2. The /receipts/{id}/points GET endpoint allows retrieval of the calculated points for a receipt, identified by the receipt's ID provided in the POST endpoint's response. This endpoint delivers a JSON object with the total points that were computed for the submitted receipt.

To see a more detailed description of the requirements, click [here](https://github.com/fetch-rewards/receipt-processor-challenge "Receipt Processor Challenge Requirements")

Note: To resolve any issues, email me at praveen.sundaram.2022@gmail.com

##Project Setup

The project can be run using either Go directly or within a Docker container. After cloning the repository, you must `cd` to the `receipt-processor-challenge` directory. 


###Setup Instructions to Run Application on Local Machine

1. Run setup commands

```bash
go mod tidy
go mod vendor
```

2. Run the following command on the CLI to start the application

```bash
go run cmd/receiptprocessor/main.go
```

Once the server starts, it will listen on port 8080 and be ready to handle requests. Look for the message "Listening on localhost:8080" in the terminal. To stop the server, press `ctrl-c`.

###Setup Instructions to Run Application on Docker Container

1. Build the Docker container using the following command, tagging it as fetch-receipt-service:

```bash
docker build . -t fetch-receipt-service
```

2.  Run the Docker container with the command:
```bash
docker run -it -p 8080:8080 fetch-receipt-service
```
Once the server starts, it will listen on port 8080 and be ready to handle requests. Look for the message "Listening on localhost:8080" in the terminal. To stop the server, press `ctrl-c`.

###Server Interaction and API Usage

The server is designed to handle two specific types of HTTP requests, as described in the project specifications. To interact with the server, you can use applications like Postman or the command-line `curl` command on UNIX-based systems (including macOS and Linux). Below are examples of `curl` commands that demonstrates how to send requests to the server.  Any JSON that follows the `api.yml` file specifications will work.

1. Endpoint: /receipts/process
   Method: POST
   Description: This endpoint is responsible for accepting and processing receipt data.

Here are 2 examples

```bash
curl --location 'http://localhost:8080/receipts/process' \
--header 'Content-Type: application/json' \
--data '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}'
```

You should get back a JSON object that looks like this:

```json
{ "points" : "28" }
```

```bash
curl --location 'http://localhost:8080/receipts/process' \
--header 'Content-Type: application/json' \
--data '{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}'
```

You should get back a JSON object that looks like this:

```json
{ "points" : "109" }
```

Note: A detailed description of how these points are calculated can be found [here](https://github.com/fetch-rewards/receipt-processor-challenge#rules)

2. Endpoint: /receipts/{id}/points
   Method: GET
   Description: Retrieves the total points for a receipt. Substitute {id} with the receipt ID obtained from the POST endpoint.

```bash
curl --location 'http://localhost:8080/receipts/{id}/points'
```

###Running Tests

To run unit tests, run the following command:

```bash
go test ./...
```
