package variables

import "time"

// Builder is a variables builder.
type Builder struct {
	zoneTag   string
	start     time.Time
	end       time.Time
	hostnames []string
}

// NewBuilder creates a new builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// WithZoneTag adds a zone tag.
func (b *Builder) WithZoneTag(zoneTag string) *Builder {
	b.zoneTag = zoneTag
	return b
}

// WithStart adds a start.
func (b *Builder) WithStart(start time.Time) *Builder {
	b.start = start
	return b
}

// WithEnd adds an end.
func (b *Builder) WithEnd(end time.Time) *Builder {
	b.end = end
	return b
}

// WithHostnames adds hostnames.
func (b *Builder) WithHostnames(hostnames []string) *Builder {
	b.hostnames = hostnames
	return b
}

// Build builds the variables struct.
func (b *Builder) Build() map[string]interface{} {
	variables := map[string]interface{}{
		"zoneTag": b.zoneTag,
		"filter": map[string]interface{}{
			"AND": []map[string]interface{}{
				{
					"datetime_geq": b.start.Format(time.RFC3339),
					"datetime_leq": b.end.Format(time.RFC3339),
				},
				{
					"requestSource": "eyeball",
				},
			},
		},
	}
	if len(b.hostnames) > 0 {
		var orFilter []map[string]interface{}
		for _, hostname := range b.hostnames {
			orFilter = append(orFilter, map[string]interface{}{
				"clientRequestHTTPHost": hostname,
			})
		}
		variables["filter"].(map[string]interface{})["OR"] = orFilter
	}
	return variables
}
