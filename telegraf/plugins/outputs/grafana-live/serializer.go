package grafanalive

import (
	"encoding/json"
	"math"

	"github.com/influxdata/telegraf"
)

// slightly different output format than the standard JSON
type serializer struct {
	TimestampUnits int64
}

func (s *serializer) Serialize(metric telegraf.Metric) ([]byte, error) {
	m := s.createObject(metric)
	serialized, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	serialized = append(serialized, '\n')

	return serialized, nil
}

func (s *serializer) SerializeBatch(metrics []telegraf.Metric) ([]byte, error) {
	objects := make([]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		m := s.createObject(metric)
		objects = append(objects, m)
	}

	obj := map[string]interface{}{
		"measures": objects,
	}

	serialized, err := json.Marshal(obj)
	if err != nil {
		return []byte{}, err
	}
	return serialized, nil
}

func (s *serializer) createObject(metric telegraf.Metric) map[string]interface{} {
	m := make(map[string]interface{}, 4)

	m["name"] = metric.Name()
	m["timestamp"] = metric.Time().UnixNano() / s.TimestampUnits

	tagCount := len(metric.TagList())
	if tagCount > 0 {
		labels := make(map[string]string, tagCount)
		for _, tag := range metric.TagList() {
			labels[tag.Key] = tag.Value
		}
		m["labels"] = labels
	}

	values := make(map[string]interface{}, len(metric.FieldList()))
	for _, field := range metric.FieldList() {
		switch fv := field.Value.(type) {
		case float64:
			// JSON does not support these special values
			if math.IsNaN(fv) || math.IsInf(fv, 0) {
				continue
			}
		}
		values[field.Key] = field.Value
	}
	m["values"] = values
	return m
}
