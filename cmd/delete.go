package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func deleteRecordFromCSVFile(filename string, id string, upd bool) {
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
		} else {
			if upd {
				// need to make it true
				if record[0] == id {
					record[len(record)-1] = "true"
				}
				if err := writer.Write(record); err != nil {
					fmt.Println("Failed to write to temporary file: ", err)
				}

			} else {
				if record[0] != id {
					if err := writer.Write(record); err != nil {
						fmt.Println("Failed to write to temporary file: ", err)
					}
				}
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

var deleteCmd = &cobra.Command{
	Use:   "delete [task to delete]",
	Short: "Deletes a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteRecordFromCSVFile("data.csv", args[0], false)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
