package services

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type RMS struct {
	DeviceId int     `json:"device_id"`
	Time     int     `json:"time"`
	Acc_x    float64 `json:"acc_x_rms"`
	Acc_y    float64 `json:"acc_y_rms"`
	Acc_z    float64 `json:"acc_z_rms"`
}

type Move struct {
	DeviceId int `json:"device_id"`
	Time     int `json:"time"`
	Move     int `json:"move"`
}

func AddMoveDB(move Move, svc *dynamodb.DynamoDB) {
	av, err := dynamodbattribute.MarshalMap(move)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("tinyml-move"),
	}
	if err != nil {
		fmt.Println("Error with marshalling.")
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}

func AddRMSDB(rms RMS, svc *dynamodb.DynamoDB) {
	av, err := dynamodbattribute.MarshalMap(rms)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("tinyml-rms"),
	}
	if err != nil {
		fmt.Println("Error with marshalling.")
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}
