package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	MoveUtil "github.com/tkeyo/tinyml-be/move_util"
	RMSUtil "github.com/tkeyo/tinyml-be/rms_util"

	"github.com/gin-gonic/gin"
)

func getRMSData(c *gin.Context) {
	csvFile, err := os.Open("data/dataRMS.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	records := csvData[len(csvData)-10:]

	time, acc_x, acc_y, acc_z := RMSUtil.GetParsedRMSData(records)

	c.JSON(200, gin.H{
		"time":      time,
		"acc_x_rms": acc_x,
		"acc_y_rms": acc_y,
		"acc_z_rms": acc_z,
	})
}

func getMoveData(c *gin.Context) {
	csvFile, err := os.Open("data/dataMove.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	var day1 []int
	var day2 []int
	var day3 []int

	for _, row := range csvData {
		day := row[0]
		move, _ := strconv.Atoi(row[1])

		if day == "300" {
			day1 = append(day1, move)
		} else if day == "400" {
			day2 = append(day2, move)
		} else if day == "500" {
			day3 = append(day3, move)
		}
	}

	countsDay1 := MoveUtil.GetMoveCounts(day1)
	countsDay2 := MoveUtil.GetMoveCounts(day2)
	countsDay3 := MoveUtil.GetMoveCounts(day3)

	c.JSON(200, gin.H{
		"day_1": countsDay1,
		"day_2": countsDay2,
		"day_3": countsDay3,
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server ON",
	})
}

func endpointRMS(c *gin.Context) {
	var rms RMSUtil.DataRMS
	c.BindJSON(&rms)
	c.JSON(200, gin.H{
		"message": "OK",
	})
	RMSUtil.SaveRMSData(rms)
}

func endpointMove(c *gin.Context) {
	var move MoveUtil.DataMove
	c.BindJSON(&move)
	c.JSON(200, gin.H{
		"message": "OK",
	})
	MoveUtil.SaveMoveData(move)
}

func main() {
	fmt.Println("Server is running....")

	r := gin.Default()
	r.GET("/api/health", healthCheck)
	r.GET("/api/get-rms", getRMSData)
	r.GET("/api/get-move", getMoveData)
	r.POST("/api/rms", endpointRMS)
	r.POST("/api/move", endpointMove)
	r.Run()

}
