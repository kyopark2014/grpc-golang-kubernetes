###################
##  build stage  ##
###################
FROM golang:1.13.0-alpine as builder
WORKDIR /grpc-golang-server
COPY . .
RUN go build -v -o grpc-golang-server

##################
##  exec stage  ##grpc-golang-server
##################
FROM alpine:3.10.2
WORKDIR /app
COPY --from=builder /grpc-golang-server /app/
CMD ["./grpc-golang-server"]
