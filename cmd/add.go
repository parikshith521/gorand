package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func countCSVFileRecords(filename string) int {
	data, err := readCSVFile(filename)
	if err != nil {
		fmt.Println("Error reading CSV file: ", err)
		return -1
	}
	reader, err := parseCSV(data)
	if err != nil {
		fmt.Println("Error creating CSV reader: ", err)
	}
	cnt := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data: ", err)
			return -1
		}
		cnt++
	}
	return cnt
}

func writeCSVRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	if err != nil {
		fmt.Println("Error writing to CSV file: ", err)
	}
}

func writeToCSVFile(filename string, record []string) {
	// writer, file, err := createCSVWriter(filename)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error creating CSV File writer: ", err)
		return
	}
	writer := csv.NewWriter(file)
	defer file.Close()
	cnt := countCSVFileRecords(filename)
	fmt.Println("COUNT: ", cnt)
	if cnt == -1 {
		fmt.Println("Error seeking CSV end")
		return
	}
	//record should be combined into a complete string...
	full_record := append([]string{strconv.Itoa(cnt + 1)}, strings.Join(record, " "), time.Now().Format(time.RFC3339), "false")
	writeCSVRecord(writer, full_record)
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing CSV writer: ", err)
	}
}

var addCmd = &cobra.Command{
	Use:   "add [task to add]",
	Short: "Adds a new task",
	Run: func(cmd *cobra.Command, args []string) {
		writeToCSVFile("data.csv", args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
