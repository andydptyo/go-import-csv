package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func OpenCsvFile(path string) (*csv.Reader, *os.File, error) {
	log.Println("=> open csv file")

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(f)
	return reader, f, nil
}
