package masterclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Bo0mer/os-agent/model"
)

type MasterClient interface {
	Register() error
}

type masterClient struct {
	self       model.Slave
	httpClient *http.Client
	masterURL  string
}

func NewMasterClient(masterUrl string, self model.Slave) MasterClient {
	return &masterClient{
		self:       self,
		httpClient: new(http.Client),
		masterURL:  strings.TrimSuffix(masterUrl, "/"),
	}
}

func (c *masterClient) Register() error {
	// no need to send any data,
	// just say that we're alive
	data, _ := json.Marshal(c.self)
	_, err := c.doPost("/register", data)
	return err
}

func (c *masterClient) doPost(relativeURL string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	url := fmt.Sprintf("%s%s", c.masterURL, relativeURL)

	request, _ := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", "application/json")
	return c.doRequest(request)
}
func (c *masterClient) doRequest(request *http.Request) ([]byte, error) {
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if 200 > response.StatusCode || response.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Response status code: %s", response.StatusCode))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
