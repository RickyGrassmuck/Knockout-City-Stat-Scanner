package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const APIURL string = "https://api.extract-table.com"

type Table [][]string

func scan(filePath string) ([]byte, error) {
	contents, _ := os.ReadFile(filePath)

	imageMimeType := guessImageMimeTypes(bytes.NewReader(contents))
	fmt.Printf("Image File Type: %s\n", imageMimeType)
	response, err := http.Post("https://api.extract-table.com", imageMimeType, bytes.NewReader(contents))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func getFiles() ([]string, error) {
	jpgFiles, err := filepath.Glob("*.j*g")
	if err != nil {
		return nil, err
	}

	pngFiles, err := filepath.Glob("*.png")
	if err != nil {
		return nil, err
	}

	allFiles := append(jpgFiles, pngFiles...)
	return allFiles, nil
}

// Guess image format from gif/jpeg/png/webp
func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// Guess image mime types from gif/jpeg/png/webp
func guessImageMimeTypes(r io.Reader) string {
	format, _ := guessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}

func toCSV(data []byte) (string, error) {
	var t [][]string

	err := json.Unmarshal(data, &t)
	if err != nil {
		return "", err
	}

	s := &bytes.Buffer{}
	writer := csv.NewWriter(s)
	for _, row := range t {
		writer.Write(row)
	}
	writer.Flush()
	return fmt.Sprintf("PLAYERS%s", s.String()), nil

}

func main() {

	files, err := getFiles()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("Found Files: %v\n", files)
	allCSV := []string{}

	for _, f := range files {
		fmt.Printf("Scanning: %s\n", f)

		results, err := scan(f)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		newFile := strings.Replace(f, filepath.Ext(f), ".csv", 1)
		csvResults, err := toCSV(results)
		if err != nil {
			fmt.Printf("Error converting to CSV: %v\n", err)
		}

		fmt.Printf("Creating File %s\n", newFile)
		writeErr := os.WriteFile(newFile, []byte(csvResults), 0644)
		if writeErr != nil {
			fmt.Printf("File Write Error: %v\n", writeErr)
		}
		allCSV = append(allCSV, csvResults)
		time.Sleep(2 * time.Second)
	}
	for i, f := range allCSV {
		fmt.Printf("\n============================\n  Result %d\n============================\n", i+1)
		fmt.Printf("\n%s\n", f)
		fmt.Print("============================\n")
	}
}

// client := &http.Client{
// 	Timeout: time.Second * 10,
// }

// req, err := http.NewRequest(http.MethodPost, APIURL, data)
// if err != nil {
// 	return nil, err
// }
// req.Header.Add("Content-Type", imageMimeType)
// req.Header.Add("Accept", "text/csv")

// fmt.Printf("%v\n", req)
// response, err := client.Do(req)
