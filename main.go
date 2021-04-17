package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	DynamoUtil "github.com/tkeyo/tinyml-be/services"
)

// func getRMSData(c *gin.Context) {
// }

var dynamo *dynamodb.DynamoDB

func connectDynamoDB() (db *dynamodb.DynamoDB) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewSharedCredentials(".aws/credentials", "default"),
	})
	if err != nil {
		fmt.Println(err)
	}
	svc := dynamodb.New(sess)
	return svc
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server ON",
	})
}

func endpointRMS(c *gin.Context) {
	var rms DynamoUtil.RMS
	c.BindJSON(&rms)

	DynamoUtil.AddRMSDB(rms, dynamo)
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func getMoveData(c *gin.Context) {
	DynamoUtil.ScanMoveDB(15, 1, dynamo)
}

func endpointMove(c *gin.Context) {
	var move DynamoUtil.Move
	c.BindJSON(&move)

	DynamoUtil.AddMoveDB(move, dynamo)
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func main() {
	fmt.Println("Server is running....")

	dynamo = connectDynamoDB()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.GET("/api/health", healthCheck)
	// r.GET("/api/get-rms", getRMSData)
	r.GET("/api/get-move", getMoveData)
	r.POST("/api/write-rms", endpointRMS)
	r.POST("/api/write-move", endpointMove)
	r.Run(":8081")
}
