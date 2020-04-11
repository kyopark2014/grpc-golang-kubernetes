# grpc-golang-kubernetes

This project shows how to create gRPC between two appliclations in kubernetes.


### GRPC server

#### Define listener

```go
if grpcErr != nil {
  log.E("Localisation function Failed to listen : %v", grpcErr)
  os.Exit(1)
}

srv := grpc.NewServer()
proto.RegisterAddServiceServer(srv, &server{})
reflection.Register(srv)

if e := srv.Serve(listener); e != nil {
  panic(e)
}
```

#### Define operations

```go
func (s *server) Add(ctx context.Context, request *proto.Request) (*proto.Response, error) {
  log.D("Server: add()...")
  a, b := request.GetA(), request.GetB()

	result := a + b

	log.D("%v + %v = %v", a, b, result)

	return &proto.Response{Result: result}, nil
}

func (s *server) Multiply(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	log.D("Server: multiply)...")
	a, b := request.GetA(), request.GetB()

	result := a * b

	log.D("%v x %v = %v", a, b, result)

	return &proto.Response{Result: result}, nil
}
```

### GRPC client

#### Define grpc client

```go
conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
if err != nil {
	panic(err)
}

client := proto.NewAddServiceClient(conn)
```  

#### Define operations

```go
g := gin.Default()
g.GET("/add/:a/:b", func(ctx *gin.Context) {
	log.I("Client: add()...")

	a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
		log.E("Invalid Parameter A: %v", err)
		return
	}

	b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
		log.E("Invalid Parameter B: %v", err)
		return
	}

	req := &proto.Request{A: int64(a), B: int64(b)}
	if response, err := client.Add(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.E("error: %v", err)
	}
})

g.GET("/mult/:a/:b", func(ctx *gin.Context) {
	log.I("Client: mult()...")
	a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
		log.E("Invalid Parameter A: %v", err)
		return
	}
	b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
		log.E("Invalid Parameter B: %v", err)
		return
	}
	req := &proto.Request{A: int64(a), B: int64(b)}

	if response, err := client.Multiply(ctx, req); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(response.Result),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.E("error: %v", err)
	}
})

if err := g.Run(":8080"); err != nil {
	log.E("Failed to run server: %v", err)
}
```



### Create EKS cluster
$ eksctl create cluster -f k8s/cluster-grpc-golang-kubernetes.yaml

### Prepare protobuf
```c
$ protoc --proto_path=proto --proto_path=/usr/local/include/ --go_out=plugins=grpc:proto server/proto/server.proto

$ protoc --proto_path=proto --proto_path=/usr/local/include/ --go_out=plugins=grpc:proto client/proto/client.proto

```

#### [server.proto]

```c
syntax = "proto3";

package proto;

message Request {
  int64 a = 1;
  int64 b = 2;
}

message Response { 
  int64 result = 1; 
}

service AddService {
  rpc Add(Request) returns (Response);
  rpc Multiply(Request) returns (Response);
}

```

#### [client.proto]

```c

syntax = "proto3";

package proto;

message Request {
  int64 a = 1;
  int64 b = 2;
}

message Response { 
  int64 result = 1; 
}

service AddService {
  rpc Add(Request) returns (Response);
  rpc Multiply(Request) returns (Response);
}
```

### Create grpc-server
Dockerfile

```c

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

```


### Build server image

```c
$ docker build -t grpc-golang-server:v1 .
```

### Tagging
```c
$ docker tag grpc-golang-server:v1 994942771862.dkr.ecr.eu-west-2.amazonaws.com/repository-grpc-server
```

### Create repository in ECR if required
```c
$ aws ecr create-repository --region eu-west-2 --repository-name repository-grpc-server
```

### Push the image to ECR
```c
$ docker push 994942771862.dkr.ecr.eu-west-2.amazonaws.com/repository-grpc-server
```

## Make deployment and server in the kubernetes cloud
```c
$ kubectl create -f grpc-golang-server-deployment.yaml 
deployment.apps/grpc-server created

$ kubectl create -f grpc-golang-server-service.yaml 
service/grpc-server created
```


### Check the earn URL
Note loadbalancer type is used for this sample in order to easy understanding

```c
$ kubectl get services
NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP                                                               PORT(S)                      AGE
grpc-server                     LoadBalancer   10.100.84.254    a70aa446b1cbd11eaaabc0a6b8228ff9-1967545675.eu-west-2.elb.amazonaws.com   4040:32169/TCP               79s
```

### Create grpc-server
Dockerfile
```c

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

```

### Build the image for grpc client
```c
$ docker build -t grpc-golang-client:v1 .

$ docker tag grpc-golang-client:v1 994942771862.dkr.ecr.eu-west-2.amazonaws.com/repository-grpc-client
```

### Upload the image in ECR
```c
$ aws ecr create-repository --region eu-west-2 --repository-name repository-grpc-client

$ docker push 994942771862.dkr.ecr.eu-west-2.amazonaws.com/repository-grpc-client
```

### Make deployment and service
```c
$ kubectl create -f grpc-golang-client-deployment.yaml 
deployment.apps/grpc-client created

$ kubectl create -f grpc-golang-client-service.yaml 
service/grpc-client created
```


### Check the client URL
```c
$ kubectl get services
NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP                                                               PORT(S)                      AGE
grpc-client                     LoadBalancer   10.100.66.58     ab46519f91cc111eaaabc0a6b8228ff9-1383798696.eu-west-2.elb.amazonaws.com   8080:31352/TCP               93s
grpc-server                     LoadBalancer   10.100.84.254    a70aa446b1cbd11eaaabc0a6b8228ff9-1967545675.eu-west-2.elb.amazonaws.com   4040:32169/TCP               32m
```

### Run the operation of grpc
```c
$ curl -i ab46519f91cc111eaaabc0a6b8228ff9-1383798696.eu-west-2.elb.amazonaws.com:8080/add/100/200
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 12 Dec 2019 13:46:18 GMT
Content-Length: 17

{"result":"300"}
```


#### Reference
https://www.youtube.com/watch?v=Y92WWaZJl24


#### Test logs

```c
$ go run server.go 

$ go run client.go 

$ curl -i localhost:8080/add/20/30
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 11 Apr 2020 01:58:21 GMT
Content-Length: 15

{"result":"50"}

$ curl -i localhost:8080/add/5/7
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 11 Apr 2020 01:58:30 GMT
Content-Length: 15

{"result":"12"}

$ curl -i localhost:8080/multi/5/3
HTTP/1.1 404 Not Found
Content-Type: text/plain
Date: Sat, 11 Apr 2020 01:58:49 GMT
Content-Length: 18

404 page not found

$ curl -i localhost:8080/mult/5/3
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 11 Apr 2020 01:58:57 GMT
Content-Length: 15

```
