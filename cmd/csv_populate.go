package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var length int
var workers int
var counter int64
var file *os.File
var m sync.Mutex

func init() {
	PopulateCmd.PersistentFlags().IntVarP(&workers, "worker", "w", 10, "total worker to be used")
	PopulateCmd.PersistentFlags().IntVarP(&length, "length", "l", 100, "total data to be populated")
}

var PopulateCmd = &cobra.Command{
	Use:   "populate",
	Short: "populate csv file",
	Long:  `populate csv file`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		file, err = os.Create("internal/database/mysql/seeders/file.csv")
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		defer file.Close()

		start := time.Now()

		jobs := make(chan []string)
		wg := new(sync.WaitGroup)

		job := []string{"ID"}
		write(0, job)
		counter++

		go dispatch(workers, jobs, wg)
		sendCsvLineToWorker(jobs, wg)

		wg.Wait()

		duration := time.Since(start)
		log.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
	},
}

func dispatch(totalWorker int, jobs <-chan []string, wg *sync.WaitGroup) {
	for workerIndex := 1; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, jobs <-chan []string, wg *sync.WaitGroup) {
			for job := range jobs {
				write(workerIndex, job)
				wg.Done()
			}
		}(workerIndex, jobs, wg)
	}
}

func sendCsvLineToWorker(jobs chan<- []string, wg *sync.WaitGroup) {
	for {
		if int(counter) > length {
			break
		}

		var row []string
		m.Lock()

		inc := strconv.FormatInt(counter, 10)
		row = []string{inc}
		counter++

		m.Unlock()

		wg.Add(1)
		jobs <- row
	}
	close(jobs)
}

func write(workerIndex int, value []string) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			w := csv.NewWriter(file)
			if err := w.Write(value); err != nil {
				log.Fatalln("error writing record to file", err)
			}
			defer w.Flush()

			log.Println("worker", workerIndex, "write", value)
		}(&outerError)

		if outerError == nil {
			break
		}
	}
}
