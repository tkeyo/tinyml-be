package move_util

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type DataMove struct {
	Time int    `json:"time"`
	Move string `json:"move"`
}

func SaveMoveData(d DataMove) {
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
