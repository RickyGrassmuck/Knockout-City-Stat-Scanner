package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const APIURL string = "https://api.extract-table.com"

func ScanImage(data []byte) ([][]string, error) {
	imageMimeType := guessImageMimeTypes(bytes.NewReader(data))
	response, err := http.Post("https://api.extract-table.com", imageMimeType, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var table [][]string

	err = json.Unmarshal(responseData, &table)
	if err != nil {
		return nil, err
	}

	normalizedData := [][]string{}
	normalizedData = append(normalizedData, table[0])
	for _, row := range table[1:6] {
		normalizedData = append(normalizedData, NormalizeStatsOutput(row))
	}

	return normalizedData, nil

}

func toCSV(table [][]string) (string, error) {
	s := &bytes.Buffer{}
	writer := csv.NewWriter(s)
	for _, row := range table {
		writer.Write(row)
	}
	writer.Flush()
	return fmt.Sprintf("PLAYERS%s", s.String()), nil
}

// This function exists to fix common issues seen when extracting data from the images
// Ex. The scanner frequently detects zeros as the letter 'o', the function replaces those
//     occurrences in columns where we would never see letters.
func NormalizeStatsOutput(rowData []string) []string {
	normalizedRow := []string{}
	normalizedRow = append(normalizedRow, rowData[0])

	for _, val := range rowData[1:] {
		normalizedRow = append(normalizedRow, strings.Replace(val, "o", "0", -1))
	}

	return normalizedRow
}
