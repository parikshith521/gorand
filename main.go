package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func readCSVFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseCSV(data []byte) (*csv.Reader, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	return reader, nil
}

func processCSV(reader *csv.Reader) {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data: ", err)
			break
		}
		fmt.Println(record)
	}
}

func createCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(f)
	return writer, f, nil
}

func writeCSVRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	if err != nil {
		fmt.Println("Error writing to CSV file: ", err)
	}
}

func printCSVFile(filename string) {
	data, err := readCSVFile("data.csv")
	if err != nil {
		fmt.Println("Error reading CSV file: ", err)
		return
	}
	reader, err := parseCSV(data)
	if err != nil {
		fmt.Println("Error creating CSV reader: ", err)
	}
	processCSV(reader)
}

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
	full_record := append([]string{strconv.Itoa(cnt + 1)}, record...)
	writeCSVRecord(writer, full_record)
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing CSV writer: ", err)
	}
}

func deleteRecordFromCSVFile(filename string, id string) {
	//create a tmp file, read from org file to tmp file, then replace original with new file
	org_file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open original file: ", err)
		return
	}
	defer org_file.Close()

	tmp_file, err := os.CreateTemp("", "data_*.csv")
	if err != nil {
		fmt.Println("Failed to create temporary file: ", err)
	}
	defer tmp_file.Close()
	defer os.Remove(tmp_file.Name())

	org_reader := csv.NewReader(org_file)
	writer := csv.NewWriter(tmp_file)

	for {
		record, err := org_reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Failed to copy to temporary file: ", err)
			break
		} else if record[0] != id {
			if err := writer.Write(record); err != nil {
				fmt.Println("Failed to write to temporary file: ", err)
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error when flushing: ", err)
	}

	org_file.Close()
	tmp_file.Close()
	if err := os.Rename(tmp_file.Name(), filename); err != nil {
		fmt.Println("Failed to rename temporary file: ", err)
	}

}

func main() {
	user_input := os.Args[1:]
	switch user_input[0] {
	case "list":
		printCSVFile("data.csv")
	case "add":
		writeToCSVFile("data.csv", user_input[1:])
	case "delete":
		deleteRecordFromCSVFile("data.csv", user_input[1])
	default:
		fmt.Println("Invalid args")
	}

}
