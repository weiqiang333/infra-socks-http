package model

import (
	"log"
)

// access log: $data $statusCode $receiver $status $grade $alertname $[alertSummary|msg]
func init(){
	log.SetFlags(log.Ldate|log.Lmicroseconds|log.LUTC)
}