package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func main() {
	const jsonStream = `{"sensor_id": "temp-1", "timestamp": 1234567890, "readings": [22.1, 22.3, 22.0], "metadata": {"id": 123}}`
	sensor := SensorParser(strings.NewReader(jsonStream))
	fmt.Println(sensor.sensorID)
	fmt.Println(sensor.value)
}

type SensorData struct {
	sensorID string
	value    float64
}

func SensorParser(r io.Reader) SensorData {
	var sd SensorData
	var count int

	dec := json.NewDecoder(r)

	for {
		t, err := dec.Token()
		if err != nil {
			break
		}

		if str, ok := t.(string); ok {
			switch str {
			case "sensor_id":
				var s string
				if err := dec.Decode(&s); err != nil {
					break
				}
				sd.sensorID = s
				count++
			case "readings":
				var f []float64
				if err := dec.Decode(&f); err != nil {
					break
				}
				sd.value = f[0]
				count++
			}
		}

		if count == 2 {
			break
		}

	}

	return sd
}
