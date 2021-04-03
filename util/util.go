package util_rms

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
