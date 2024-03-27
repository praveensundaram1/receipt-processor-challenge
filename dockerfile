FROM golang:1.22.1

WORKDIR /app

COPY . /app

RUN apt-get update && go mod tidy && go mod vendor
EXPOSE 8080
CMD ["go", "run", "cmd/receiptprocessor/main.go"]



