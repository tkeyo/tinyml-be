package services

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Root mean square (RMS) struct
type RMS struct {
	DeviceId  int     `json:"device_id"`
	Timestamp int64   `json:"time"`
	Acc_x     float64 `json:"acc_x_rms"`
	Acc_y     float64 `json:"acc_y_rms"`
	Acc_z     float64 `json:"acc_z_rms"`
}

// Move struct
type Move struct {
	DeviceId  int `json:"device_id"`
	Timestamp int `json:"time"`
	Move      int `json:"move"`
}

type M map[string]interface{}

// Returns `move` values from DynamoDB
func ScanMoveDB(minTime int64, deviceId int, svc *dynamodb.DynamoDB) ([]int, []M, []M, []M) {
	filt := expression.Name("device_id").Equal(expression.Value(deviceId)).And(expression.Name("time").GreaterThan(expression.Value(minTime)))
	proj := expression.NamesList(
		expression.Name("device_id"),
		expression.Name("time"),
		expression.Name("move"))

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

		timestamps = append(timestamps, move.Timestamp)

		switch move.Move {
		case 1:
			xMovementSlice = append(xMovementSlice, M{"x": move.Timestamp, "y": 1})
		case 2:
			yMovementSlice = append(yMovementSlice, M{"x": move.Timestamp, "y": 2})
		case 3:
			circleMovementSlice = append(circleMovementSlice, M{"x": move.Timestamp, "y": 3})
		}
	}
	return timestamps, xMovementSlice, yMovementSlice, circleMovementSlice
}

// Returns `RMS` values from DynamoDB
func ScanRMSDB(minTime int64, deviceId int, svc *dynamodb.DynamoDB) ([]int64, []float32, []float32, []float32) {
	filt := expression.Name("device_id").Equal(expression.Value(deviceId)).And(expression.Name("time").GreaterThan(expression.Value(minTime)))
	proj := expression.NamesList(
		expression.Name("device_id"),
		expression.Name("time"),
		expression.Name("acc_x_rms"),
		expression.Name("acc_y_rms"),
		expression.Name("acc_z_rms"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("tinyml-rms"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Got error retrieving data: %s", err)
	}

	var xRMSSlice []float32
	var yRMSSlice []float32
	var zRMSSlice []float32
	var timestamps []int64

	for _, i := range result.Items {
		rmsItem := RMS{}
		err := dynamodbattribute.UnmarshalMap(i, &rmsItem)
		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		timestamps = append(timestamps, rmsItem.Timestamp)
		xRMSSlice = append(xRMSSlice, float32(rmsItem.Acc_x))
		yRMSSlice = append(yRMSSlice, float32(rmsItem.Acc_y))
		zRMSSlice = append(zRMSSlice, float32(rmsItem.Acc_z))
	}

	return timestamps, xRMSSlice, yRMSSlice, zRMSSlice

}

// Adds move value to DynamoDB table
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

// Adds RMS values to DynamoDB table
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
