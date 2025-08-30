package main

import (
	"fmt"
	"time"
)

type dash0LogsProcessor struct {
	logStats       map[string]uint64
	logIntake      <-chan string
	durationWindow time.Duration
}

func (lp *dash0LogsProcessor) StartLogProcessing() {
	ticker := time.NewTicker(lp.durationWindow).C

	for {
		select {
		case <-ticker:
			fmt.Println("Log stats:")
			for logValue, count := range lp.logStats {
				fmt.Printf("%s - %d\n", logValue, count)
			}
		case logValue := <-lp.logIntake:
			if logValue == "" {
				logValue = "unknown"
			}

			lp.logStats[logValue]++
		}
	}
}
