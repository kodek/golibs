package thermostat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

type Client interface {
	GetStatus() (ThermostatStatus, error)
}

type ThermostatStatus struct {
	Temperature    float64
	ThermostatMode ThermostatMode
	FanMode        FanMode
	Override       bool
	Hold           bool
}

// The thermostat mode
type ThermostatMode int

const (
	THERMOSTAT_OFF ThermostatMode = iota
	THERMOSTAT_HEAT
	THERMOSTAT_COOL
	THERMOSTAT_AUTO
)

//go:generate stringer -type ThermostatMode

func (tm *ThermostatMode) UnmarshalJSON(data []byte) error {
	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("Can't parse ThermostatMode: %s", data)
	}

	if v < float64(THERMOSTAT_OFF) || v > float64(THERMOSTAT_AUTO) {
		return fmt.Errorf("Invalid ThermostatMode value: %f", v)
	}

	*tm = ThermostatMode(v)
	return nil
}

// The fan mode
type FanMode int

const (
	FAN_AUTO FanMode = iota
	FAN_AUTO_CIRCULATE
	FAN_ON
)

//go:generate stringer -type FanMode

func (fm *FanMode) UnmarshalJSON(data []byte) error {
	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("Can't parse FanMode: %s", data)
	}

	if v < float64(FAN_AUTO) || v > float64(FAN_ON) {
		return fmt.Errorf("Invalid FanMode value: %f", v)
	}

	*fm = FanMode(v)
	return nil
}

type thermostatClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(host string) (Client, error) {
	// TODO: Validate host. Make sure it is a valid URL.
	return &thermostatClient{
		baseURL:    host,
		httpClient: &http.Client{},
	}, nil
}

func (c *thermostatClient) GetStatus() (status ThermostatStatus, err error) {
	resp, err := c.httpClient.Get(c.baseURL + "/tstat")
	if err != nil {
		return ThermostatStatus{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return status, err
	}

	if resp.StatusCode >= 400 {
		glog.Error("API error: ", resp)
		return status, errors.New(fmt.Sprintf("API status: %i", resp.StatusCode))
	}

	status, err = parseStatus(body)

	return status, err
}

func parseStatus(body []byte) (status ThermostatStatus, err error) {
	var aux struct {
		Temp     float64
		Tmode    ThermostatMode
		FMode    FanMode
		Override int
		Hold     int
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	if err = dec.Decode(&aux); err != nil {
		return status, err
	}

	status.Temperature = aux.Temp
	status.ThermostatMode = aux.Tmode
	status.FanMode = aux.FMode
	status.Override = aux.Override == 1
	status.Hold = aux.Hold == 1
	return status, err
}
