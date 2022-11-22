package cachestats

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

	assert.Len(t, data, 6)

	dynamic := data[0]
	assert.Equal(t, "cacheStatusRequests", *dynamic.MetricName)
	assert.Equal(t, "cacheStatus", *dynamic.Dimensions[0].Name)
	assert.Equal(t, "dynamic", *dynamic.Dimensions[0].Value)
	assert.EqualValues(t, 362, *dynamic.Value)

	expired := data[1]
	assert.Equal(t, "cacheStatusRequests", *expired.MetricName)
	assert.Equal(t, "cacheStatus", *expired.Dimensions[0].Name)
	assert.Equal(t, "expired", *expired.Dimensions[0].Value)
	assert.EqualValues(t, 13, *expired.Value)

	hit := data[2]
	assert.Equal(t, "cacheStatusRequests", *hit.MetricName)
	assert.Equal(t, "cacheStatus", *hit.Dimensions[0].Name)
	assert.Equal(t, "hit", *hit.Dimensions[0].Value)
	assert.EqualValues(t, 2330, *hit.Value)

	miss := data[3]
	assert.Equal(t, "cacheStatusRequests", *miss.MetricName)
	assert.Equal(t, "cacheStatus", *miss.Dimensions[0].Name)
	assert.Equal(t, "miss", *miss.Dimensions[0].Value)
	assert.EqualValues(t, 845, *miss.Value)

}
