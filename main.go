package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	dynamoUtil "github.com/tkeyo/tinyml-be/services"
	util "github.com/tkeyo/tinyml-be/util"
)

var dynamo *dynamodb.DynamoDB
var APIAuthKey string

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

func endpointRMS(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		var rms dynamoUtil.RMS
		c.BindJSON(&rms)

		dynamoUtil.AddRMSDB(rms, dynamo)
		c.JSON(200, gin.H{
			"message": "OK",
		})
	}
}

func getMoveData(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		minTimeSet := 1618225200 // 12.4.2021 13:00:00
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

func getRMSData(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		minTimeSet := 1618225200 // 12.4.2021 13:00:00
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

func endpointMove(c *gin.Context) {
	requestAuthKey := c.Request.Header["Authorization"][0]

	if !util.IsAuthorized(requestAuthKey, APIAuthKey) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	} else {
		var move dynamoUtil.Move
		c.BindJSON(&move)

		dynamoUtil.AddMoveDB(move, dynamo)
		c.JSON(200, gin.H{
			"message": "OK",
		})
	}
}

func main() {
	fmt.Println("Server is running....")

	dynamo = connectDynamoDB()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	APIAuthKey = os.Getenv("API_AUTH_KEY")

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
	r.GET("/api/get-rms", getRMSData)
	r.GET("/api/get-move", getMoveData)
	r.POST("/api/write-rms", endpointRMS)
	r.POST("/api/write-move", endpointMove)
	r.Run(":8081")
}
