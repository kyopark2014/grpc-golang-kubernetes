###################
##  build stage  ##
###################
FROM golang:1.13.0-alpine as builder
WORKDIR /grpc-golang-client
COPY . .
RUN go build -v -o grpc-golang-client

##################
##  exec stage  ##grpc-golang-client
##################
FROM alpine:3.10.2
WORKDIR /app
COPY --from=builder /grpc-golang-client /app/
CMD ["./grpc-golang-client"]
