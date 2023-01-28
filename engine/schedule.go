package engine

import "mycrawler/collect"

type ScheduleEngine struct {
	requestCh chan *collect.Request
}