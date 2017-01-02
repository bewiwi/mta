package models

import (
	"fmt"
	"time"
	"os"
)

type CheckAnswer struct {
	RequestID string
	Timestamp int64
	Error     string
	Hostname  string
	Values    map[string]float64
}

func (c CheckAnswer) Print() {
	fmt.Println("RequestID: ", c.RequestID)
	fmt.Println("Hostname: ", c.Hostname)
	if c.Error != ""{
		fmt.Println("Error: ", c.Error)
	}
	for key, value := range c.Values {
		fmt.Println(key, ": ",value)
	}
}
func NewCheckAnswer() CheckAnswer {
	answer := CheckAnswer{}
	answer.Timestamp = time.Now().UnixNano()
	hostname, _ := os.Hostname()
	answer.Hostname = hostname
	return answer
}