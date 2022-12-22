package scraper

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/newrelic-kubernetes-endpoint-scraper/pkg/config"
)

// Object which is responsible for scraping
type EndpointScraper struct {
	config *config.Config
	client *http.Client
	evs    *config.EndpointValues
}

var readResponseBody = func(
	body io.ReadCloser,
) (
	[]byte,
	error,
) {
	return ioutil.ReadAll(body)
}

// Creates new scraper for endpoints
func NewScraper(
	cfg *config.Config,
) *EndpointScraper {

	// Create HTTP client
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	evs := config.NewEndpointValues()

	cfg.Logger.Log(logrus.DebugLevel, "Config file is succesfully created.")

	return &EndpointScraper{
		config: cfg,
		client: &client,
		evs:    evs,
	}
}

// Scrape endpoints
func (s *EndpointScraper) Run() {

	// Loop & scrape all endpoints
	s.config.Logger.Log(logrus.DebugLevel, "Looping over the endpoints to scrape...")
	for _, endpoint := range s.config.Endpoints {

		s.config.Logger.LogWithFields(logrus.DebugLevel, "Scraping endpoint...",
			map[string]string{
				"endpointType": endpoint.Type,
				"endpointName": endpoint.Name,
				"endpointUrl":  endpoint.URL,
			})

		// Create HTTP request
		req, err := http.NewRequest(http.MethodGet, endpoint.URL, nil)
		if err != nil {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, "HTTP request could not be created.",
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
					"error":        err.Error(),
				})
			return
		}

		// Perform HTTP request
		res, err := s.client.Do(req)
		if err != nil {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, "HTTP request could not be created.",
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
					"error":        err.Error(),
				})
			return
		}
		defer res.Body.Close()

		// Check if call was successful
		if res.StatusCode != http.StatusOK {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, "HTTP request has returned not OK status.",
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
				})
			return
		}

		// Extract response body
		body, err := readResponseBody(res.Body)
		if err != nil {
			s.config.Logger.LogWithFields(logrus.ErrorLevel, "Response body could not be parsed.",
				map[string]string{
					"endpointType": endpoint.Type,
					"endpointName": endpoint.Name,
					"endpointUrl":  endpoint.URL,
					"error":        err.Error(),
				})
			return
		}

		// Parse response body
		switch endpoint.Type {
		case "kvp":
			s.parse(&KvpParser{}, endpoint, body)
		}
	}
}

func (s *EndpointScraper) parse(
	p Parser,
	endpoint config.Endpoint,
	data []byte,
) {
	s.evs.AddEndpointValues(endpoint, p.Run(data))

	s.config.Logger.LogWithFields(logrus.DebugLevel, "Endpoint values are parsed.",
		map[string]string{
			"endpointType": endpoint.Type,
			"endpointName": endpoint.Name,
			"endpointUrl":  endpoint.URL,
		})
}

func (s *EndpointScraper) GetEndpointValues() *config.EndpointValues {
	return s.evs
}
