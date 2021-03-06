package main

import (
	"flag"
	"log"
	"time"

	"github.com/beloblotskiy/BigQuery-Logs/bqldr"
	"github.com/beloblotskiy/BigQuery-Logs/dmaker"
	"github.com/beloblotskiy/BigQuery-Logs/etlutils"
	"github.com/beloblotskiy/BigQuery-Logs/scanner"
	"github.com/beloblotskiy/BigQuery-Logs/scorer"
)

func main() {
	t0 := time.Now()

	sysPtr := flag.String("sys", "", "System tag, typically server name")
	dirPtr := flag.String("dir", "", "Start dir for log scan")
	flag.Parse()
	if len(*sysPtr) == 0 || len(*dirPtr) == 0 {
		log.Panic("Not all command-line arguments are defined. Please use -h to get more information.")
	} else {
		log.Printf("Started at %v with command-line arguments: sys=%s, dir=%s", t0, *sysPtr, *dirPtr)
	}

	cdcstart := etlutils.CalcCDSStartDate(*sysPtr)
	bqldr.PrepareCDC(*sysPtr, cdcstart)
	// for debugging: etlutils.PrintSR(dmaker.Decide(...))
	<-bqldr.Upload(5, *sysPtr, dmaker.Decide(1, scorer.Score(15, scanner.Scan(cdcstart, *dirPtr))))
	log.Printf("Execution time: %v", time.Since(t0))
}
