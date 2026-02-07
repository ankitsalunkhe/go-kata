package main

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fsys := fstest.MapFS{
			"test.conf": &fstest.MapFile{
				Data: []byte(`test=data`),
			},
			"folder/test.conf": &fstest.MapFile{
				Data: []byte(`test=data`),
			},
			"test2.random": &fstest.MapFile{
				Data: []byte(`test=data`),
			},
		}

		resultByte, err := LoadConfigs(fsys, ".")
		require.NoError(t, err)
		require.Equal(t, map[string][]byte{
			"test.conf":        []byte(`test=data`),
			"folder/test.conf": []byte(`test=data`),
		}, resultByte)
	})
}
