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
