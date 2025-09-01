package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data: ", err)
			break
		}
		t, err := time.Parse(time.RFC3339, record[len(record)-2])
		if err != nil {
			fmt.Println("Error parsing time from CSV: ", err)
		}
		timeDiff := timediff.TimeDiff(t)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", record[0], record[1], timeDiff, record[3])
	}
	w.Flush()
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all items",
	Run: func(cmd *cobra.Command, args []string) {
		printCSVFile("data.csv")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
