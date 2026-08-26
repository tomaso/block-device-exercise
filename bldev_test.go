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
			initData: []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
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
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    -BLOCK_SIZE,
			bufSize:   BLOCK_SIZE,
			want:      nil,
			wantError: true,
		},
		{
			name:      "ReadAt errs when offset is out of bounds",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    BLOCK_SIZE * 4,
			bufSize:   BLOCK_SIZE,
			want:      nil,
			wantError: true,
		},
		{
			name:      "ReadAt errs when offset + buffer length exceeds data length",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    2 * BLOCK_SIZE,
			bufSize:   3 * BLOCK_SIZE,
			want:      nil,
			wantError: true,
		},
		{
			name:      "ReadAt works",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    BLOCK_SIZE,
			bufSize:   BLOCK_SIZE,
			want:      []byte{'t', 'o', ' ', 's'},
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

func TestBlockDeviceWriteAt(t *testing.T) {
	// Test scenarios table
	tests := []struct {
		name      string // Name of the scenario
		initData  []byte // Initial data for the block device
		offset    int    // Offset to write to
		buffer    []byte // Data to be written
		afterData []byte // Data expected after write
		wantError bool   // Whether we expect an error or not

	}{
		{
			name:      "WriteAt errs when offset is negative",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    -BLOCK_SIZE,
			buffer:    []byte{'a', 's', 'd', 'f'},
			afterData: nil,
			wantError: true,
		},
		{
			name:      "WriteAt errs when offset is out of bounds",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    4 * BLOCK_SIZE,
			buffer:    []byte{'a', 's', 'd', 'f'},
			afterData: nil,
			wantError: true,
		},
		{
			name:      "WriteAt errs when buffer does not fit in the backend data",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    2 * BLOCK_SIZE,
			buffer:    []byte{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k'},
			afterData: nil,
			wantError: true,
		},
		{
			name:      "WriteAt writes buffer to backend data",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offset:    BLOCK_SIZE,
			buffer:    []byte{'a', 's', 'd', 'f'},
			afterData: []byte{'p', 'o', 't', 'a', 'a', 's', 'd', 'f', 'a', 'l', 'a', 'd'},
			wantError: false,
		},
	}

	// Iterate over the table and run each scenario in isolation
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &blockDevice{}
			bd.Initialize(tt.initData)

			err := bd.WriteAt(tt.offset, tt.buffer)

			if tt.wantError && err == nil {
				t.Errorf("Expected error but didn't get one")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
			if !tt.wantError && err == nil {
				dumpedData := bd.Dump()
				if string(dumpedData) != string(tt.afterData) {
					t.Errorf("Expected new data to be: %v, but got: %v", string(tt.afterData), string(dumpedData))
				}
			}
		})
	}
}

func TestBlockDeviceWriteAt_MultipleWrites(t *testing.T) {
	// Test scenarios table
	tests := []struct {
		name      string   // Name of the scenario
		initData  []byte   // Initial data for the block device
		offsets   []int    // Offsets to write to
		buffers   [][]byte // Data to be written
		afterData []byte   // Data expected after write

	}{
		{
			name:      "WriteAt writes 2 buffers to the same offset",
			initData:  []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
			offsets:   []int{BLOCK_SIZE, BLOCK_SIZE},
			buffers:   [][]byte{{'a', 's', 'd', 'f'}, {'g', 'h', 'j', 'k'}},
			afterData: []byte{'p', 'o', 't', 'a', 'g', 'h', 'j', 'k', 'a', 'l', 'a', 'd'},
		},
	}

	// Iterate over the table and run each scenario in isolation
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd := &blockDevice{}
			bd.Initialize(tt.initData)

			for i, offset := range tt.offsets {
				err := bd.WriteAt(offset, tt.buffers[i])
				if err != nil {
					t.Errorf("Error occurred while writing at offset %d: %v", offset, err)
				}
			}
		})
	}
}

// func TestBlockDeviceSerialize(t *testing.T) {
// 	// Test scenarios table
// 	tests := []struct {
// 		name     string // Name of the scenario
// 		initData []byte // Initial data for the block device
// 		expOut   []blockDeviceChange
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "Serialize returns expected output",
// 			initData: nil,
// 			expOut:   []blockDeviceChange{},
// 			wantErr:  false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			bd := &blockDevice{}
// 			bd.Initialize(tt.initData)

// 			changes, err := bd.Serialize()
// 			if (err == nil) && !tt.wantErr {
// 				for i, _ := range changes {
// 					if changes[i].Offset != tt.expOut[i].Offset || changes[i].Data == tt.expOut[i].Data {
// 						t.Errorf(
// 							"Expected change %d: (offset: %d, bytes: %v), but we got: (offset: %d, bytes: %v)",
// 							i, tt.expOut[i].Offset, tt.expOut[i].Data, changes[i].Offset, changes[i].Data)
// 					}
// 				}
// 			} else {
// 				t.Errorf("Error occured: %v", err)
// 			}
// 		})
// 	}
// }

// func TestBlockDeviceDeserialize(t *testing.T) {
// 	tests := []struct {
// 		name              string
// 		initSerializeData []byte
// 		bdcSerialize      []blockDeviceChange
// 		expData           []byte // Expected data in the newly constructed blockDevice
// 	}{
// 		{
// 			name:              "Deserialize creates expected block device",
// 			initSerializeData: []byte{'p', 'o', 't', 'a', 't', 'o', ' ', 's', 'a', 'l', 'a', 'd'},
// 			bdcSerialize: []blockDeviceChange{
// 				{Offset: 0, Data: [4]byte{'a', 'b', 'c', 'd'}},
// 				{Offset: 4, Data: [4]byte{'e', 'f', 'g', 'h'}},
// 			},
// 			expData: []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'a', 'l', 'a', 'd'},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			bd, err := Deserialize(tt.bdcSerialize, tt.initSerializeData)
// 			if err != nil {
// 				t.Errorf("Deserialization failed: %v", err)
// 			}
// 			dumpedBytes := bd.Dump()
// 			if !bytes.Equal(tt.expData, dumpedBytes) {
// 				t.Errorf(
// 					"Expected data from deserialize: %v, actual dump of the blockDevice: %v",
// 					string(tt.expData), string(dumpedBytes),
// 				)
// 			}
// 		})
// 	}
// }
