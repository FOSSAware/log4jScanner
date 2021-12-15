package cmd

import (
	"encoding/csv"
	"os"
	"strings"
	"sync"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

var csvPath string
var csvMU sync.Mutex

func checkCsvPath() {
	if csvPath == "" {
		csvPath = "log4jScanner-results.csv"
	}

	if !strings.HasSuffix(strings.ToLower(csvPath), ".csv") {
		pterm.Warning.Println("csv-output path is not a CSV file. Output will be saved to running folder as log4jScanner-results.csv")
		csvPath = "log4jScanner-results.csv"
	}
}

func initCSV() {
	checkCsvPath()

	// Set headers
	csvRecords := [][]string{
		{"type", "ip"},
	}

	// create a CSV file and write headers
	f, err := os.Create(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range csvRecords {
		if err = w.Write(record); err != nil {
			log.Fatal(err)
		}
	}
}

func readCsv() (csvRecords [][]string, err error) {
	// open and read csv
	f, err := os.Open(csvPath)
	if err == os.ErrNotExist {
		f, _ = os.Create(csvPath)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	csvRecords, err = r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func updateCsvRecords(resultMsg string) {
	csvMU.Lock()
	defer csvMU.Unlock()

	// open and read csv
	csvRecords, err := readCsv()
	if err != nil {
		log.Fatal(err)
	}

	// append new result to existing CSV content
	fullRes := strings.Split(resultMsg, ",")
	csvRes := []string{fullRes[0], fullRes[1]} // only write to CSV 'request' and the IP address
	csvRecords = append(csvRecords, csvRes)

	// write current and new content to CSV
	writeCSV(csvRecords)
}

func writeCSV(csvRecords [][]string) {
	// load CSV file and write
	f, err := os.Create(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range csvRecords {
		if err = w.Write(record); err != nil {
			log.Fatal(err)
		}
	}
}
