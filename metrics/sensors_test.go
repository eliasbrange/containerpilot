package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func TestSensorParse(t *testing.T) {
	jsonFragment := `{
	"namespace": "namespace_text",
	"subsystem": "subsystem_text",
	"name": "%s",
	"help": "help text",
	"type": "%s",
	"poll": 10,
	"check": ["/bin/sensor.sh"]
}`

	test1Json := []byte(fmt.Sprintf(jsonFragment, "sensor_counter", "counter"))
	collector := parseAndGetCollector(t, test1Json)
	if _, ok := collector.(prometheus.Counter); !ok {
		t.Fatalf("Incorrect collector; expected Counter but got %v", collector)
	}

	test2Json := []byte(fmt.Sprintf(jsonFragment, "sensor_gauge", "gauge"))
	collector = parseAndGetCollector(t, test2Json)
	if _, ok := collector.(prometheus.Gauge); !ok {
		t.Fatalf("Incorrect collector; expected Gauge but got %v", collector)
	}

	test3Json := []byte(fmt.Sprintf(jsonFragment, "sensor_histogram", "histogram"))
	collector = parseAndGetCollector(t, test3Json)
	if _, ok := collector.(prometheus.Histogram); !ok {
		t.Fatalf("Incorrect collector; expected Histogram but got %v", collector)
	}

	test4Json := []byte(fmt.Sprintf(jsonFragment, "sensor_summary", "summary"))
	collector = parseAndGetCollector(t, test4Json)
	if _, ok := collector.(prometheus.Summary); !ok {
		t.Fatalf("Incorrect collector; expected Summary but got %v", collector)
	}

}

func parseAndGetCollector(t *testing.T, testJson []byte) prometheus.Collector {
	sensor := &Sensor{}
	if err := json.Unmarshal(testJson, &sensor); err != nil {
		t.Fatalf("Could not parse sensor JSON: %s", err)
	} else if err := sensor.Parse(); err != nil {
		t.Fatalf("Could not parse sensor check or collector type: %s", err)
	}
	return sensor.collector
}

// invalid collector type
func TestSensorBadType(t *testing.T) {
	jsonFragment := []byte(`{
	"namespace": "namespace_text",
	"subsystem": "subsystem_text",
	"name": "sensor_bad_type",
	"type": "nonsense"}`)

	sensor := &Sensor{}
	if err := json.Unmarshal(jsonFragment, &sensor); err != nil {
		t.Fatalf("Could not parse sensor JSON: %s", err)
	}
	if err := sensor.Parse(); err == nil {
		t.Fatalf("Did not get error from sensor.Parse(): %v", sensor)
	}
}

// invalid metric name
func TestSensorBadName(t *testing.T) {
	jsonFragment := []byte(`{
	"namespace": "namespace_text",
	"subsystem": "subsystem_text",
	"name": "sensor.bad.type",
	"type": "counter"}`)

	sensor := &Sensor{}
	if err := json.Unmarshal(jsonFragment, &sensor); err != nil {
		t.Fatalf("Could not parse sensor JSON: %s", err)
	}
	if err := sensor.Parse(); err == nil {
		t.Fatalf("Did not get error from sensor.Parse(): %v", sensor)
	}
}