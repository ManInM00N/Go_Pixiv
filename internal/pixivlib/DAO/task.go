package DAO

import "time"

type TaskInfo struct {
	BeginTime time.Time
	EndTime   time.Time
	Current   uint64
	Total     uint64
	Status    string
	Message   string
	Name      string
	ID        string
}
