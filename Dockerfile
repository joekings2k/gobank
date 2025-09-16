FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app 
COPY . .
RUN go build -o main main.go 

EXPOSE 8080
CMD ["/app/main"]