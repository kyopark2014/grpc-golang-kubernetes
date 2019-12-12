# grpc-golang-kubernetes

This project shows how to create gRPC between two appliclations in kubernetes.

### Create EKS cluster
$ eksctl create cluster -f k8s/cluster-grpc-golang-kubernetes.yaml

### Prepare protobuf
```c
$ git clone https://github.com/tensor-programming/grpc_tutorial
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.1/protoc-3.11.1-linux-x86_64.zip
$ unzip protoc-3.11.1-linux-x86_64.zip
$ sudo cp bin/protoc /usr/bin/
$ sudo cp -R include/* /usr/local/include/
$ sudo chmod 777 /usr/bin/protoc
$ sudo add-apt-repository ppa:longsleep/golang-backports
$ sudo apt-get update
$ sudo apt-get install golang-go
$ go get -u google.golang.org/grpc
$ go get -u github.com/golang/protobuf/protoc-gen-go  
$ go get github.com/gin-gonic/gin
$ protoc --proto_path=proto --proto_path=/usr/local/include/ --go_out=plugins=grpc:proto proto/service.proto
```

[service.proto]
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
FROM golang:1.13.0
WORKDIR /usr/src/app
RUN apt-get update && apt-get install -y unzip
RUN go get -u google.golang.org/grpc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.1/protoc-3.11.1-linux-x86_64.zip
RUN unzip protoc-3.11.1-linux-x86_64.zip
RUN rm protoc-3.11.1-linux-x86_64.zip
RUN mv /usr/src/app/bin/protoc /usr/bin/
RUN cp -R include/* /usr/local/include/
RUN chmod 777 /usr/bin/protoc
RUN go get github.com/gin-gonic/gin
RUN go get -u github.com/golang/protobuf/protoc-gen-go
# RUN protoc --proto_path=proto --proto_path=/usr/local/include/ --go_out=plugins=grpc:proto proto/service.proto
COPY . .
EXPOSE 4040
CMD ["go","run", "main.go"]
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
FROM golang:1.13.0
WORKDIR /usr/src/app
RUN apt-get update && apt-get install -y unzip
RUN go get -u google.golang.org/grpc
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.1/protoc-3.11.1-linux-x86_64.zip
RUN unzip protoc-3.11.1-linux-x86_64.zip
RUN rm protoc-3.11.1-linux-x86_64.zip
RUN mv /usr/src/app/bin/protoc /usr/bin/
RUN cp -R include/* /usr/local/include/
RUN chmod 777 /usr/bin/protoc
RUN go get github.com/gin-gonic/gin
RUN go get -u github.com/golang/protobuf/protoc-gen-go
# RUN protoc --proto_path=proto --proto_path=/usr/local/include/ --go_out=plugins=grpc:proto proto/service.proto
COPY . .
EXPOSE 8080
CMD ["go","run", "main.go"]
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
