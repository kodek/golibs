package thermostat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func testHttpRequest(code int, body string) (*httptest.Server, *thermostatClient) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	client := &thermostatClient{
		baseURL:    server.URL,
		httpClient: httpClient,
	}

	return server, client
}

func TestTemperatureStatus(t *testing.T) {
	var tests = []struct {
		req      string
		expected ThermostatStatus
	}{
		{`{"temp":123.45}`, ThermostatStatus{Temperature: 123.45}},
		{`{"tmode":0}`, ThermostatStatus{ThermostatMode: THERMOSTAT_OFF}},
		{`{"tmode":1}`, ThermostatStatus{ThermostatMode: THERMOSTAT_HEAT}},
		{`{"tmode":2}`, ThermostatStatus{ThermostatMode: THERMOSTAT_COOL}},
		{`{"tmode":3}`, ThermostatStatus{ThermostatMode: THERMOSTAT_AUTO}},
		{`{"fmode":0}`, ThermostatStatus{FanMode: FAN_AUTO}},
		{`{"fmode":1}`, ThermostatStatus{FanMode: FAN_AUTO_CIRCULATE}},
		{`{"fmode":2}`, ThermostatStatus{FanMode: FAN_ON}},
		{`{"hold":0}`, ThermostatStatus{Hold: false}},
		{`{"hold":1}`, ThermostatStatus{Hold: true}},
		{`{"override":0}`, ThermostatStatus{Override: false}},
		{`{"override":1}`, ThermostatStatus{Override: true}},
	}
	for _, test := range tests {
		server, client := testHttpRequest(200, test.req)
		defer server.Close()
		actual, err := client.GetStatus()
		if err != nil {
			t.Errorf("Error GetStatus(%v): %s", test.req, err)
		}
		if actual != test.expected {
			t.Errorf("GetStatus(%v) = %v; want %v", test.req, actual, test.expected)
		}
	}
}

//func ExampleBasic() {
//	c, err := NewClient("http://thermostat.lan")
//	if err != nil {
//		panic(err)
//	}
//
//	ts, err := c.GetStatus()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Thermostat status: %v", ts)
//	// Output:
//	// Thermostat status: {68.5 HEAT}
//}
