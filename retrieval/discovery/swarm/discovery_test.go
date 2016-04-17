package swarm

import (
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
)

func TestSources(t *testing.T) {
	d, err := createDiscovery()
	if err != nil {
		t.Fatalf("Fail to create discovery: %s.", err.Error())
	}

	sources := d.Sources()
	if len(sources) == 0 {
		t.Fatalf("Fail to get sources: sources are empty.")
	}
}

func TestRun(t *testing.T) {
	d, err := createDiscovery()
	if err != nil {
		t.Fatalf("Fail to create discovery: %s", err.Error())
	}
	updateOneNodes(d)

	up := make(chan config.TargetGroup, 10)
	done := make(chan struct{}, 1)

	go d.Run(up, done)

	ticker := time.Tick(3 * time.Second)
	go func() {
		nodeNum := 1
		for {
			select {
			case <-done:
				return
			case <-ticker:
				if nodeNum == 1 {
					updateThreeNodes(d)
					nodeNum = 3
				} else if nodeNum == 3 {
					updateTwoNodes(d)
					nodeNum = 2
				}
			}
		}
	}()

	timeout := time.Tick(10 * time.Second)
	hasScaleDown := false
	hasScaleUp := false
	for {
		isContinue := true
		select {
		case <-timeout:
			done <- struct{}{}
			done <- struct{}{}
			isContinue = false
			break
		case tg := <-up:
			nodeNum := len(tg.Targets)
			if nodeNum == 1 {
				hasScaleUp = false
				hasScaleDown = false
			} else if nodeNum == 2 {
				hasScaleDown = true
			} else if nodeNum == 3 {
				hasScaleUp = true
			}
		}

		if !isContinue {
			break
		}
	}

	if !hasScaleUp {
		t.Fatalf("Has not scaled up.")
	}

	if !hasScaleDown {
		t.Fatalf("Has not scaled down.")
	}
}

func createDiscovery() (*Discovery, error) {
	cfg := &config.DefaultSwarmSDConfig
	cfg.RefreshInterval = model.Duration(1 * time.Second)

	masterUrls := []config.URL{}
	for _, m := range masters {
		url, err := neturl.Parse(m)
		if err != nil {
			return nil, err
		}
		masterUrls = append(masterUrls, config.URL{url})
	}
	cfg.Masters = masterUrls

	sd := &Discovery{
		Conf: cfg,
	}
	err := sd.Initialize()
	if err != nil {
		return nil, err
	}

	sd.client.do = mockOneNodeDo
	return sd, nil
}

func updateOneNodes(d *Discovery) {
	d.client.masterMu.Lock()
	defer d.client.masterMu.Unlock()
	d.client.do = mockOneNodeDo
}

func updateTwoNodes(d *Discovery) {
	d.client.masterMu.Lock()
	defer d.client.masterMu.Unlock()
	d.client.do = mockTwoNodeDo
}

func updateThreeNodes(d *Discovery) {
	d.client.masterMu.Lock()
	defer d.client.masterMu.Unlock()
	d.client.do = mockThreeNodeDo
}

func mockOneNodeDo(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}

	resp.StatusCode = 200
	resp.Body = &MockReadCloser{strings.NewReader(oneNode)}
	return
}

func mockTwoNodeDo(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}

	resp.StatusCode = 200
	resp.Body = &MockReadCloser{strings.NewReader(twoNode)}
	return
}

func mockThreeNodeDo(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{}

	resp.StatusCode = 200
	resp.Body = &MockReadCloser{strings.NewReader(threeNode)}
	return
}

type MockReadCloser struct {
	io.Reader
}

func (r MockReadCloser) Close() error {
	return nil
}
