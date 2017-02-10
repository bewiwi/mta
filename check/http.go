package check

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/tcnksm/go-httpstat"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"errors"
)

type HttpCheck struct {
	Host   string `json:"host"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (h *HttpCheck) ValidParam() (error) {
	if h.Host == "" {
		return errors.New("host can't be null")
	}
	return nil
}

func (h *HttpCheck) Run() (*models.CheckResponse, error) {
	var err error
	response := models.NewCheckResponse()

	log.Debug("HTTP: ", h.Host)
	var result httpstat.Result

	req, err := http.NewRequest(h.Method, "https://"+h.Host, nil)
	if err != nil {
		return HandleError(&response, err)
	}
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return HandleError(&response, err)
	}

	// Body in /dev/null
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		return HandleError(&response, err)
	}
	res.Body.Close()
	result.End(time.Now())

	response.Values = map[string]float64{
		"DNSLookup":        result.DNSLookup.Seconds(),
		"TCPConnection":    result.TCPConnection.Seconds(),
		"TLSHandshake":     result.TLSHandshake.Seconds(),
		"ServerProcessing": result.ServerProcessing.Seconds(),
		"NameLookup":       result.NameLookup.Seconds(),
		"Connect":          result.Connect.Seconds(),
		"Pretransfer":      result.Pretransfer.Seconds(),
		"StartTransfer":    result.StartTransfer.Seconds(),
		"Total":            result.Total(time.Now()).Seconds(),
	}

	return &response, err

}
