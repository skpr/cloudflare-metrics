package totals

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

	assert.Len(t, data, 2)

	d1 := data[0]
	assert.Equal(t, "totalRequests", *d1.MetricName)
	assert.EqualValues(t, float64(4232), *d1.Value)

	d2 := data[1]
	assert.Equal(t, "totalResponseBytes", *d2.MetricName)
	assert.EqualValues(t, int64(91319486), *d2.Value)
}
