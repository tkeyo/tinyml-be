package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	RMSUtil "github.com/tkeyo/tinyml-be/util"

	"github.com/gin-gonic/gin"
)

type dataMove struct {
	Time int    `json:"time"`
	Move string `json:"move"`
}

func saveMoveData(d dataMove) {
	csvFile, err := os.OpenFile("data/dataMove.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	var row []string

	row = append(row, strconv.Itoa(d.Time))
	row = append(row, d.Move)

	writer.Write(row)
	writer.Flush()
}

func getParsedRMSData(records [][]string) ([]int, []float64, []float64, []float64) {
	var time []int
	var acc_x []float64
	var acc_y []float64
	var acc_z []float64

	for _, record := range records {

		timeVal, _ := strconv.Atoi(record[0])
		acc_xVal, _ := strconv.ParseFloat(record[1], 4)
		acc_yVal, _ := strconv.ParseFloat(record[2], 4)
		acc_zVal, _ := strconv.ParseFloat(record[3], 4)

		time = append(time, timeVal)
		acc_x = append(acc_x, acc_xVal)
		acc_y = append(acc_y, acc_yVal)
		acc_z = append(acc_z, acc_zVal)
	}
	return time, acc_x, acc_y, acc_z
}

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

	time, acc_x, acc_y, acc_z := getParsedRMSData(records)

	c.JSON(200, gin.H{
		"time":      time,
		"acc_x_rms": acc_x,
		"acc_y_rms": acc_y,
		"acc_z_rms": acc_z,
	})
}

func getMoveCounts(arr []int) map[int]int {
	dict := make(map[int]int)
	for _, num := range arr {
		dict[num] = dict[num] + 1
	}
	return dict
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

	countsDay1 := getMoveCounts(day1)
	countsDay2 := getMoveCounts(day2)
	countsDay3 := getMoveCounts(day3)

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
	var move dataMove
	c.BindJSON(&move)
	c.JSON(200, gin.H{
		"message": "OK",
	})
	saveMoveData(move)
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
