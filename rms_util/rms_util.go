package rms_util

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type DataRMS struct {
	Time  int     `json:"time"`
	Acc_x float64 `json:"acc_x"`
	Acc_y float64 `json:"acc_y"`
	Acc_z float64 `json:"acc_z"`
}

func SaveRMSData(d DataRMS) {
	csvFile, err := os.OpenFile("data/dataRMS.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	var row []string

	row = append(row, strconv.Itoa(d.Time))
	row = append(row, strconv.FormatFloat(d.Acc_x, 'f', 4, 64))
	row = append(row, strconv.FormatFloat(d.Acc_y, 'f', 4, 64))
	row = append(row, strconv.FormatFloat(d.Acc_z, 'f', 4, 64))

	writer.Write(row)
	writer.Flush()
}

func GetParsedRMSData(records [][]string) ([]int, []float64, []float64, []float64) {
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
