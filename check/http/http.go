package http_check

import (
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/bewiwi/mta/check"
	"github.com/tcnksm/go-httpstat"
	"net/http"
	"io"
	"io/ioutil"
)

type HttpCheck struct {
	Host   string `json:"host"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (h *HttpCheck) Run() (*models.CheckResponse, error) {
	var err error
	response := models.NewCheckResponse()

	log.Debug("HTTP: ", h.Host)
	var result httpstat.Result

	req, err := http.NewRequest("GET", "https://"+h.Host, nil)
	if err != nil {
		return check.HandleError(&response, err)
	}
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return check.HandleError(&response, err)
	}

	// Body in /dev/null
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		return check.HandleError(&response, err)
	}
	res.Body.Close()
	result.End(time.Now())

	response.Values = map[string]float64{
		"DNSLookup": result.DNSLookup.Seconds(),
		"TCPConnection": result.TCPConnection.Seconds(),
		"TLSHandshake": result.TLSHandshake.Seconds(),
		"ServerProcessing": result.ServerProcessing.Seconds(),
		"NameLookup": result.NameLookup.Seconds(),
		"Connect": result.Connect.Seconds(),
		"Pretransfer": result.Pretransfer.Seconds(),
		"StartTransfer": result.StartTransfer.Seconds(),
		"Total": result.Total(time.Now()).Seconds(),
	}

	return &response, err

}
