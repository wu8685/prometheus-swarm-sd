package swarm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/prometheus/prometheus/config"
)

type Do func(req *http.Request) (resp *http.Response, err error)

type swarmClient struct {
	client   *http.Client
	masters  []config.URL
	masterMu sync.Mutex

	do Do
}

func newSwarmClient(conf *config.SwarmSDConfig) *swarmClient {
	client := &http.Client{}

	return &swarmClient{
		client:  client,
		masters: conf.Masters,
	}
}

func newSwarmClientForTest(conf *config.SwarmSDConfig, do Do) *swarmClient {
	client := newSwarmClient(conf)
	client.do = do
	return client
}

func (c *swarmClient) getNodeInfo() (*Info, error) {
	c.masterMu.Lock()
	defer c.masterMu.Unlock()

	for _, master := range c.masters {
		urlStr := fmt.Sprintf("%s/info", master.String())
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return nil, err
		}

		var resp *http.Response
		if c.do != nil {
			// test code
			resp, err = c.do(req)
		} else {
			resp, err = c.client.Do(req)
		}
		if err == nil {
			return c.processNodeInfo(resp)
		}

		c.rotateMaster()
	}
	return nil, errors.New("No available master.")
}

func (c *swarmClient) processNodeInfo(resp *http.Response) (*Info, error) {
	if resp == nil {
		return nil, errors.New("Response is nil.")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Fail to access swarm master: %d - %s.", resp.StatusCode, resp.Status)
	}

	var infoEntity InfoEntity
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&infoEntity); err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unable to parse node info response: %s", string(body))
	}

	return infoEntity.toInfo()
}

func (c *swarmClient) rotateMaster() {
	if len(c.masters) > 1 {
		c.masters = append(c.masters[1:], c.masters[0])
	}
}
