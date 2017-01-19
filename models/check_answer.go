package models

import (
	"fmt"
	"os"
	"time"
)

type CheckResponse struct {
	CheckMetadata CheckMetadaV1
	Timestamp int64
	Error     string
	Hostname  string
	Values    map[string]float64
}

func (c CheckResponse) Print() {
	fmt.Println("CheckID: ", c.CheckMetadata.Id)
	fmt.Println("Hostname: ", c.Hostname)
	if c.Error != "" {
		fmt.Println("Error: ", c.Error)
	}
	for key, value := range c.Values {
		fmt.Println(key, ": ", value)
	}
}
func NewCheckResponse() CheckResponse {
	response := CheckResponse{}
	response.Timestamp = time.Now().UnixNano()
	hostname, _ := os.Hostname()
	response.Hostname = hostname
	return response
}
