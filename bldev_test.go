package bldev

import (
	"testing"
)

func TestBlockDeviceInitialize(t *testing.T) {
	// Test scenarios table
	tests := []struct {
		name     string // Name of the scenario
		initData []byte // Initial data for the block device
	}{
		{
			name:     "Initialize works",
			initData: []byte{'p', 'o', 't', 'a', 't', 'o'},
		},
	}

	// Iterate over the table and run each scenario in isolation
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &blockDevice{}
			bd.Initialize(tt.initData)

			if string(bd.Data) != string(tt.initData) {
				t.Errorf("Expected data: %v, but got: %v", tt.initData, bd.Data)
			}
		})
	}
}

func TestBlockDeviceReadAt(t *testing.T) {
	// Test scenarios table
	tests := []struct {
		name      string // Name of the scenario
		initData  []byte // Initial data for the block device
		offset    int    // Offset to read from
		bufSize   int    // Size of the buffer to read into
		want      []byte // Expected data read
		wantError bool   // Whether we expect an error or not

	}{
		{
			name:      "ReadAt errs when offset is negative",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o'},
			offset:    -1,
			bufSize:   3,
			want:      nil,
			wantError: true,
		},
		{
			name:      "ReadAt errs when offset is out of bounds",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o'},
			offset:    10,
			bufSize:   3,
			want:      nil,
			wantError: true,
		},
		{
			name:      "ReadAt errs when offset + buffer length exceeds data length",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o'},
			offset:    4,
			bufSize:   3,
			want:      nil,
			wantError: true,
		},
		{
			name:      "ReadAt works",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o'},
			offset:    1,
			bufSize:   3,
			want:      []byte{'o', 't', 'a'},
			wantError: false,
		},
	}

	// Iterate over the table and run each scenario in isolation
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &blockDevice{}
			bd.Initialize(tt.initData)

			buffer := make([]byte, tt.bufSize)
			err := bd.ReadAt(tt.offset, buffer, tt.bufSize)

			if tt.wantError && err == nil {
				t.Errorf("Expected error but didn't get one")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
			if !tt.wantError && err == nil {
				if string(buffer) != string(tt.want) {
					t.Errorf("Expected data: %v, but got: %v", tt.want, buffer)
				}
			}
		})
	}
}
