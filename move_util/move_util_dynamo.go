package move_util

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
