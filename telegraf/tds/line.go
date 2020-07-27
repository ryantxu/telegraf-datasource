package tds

import (
	"encoding/json"
	"math"
	"time"
)

// InfluxLine looks like a line from telegraph json
type InfluxLine struct {
	Name      string                 `json:"name,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Tags      map[string]string      `json:"tags,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}

type serializer struct {
	TimestampUnits time.Duration
}

func NewSerializer(timestampUnits time.Duration) (*serializer, error) {
	s := &serializer{
		TimestampUnits: truncateDuration(timestampUnits),
	}
	return s, nil
}

func (s *serializer) Serialize(metric *InfluxLine) ([]byte, error) {
	m := s.createObject(metric)
	serialized, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	serialized = append(serialized, '\n')

	return serialized, nil
}

func (s *serializer) SerializeBatch(metrics []*InfluxLine) ([]byte, error) {
	objects := make([]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		m := s.createObject(metric)
		objects = append(objects, m)
	}

	obj := map[string]interface{}{
		"metrics": objects,
	}

	serialized, err := json.Marshal(obj)
	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}

func (s *serializer) createObject(metric *InfluxLine) map[string]interface{} {
	m := make(map[string]interface{}, 4)

	tags := make(map[string]string, len(metric.Tags))
	for key, value := range metric.Tags {
		tags[key] = value
	}
	m["tags"] = tags

	fields := make(map[string]interface{}, len(metric.Fields))
	for fieldKey, field := range metric.Fields {
		switch fv := field.(type) {
		case float64:
			// JSON does not support these special values
			if math.IsNaN(fv) || math.IsInf(fv, 0) {
				continue
			}
		}
		fields[fieldKey] = field
	}
	m["fields"] = fields

	m["name"] = metric.Name
	m["timestamp"] = metric.Timestamp.UnixNano() / int64(s.TimestampUnits)
	return m
}

func truncateDuration(units time.Duration) time.Duration {
	// Default precision is 1s
	if units <= 0 {
		return time.Second
	}

	// Search for the power of ten less than the duration
	d := time.Nanosecond
	for {
		if d*10 > units {
			return d
		}
		d = d * 10
	}
}
