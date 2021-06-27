package main

import (
	"math"
	"testing"
)

func TestStatus(t *testing.T) {
	tables := []struct {
		response int
		code     int
	}{
		{200, 1},
		{201, 0},
		{503, 0},
		{404, 0},
		{int(math.NaN()), -1},
	}

	for _, table := range tables {
		code, _ := Status(table.response)
		if code != table.code {
			t.Errorf("Up/Down Code of %d was incorrect, got: %d, want: %d", table.response, code, table.code)
		}
	}
}
