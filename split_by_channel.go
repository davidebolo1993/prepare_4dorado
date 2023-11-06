package main

import (
	"context"
	"fmt"
	"sync"
	"strconv"
	"runtime"
	"os"
	"strings"

	"github.com/actforgood/bigcsvreader"
)


func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc / 1024 / 1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc / 1024 / 1024)
	fmt.Printf("\tSys = %v MiB", m.Sys / 1024 / 1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}


const (
	column_filename_fastq = iota
	column_filename_fast5
	column_filename_pod5
	column_parent_read_id
	column_read_id
	column_run_id
	column_channel
	column_mux
	column_minknow_events
	column_start_time
	column_duration
	column_passes_filtering
	column_template_start
	column_num_events_template
	column_template_duration
	column_sequence_length_template
	column_mean_qscore_template
	column_strand_score_template
	column_median_template
	column_mad_template
	column_pore_type
	column_experiment_id
	column_sample_id
	column_end_reason
)

const noOfColumns = 24

type Product struct {

	filename_pod5 string
	read_id string
	channel string
}

func main() {
	// initialize the big csv reader

	summary:=os.Args[1]
	outdir:=os.Args[2]


	bigCSV := bigcsvreader.New()
	bigCSV.SetFilePath(summary)
	bigCSV.FileHasHeader = true
	bigCSV.ColumnsDelimiter = '\t'
	bigCSV.ColumnsCount = noOfColumns
	bigCSV.MaxGoroutinesNo = 1

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	var wg sync.WaitGroup

	// start multi-thread reading
	rowsChans, errsChan := bigCSV.Read(ctx)

	//debug
	//PrintMemUsage()

	info:=make(map[int][]string)
	
	// process rows and errors:

	for i := 0; i < len(rowsChans); i++ {
		wg.Add(1)
		go rowWorker(rowsChans[i], info, &wg)
	}

	wg.Add(1)
	go errWorker(errsChan, &wg)

	wg.Wait()

	//debug
	//PrintMemUsage()

	for key,vals:=range info {

		key_s:=strconv.Itoa(key)

		//fmt.Printf("Storing info for channel %d\n", key)

		f, err := os.Create(outdir+"/"+key_s+".tsv")

		if err != nil {
	
			panic(err)

		}
	
		defer f.Close()

		for _,v:=range vals {

			res:=strings.Split(v,"$")

			if _, err = f.WriteString(res[0] + "\t" + res[1] + "\t" + key_s + "\n"); err != nil {

				panic(err)
			
			}
		
		}

	}


}

func rowWorker(rowsChan bigcsvreader.RowsChan, info map[int][]string, waitGr *sync.WaitGroup) {
	
	i:=0
	
	for row := range rowsChan {

		processRow(row,info)

		//print number of entries processed
		if i%100000 == 0 {

			fmt.Printf("Processed %d entries\n", i)	

		}

		i++

	}
	waitGr.Done()
}

func errWorker(errsChan bigcsvreader.ErrsChan, waitGr *sync.WaitGroup) {

	for err := range errsChan {

		
		handleError(err)
	
	}
	waitGr.Done()
}

// processRow can be used to implement business logic
// like validation / converting to a struct / persisting row into a storage.
func processRow(row []string, info map[int][]string) {

	pod5 := row[column_filename_pod5]
	id := row[column_parent_read_id]
	channel,_:= strconv.Atoi(row[column_channel])

	if row[column_passes_filtering] == "TRUE" {

		info[channel] = append(info[channel], pod5 + "$" + id) //use a single string that will be splitted after

	}

}

// handleError handles the error.
// errors can be fatal like file does not exist, or row related like a given row could not be parsed, etc...
func handleError(err error) {
	fmt.Println(err)
}
