package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/andydptyo/go-import-csv/internal/config"
	internalCsv "github.com/andydptyo/go-import-csv/internal/csv"
	"github.com/andydptyo/go-import-csv/internal/database/mysql"
	"github.com/spf13/cobra"
)

var totalWorker = 100
var csvFile string
var dataHeaders = make([]string, 0)
var databaseStore *mysql.Mysql

func init() {
	ImportCmd.PersistentFlags().StringVarP(&csvFile, "file", "f", "internal/database/mysql/seeders/file.csv", "csv file to be imported")
}

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "import csv file",
	Long:  `import csv file`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		c, err := config.FromFile(Cfg)
		if err != nil {
			log.Fatalf("error creating config %v", err)
		}

		if c.Database != nil {
			db, err := mysql.New(c.Database)
			if err != nil {
				log.Fatalf("error connecting to database %v", err)
			}

			databaseStore = db
		}

		start := time.Now()

		csvReader, csvFile, err := internalCsv.OpenCsvFile(csvFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer csvFile.Close()

		jobs := make(chan []interface{})
		wg := new(sync.WaitGroup)

		go dispatchWorkers(totalWorker, databaseStore, jobs, wg)
		readCsvFilePerLineThenSendToWorker(csvReader, jobs, wg)

		wg.Wait()

		duration := time.Since(start)
		log.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
	},
}

func dispatchWorkers(totalWorker int, db *mysql.Mysql, jobs <-chan []interface{}, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, db *mysql.Mysql, jobs <-chan []interface{}, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				doTheJob(workerIndex, counter, db, job)
				wg.Done()
				counter++
			}
		}(workerIndex, db, jobs, wg)
	}
}

func readCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- []interface{}, wg *sync.WaitGroup) {
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		if len(dataHeaders) == 0 {
			dataHeaders = row
			continue
		}

		rowOrdered := make([]interface{}, 0)
		for _, each := range row {
			rowOrdered = append(rowOrdered, each)
		}

		wg.Add(1)
		jobs <- rowOrdered
	}
	close(jobs)
}

func doTheJob(workerIndex int, counter int, db *mysql.Mysql, values []interface{}) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			conn, err := db.DB.Conn(context.Background())
			if err != nil {
				log.Fatal(err.Error())
			}

			query := fmt.Sprintf("INSERT INTO csv (%s) VALUES (%s)",
				strings.Join(dataHeaders, ","),
				strings.Join(generateQuestionsMark(len(dataHeaders)), ","),
			)

			_, err = conn.ExecContext(context.Background(), query, values...)
			if err != nil {
				log.Fatal(err.Error())
			}

			err = conn.Close()
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Println("=> worker", workerIndex, "insert", values)
		}(&outerError)

		if outerError == nil {
			break
		}
	}
}

func generateQuestionsMark(n int) []string {
	s := make([]string, 0)
	for i := 0; i < n; i++ {
		s = append(s, "?")
	}
	return s
}
