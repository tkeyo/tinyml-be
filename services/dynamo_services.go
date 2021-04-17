package services

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type RMS struct {
	DeviceId int     `json:"device_id"`
	Time     int     `json:"time"`
	Acc_x    float64 `json:"acc_x_rms"`
	Acc_y    float64 `json:"acc_y_rms"`
	Acc_z    float64 `json:"acc_z_rms"`
}

type Move struct {
	DeviceId  int `json:"device_id"`
	Timestamp int `json:"time"`
	Move      int `json:"move"`
}

type M map[string]interface{}

func ScanMoveDB(minTime int, deviceId int, svc *dynamodb.DynamoDB) ([]M, []M, []M, []int) {
	filt := expression.Name("device_id").Equal(expression.Value(deviceId)).And(expression.Name("time").GreaterThan(expression.Value(minTime)))
	proj := expression.NamesList(expression.Name("device_id"), expression.Name("time"), expression.Name("move"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("tinyml-move"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Got error retrieving data: %s", err)
	}

	var xMovementSlice []M
	var yMovementSlice []M
	var circleMovementSlice []M
	var timestamps []int

	for _, i := range result.Items {
		move := Move{}
		err := dynamodbattribute.UnmarshalMap(i, &move)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		timestamps = append(timestamps, move.Timestamp*1000)

		switch move.Move {
		case 1:
			xMovementSlice = append(xMovementSlice, M{"x": move.Timestamp * 1000, "y": 1})
		case 2:
			yMovementSlice = append(yMovementSlice, M{"x": move.Timestamp * 1000, "y": 2})
		case 3:
			circleMovementSlice = append(circleMovementSlice, M{"x": move.Timestamp * 1000, "y": 3})
		}
	}
	// fmt.Println(XMovementSlice)
	// fmt.Println(YMovementSlice)
	// fmt.Println(CircleMovementSlice)
	return xMovementSlice, yMovementSlice, circleMovementSlice, timestamps
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
