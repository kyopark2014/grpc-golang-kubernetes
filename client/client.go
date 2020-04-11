package main

import (
	"fmt"
	"grpc-golang-client/log"
	"grpc-golang-client/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	log.I("Start the client...")

	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

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
}
