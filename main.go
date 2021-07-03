package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	dynamoUtil "github.com/tkeyo/tinyml-be/services"
	util "github.com/tkeyo/tinyml-be/util"
)

var dynamo *dynamodb.DynamoDB // DynamoDB instance
var APIAuthKey string         // API auth key

// Connects to DynamoDB
func connectDynamoDB() (db *dynamodb.DynamoDB) {
	creds := credentials.NewEnvCredentials()
	creds.Get()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: creds,
	})
	if err != nil {
		fmt.Println(err)
	}
	svc := dynamodb.New(sess)
	return svc
}

// Endpoint to check if service is up and running
func healthCheck(c *gin.Context) {
	fmt.Println("Health check request")
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Server ON",
		})
	}
}

// Endpoint to add RMS values to DB
func endpointRMS(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		var rms dynamoUtil.RMS
		c.BindJSON(&rms)

		go dynamoUtil.AddRMSDB(rms, dynamo)
		c.JSON(202, gin.H{
			"message": "Accepted",
		})
	}
}

// Endpoint to get move data
func getMoveData(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		minTimeSet := util.GetCurrentTime() - (3600 * 1000)
		deviceIdSet := 1

		timestamps, x, y, circle := dynamoUtil.ScanMoveDB(minTimeSet, deviceIdSet, dynamo)
		c.JSON(200, gin.H{
			"move_x":     x,
			"move_y":     y,
			"circle":     circle,
			"timestamps": timestamps,
		})
	}
}

// Endpoint to get RMS data
func getRMSData(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		minTimeSet := util.GetCurrentTime() - (3600 * 1000)
		deviceIdSet := 1

		timestamps, accXRMS, accYRMS, accZRMS := dynamoUtil.ScanRMSDB(minTimeSet, deviceIdSet, dynamo)
		c.JSON(200, gin.H{
			"timestamps": timestamps,
			"acc_x_rms":  accXRMS,
			"acc_y_rms":  accYRMS,
			"acc_z_rms":  accZRMS,
		})
	}
}

// Endpoint to add move data values to DB
func endpointMove(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		var move dynamoUtil.Move
		c.BindJSON(&move)

		go dynamoUtil.AddMoveDB(move, dynamo)
		c.JSON(202, gin.H{
			"message": "Accepted",
		})
	}
}

func main() {
	// Run code
	// API_AUTH_KEY=123 AWS_ACCESS_KEY_ID=xyz AWS_SECRET_ACCESS_KEY=xyz GIN_MODE=release go run *.go
	fmt.Println("Server is running....")

	dynamo = connectDynamoDB()
	APIAuthKey = os.Getenv("API_AUTH_KEY")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.GET("/api/health", healthCheck)
	r.GET("/api/get-rms", getRMSData)
	r.GET("/api/get-move", getMoveData)
	r.POST("/api/write-rms", endpointRMS)
	r.POST("/api/write-move", endpointMove)
	r.Run(":8081")
}
