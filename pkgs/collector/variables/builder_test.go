package variables

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	zoneTag := "foo"
	end := time.Date(2022, 11, 22, 13, 30, 0, 0, time.UTC)
	start := end.Add(-time.Minute * 5)
	hostnames := []string{
		"example.com",
		"foo.com",
	}
	v := NewBuilder().
		WithZoneTag(zoneTag).
		WithStart(start).
		WithEnd(end).
		WithHostnames(hostnames).
		Build()

	assert.Equal(t, zoneTag, v["zoneTag"])

	filter := v["filter"].(map[string]interface{})
	andFilter := filter["AND"].([]map[string]interface{})
	assert.Equal(t, start.Format(time.RFC3339), andFilter[0]["datetime_geq"])
	assert.Equal(t, end.Format(time.RFC3339), andFilter[0]["datetime_leq"])

	orFilter := filter["OR"].([]map[string]interface{})
	assert.Len(t, orFilter, 2)
	assert.Equal(t, "example.com", orFilter[0]["clientRequestHTTPHost"])
	assert.Equal(t, "foo.com", orFilter[1]["clientRequestHTTPHost"])

}
