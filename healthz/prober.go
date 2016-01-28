package healthz

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
)

// Checks the host's healthz page at:
// http://host/healthz
func IsAlive(host string) bool {
	client := http.Client{
		Timeout: time.Duration(1 * time.Second),
	}
	resp, err := client.Get(fmt.Sprintf("http://%s/healthz", host))
	if err != nil {
		glog.Info("healthz connection failed: ", err)
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Info("healthz \"OK\" read fail: ", err)
		return false
	}

	if string(body) != "OK" {
		glog.Info("healthz check not \"OK\": ", string(body))
		return false
	}

	return true
}
