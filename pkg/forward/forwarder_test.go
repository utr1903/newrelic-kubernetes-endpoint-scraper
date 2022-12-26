package forward

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

func Test_HttpRequestFailed(t *testing.T) {
	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig("", endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	assert.Panics(t, func() { forwarder.Run() })
}

func Test_NewRelicReturnsNotOkResponse(t *testing.T) {
	newrelicEventServerMock := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
	defer newrelicEventServerMock.Close()

	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig(newrelicEventServerMock.URL, endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	assert.Panics(t, func() { forwarder.Run() })
}

func Test_EventsAreSent(t *testing.T) {
	newrelicEventServerMock := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	defer newrelicEventServerMock.Close()

	endpointInfoMock := createEndpointInfoMock()
	cfg := createConfig(newrelicEventServerMock.URL, endpointInfoMock)
	evs := createEndpointValues(cfg, endpointInfoMock)

	forwarder := NewForwarder(cfg, evs)
	assert.NotPanics(t, func() { forwarder.Run() })
}

func createEndpointValues(
	cfg *config.Config,
	endpointInfoMock map[string](map[string]string),
) *config.EndpointValues {
	evs := config.NewEndpointValues()
	for _, endpoint := range cfg.Endpoints {
		evs.AddEndpointValues(endpoint, endpointInfoMock[endpoint.URL])
	}
	return evs
}

func createEndpointInfoMock() map[string](map[string]string) {
	return map[string](map[string]string){
		"ep1Url": map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
		"ep2Url": map[string]string{
			"k3": "v3",
			"k4": "v4",
		},
	}
}

func createConfig(
	newrelicEventsUrl string,
	endpointInfo map[string](map[string]string),
) *config.Config {
	logLevel := "ERROR"
	eps := []config.Endpoint{}
	for url := range endpointInfo {
		eps = append(eps, config.Endpoint{
			Type: "kvp",
			Name: "my_endpoint_" + url,
			URL:  url,
		})
	}
	return &config.Config{
		Newrelic: &config.NewRelicInput{
			LogLevel:       logLevel,
			EventsEndpoint: newrelicEventsUrl,
			LicenseKey:     "",
		},
		Logger:    config.NewLogger(logLevel),
		Endpoints: eps,
	}
}