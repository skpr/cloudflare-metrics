package statuscodes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/cloudflare-metrics/pkgs/util"
)

func TestCollect(t *testing.T) {
	// t.Skip("Do not test using real config")
	responseData, err := os.ReadFile("testdata/response.json")
	assert.NoError(t, err)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, err = res.Write(responseData)
		assert.NoError(t, err)
	}))
	defer func() { testServer.Close() }()

	config, err := util.LoadConfig("testdata")
	assert.NoError(t, err)
	httpClient := http.DefaultClient

	client := graphql.NewClient(testServer.URL, httpClient)

	fetcher := NewCollector(config, client)
	end := time.Date(2022, 11, 22, 13, 30, 0, 0, time.UTC)
	start := end.Add(-time.Minute * 5)

	data, err := fetcher.CollectMetrics(context.Background(), start, end)
	assert.NoError(t, err)

	assert.Len(t, data, 8)

	status200 := data[0]
	assert.Equal(t, "statusCodeRequests", *status200.MetricName)
	assert.Equal(t, "statusCode", *status200.Dimensions[0].Name)
	assert.Equal(t, "200", *status200.Dimensions[0].Value)
	assert.EqualValues(t, 3182, *status200.Value)

	status202 := data[1]
	assert.Equal(t, "statusCodeRequests", *status202.MetricName)
	assert.Equal(t, "statusCode", *status202.Dimensions[0].Name)
	assert.Equal(t, "202", *status202.Dimensions[0].Value)
	assert.EqualValues(t, 235, *status202.Value)

	status301 := data[2]
	assert.Equal(t, "statusCodeRequests", *status301.MetricName)
	assert.Equal(t, "statusCode", *status301.Dimensions[0].Name)
	assert.Equal(t, "301", *status301.Dimensions[0].Value)
	assert.EqualValues(t, 146, *status301.Value)

}
