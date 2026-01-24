package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSensorParser(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		const jsonStream = `{"sensor_id": "temp-1", "timestamp": 1234567890, "readings": [22.1, 22.3, 22.0], "metadata": {"id": 123}}`
		sd := SensorParser(strings.NewReader(jsonStream))
		require.Equal(t, sd.sensorID, "temp-1")
		require.Equal(t, sd.value, 22.1)
	})

	// t.Run("The Corruption Test", func(t *testing.T) {
	// 	const jsonStream = `{"sensor_id": "a"} {"bad json here`
	// 	sensor, reading := SensorParser(strings.NewReader(jsonStream))
	// 	require.Equal(t, sensor, "a")
	// 	require.Equal(t, reading, float64(0))
	// })
}

func BenchmarkSensorParser(b *testing.B) {
	for b.Loop() {
		const jsonStream = `{"sensor_id": "temp-1", "timestamp": 1234567890, "readings": [22.1, 22.3, 22.0], "metadata": {"id": 123}}`
		SensorParser(strings.NewReader(jsonStream))
	}
}
