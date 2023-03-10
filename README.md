# New Relic Kubernetes Endpoint Scraper

This repository is meant to scrape the values which are exposed
by the custom endpoints of your applications running in any
pod.

The scraper is set up to run as a cron job every minute.
It will be triggered automatically by Kubernetes and fetch
all the endpoints that you have defined in the configuration.

## Configuration

In order to let the scraper know which endpoints to check, the
following configuration options are provided within the
[values.yaml](/charts/scraper/values.yaml) file.

**Example:**

```yaml
# Configuration for the scraper (main application)
scraper:
  image:
    # For AMD -> newrelic-kubernetes-endpoint-scraper-amd
    # For ARM -> newrelic-kubernetes-endpoint-scraper-arm
    repository: uturkarslan/newrelic-kubernetes-endpoint-scraper-amd
    pullPolicy: Always
    # Overrides the image tag whose default is the chart appVersion.
    tag: "1.1.0"
  # Override name of the application
  nameOverride: ""
  # Override fullname of the application
  fullnameOverride: ""
  # Mount path for the container
  mountPathConfig: /etc/config
  # Configuration data itself
  config:
    newrelic:
      # New Relic account ID
      accountId: ""
      # New Relic license key
      licenseKey: ""
      # Log level can be: DEBUG, ERROR
      logLevel: ERROR
      # Flag to enable log forwarding to New Relic
      logForwarding: true
    # Endpoints which are to be scraped
    # - type
    #   - kvp: key value pair
    #   - json: json
    endpoints: []
      # - type: "kvp"
      #   name: "MyEndpoint1"
      #   url: "http://<SERVICE>.<NAMESPACE>.svc.cluster.local:<PORT>/<ENDPOINT>"
      # - type: "kvp"
      #   name: "MyEndpoint2"
      #   url: "http://<IP_ADDRESS_OF_POD>:<PORT>/<ENDPOINT>"
```

## Scraping

Currently only endpoints which are exposing key-value pairs (`kvp`)
can be scraped and formatted. The key and the value should be separated
by a semicolon (`:`).

## Building your Docker image

If you would like to make your changes to the code and create your
own image, refer to the [`build.sh`](/scripts/build.sh). You can
build to `amd` or `arm` processors.

## Deploying the Helm chart

In order to deploy the solution, refer to the [`deploy.sh`](/scripts/deploy.sh).
You can build to `amd` or `arm` processors.

## Querying your data in New Relic

The scraped endpoints will be forwarded to New Relic as custom
events which will have the table name that you have defined within
the corresponding endpoint block in the config file.

You can query the individual endpoint data as follows:

```
FROM MyEndpoint1 SELECT *
FROM MyEndpoint2 SELECT *
```
