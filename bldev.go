package bldev

import (
	"fmt"
)

type blockDevice struct {
	// Data of the block device
	Initialized bool
	Data        []byte
	Changes     map[int]blockDeviceChange // block (offset % BLOCK_SIZE) -> block data
}

type blockDeviceChange struct {
	Offset int
	Data   [BLOCK_SIZE]byte
}

/*
* initial data

** This is the “base” representing the shared filesystem
** It should not be modified
** Assume the initial data is divisible by the “block size”
** Assume that we won’t be reading/writing outside of the length of the data
 */
/* We assume here that the data will be owned by the block device and that it will not be modified outside of the block device */
func (bd *blockDevice) Initialize(initData []byte) error {
	if bd.Initialized {
		return fmt.Errorf("Block device already initialized")
	}
	bd.Data = initData
	bd.Changes = make(map[int]blockDeviceChange)
	bd.Initialized = true
	return nil
}

/*
* ReadAt

** Takes a byte offset and a data byte slice with a length corresponding to the number of bytes we want to read. The data from the device backend will be read into the data
** Assume that the offset and length to read will always be a multiple of the “block size”
 */
func (bd *blockDevice) ReadAt(offset int, buffer []byte, length int) error {
	if !bd.Initialized {
		return fmt.Errorf("Block device needs to be initialized first")
	}
	if offset < 0 || offset >= len(bd.Data) {
		return fmt.Errorf("invalid offset")
	}
	if length < 0 {
		return fmt.Errorf("invalid length")
	}
	if offset+length > len(bd.Data) {
		return fmt.Errorf("read beyond device bounds")
	}
	if len(buffer)%BLOCK_SIZE != 0 {
		return fmt.Errorf("The buffer length %d is not multiple of BLOCK_SIZE: %d", len(buffer), BLOCK_SIZE)
	}
	if offset%BLOCK_SIZE != 0 {
		return fmt.Errorf("The offset %d is not multiple of BLOCK_SIZE: %d", offset, BLOCK_SIZE)
	}
	if length%BLOCK_SIZE != 0 {
		return fmt.Errorf("The length %d is not multiple of BLOCK_SIZE: %d", length, BLOCK_SIZE)
	}

	block_index_offset := offset / BLOCK_SIZE
	for i := 0; i < length/BLOCK_SIZE; i++ {
		change, exists := bd.Changes[i+block_index_offset]
		if exists {
			copy(buffer[i*BLOCK_SIZE:(i+1)*BLOCK_SIZE], change.Data[:])
		} else {
			copy(buffer[i*BLOCK_SIZE:(i+1)*BLOCK_SIZE], bd.Data[offset:offset+length])
		}
	}

	return nil
}

/*
* WriteAt

** It takes a byte offset and a data byte slice containing the data we want to write. The content of the data will be written into the device backend
** Assume that the offset and length to write will always be a multiple of the “block size”
 */
func (bd *blockDevice) WriteAt(offset int, buffer []byte) error {
	if !bd.Initialized {
		return fmt.Errorf("Block device needs to be initialized first")
	}
	if len(buffer)%BLOCK_SIZE != 0 {
		return fmt.Errorf("The buffer length %d is not multiple of BLOCK_SIZE: %d", len(buffer), BLOCK_SIZE)
	}
	if offset%BLOCK_SIZE != 0 {
		return fmt.Errorf("The offset %d is not multiple of BLOCK_SIZE: %d", offset, BLOCK_SIZE)
	}

	length := len(buffer)

	if offset < 0 || offset >= len(bd.Data) {
		return fmt.Errorf("invalid offset")
	}
	if length < 0 {
		return fmt.Errorf("invalid length")
	}
	if offset+length > len(bd.Data) {
		return fmt.Errorf("write beyond device bounds")
	}

	block_index_offset := offset / BLOCK_SIZE
	for i := 0; i < length/BLOCK_SIZE; i++ {
		// fmt.Printf("i: %d\n", i)
		change, exists := bd.Changes[i+block_index_offset]
		if exists {
			copy(change.Data[:], buffer[i*BLOCK_SIZE:(i+1)*BLOCK_SIZE])
		} else {
			change.Offset = (i + block_index_offset) * BLOCK_SIZE
			copy(change.Data[:], buffer[i*BLOCK_SIZE:(i+1)*BLOCK_SIZE])
			bd.Changes[i+block_index_offset] = change
		}
	}
	// fmt.Printf("bd.Changes: %v\n", bd.Changes)

	return nil
}

/*
* Serialize

** Return a byte slice with serialized changes made to the block device backend
** The format should have the smallest possible overhead
 */
func (bd *blockDevice) Serialize() ([]blockDeviceChange, error) {
	if !bd.Initialized {
		return nil, fmt.Errorf("Block device needs to be initialized first")
	}

	// return []blockDeviceChange{blockDeviceChange{Offset: 0, Data: bd.Data}}, nil
	return nil, nil
}

/*
* Deserialize
(this is not a method of blockDevice, it's a global function because it creates an instance of blockDevice,
i.e. this is a _factory_ method)

** Takes a byte slice created by Serialize and the initial data
** Creates a new block device backend
** It will deserialize the changes and load them on top of the initial data
*/
func Deserialize(changes []blockDeviceChange, initData []byte) (blockDevice, error) {
	bd := blockDevice{}
	bd.Initialize(initData)
	for _, change := range changes {
		err := bd.WriteAt(change.Offset, change.Data[:])
		if err != nil {
			return bd, err
		}
	}
	return bd, nil
}

/* Debug helper function. Can be slow.
* Returns a copy of the backend data. It's the backened data
* with the changes layed "on top" of it.
 */
func (bd *blockDevice) Dump() []byte {
	// fmt.Printf("Dump(): length: %d \n", len(bd.Changes))
	// fmt.Printf("Dump() : %d \n", len(bd.Changes))
	c := make([]byte, len(bd.Data))
	copy(c, bd.Data)
	for _, change := range bd.Changes {
		// fmt.Printf("Dump(): k: %d, BLOCK_SIZE: %d, change.Offset: %d, change.Data: %v\n", k, BLOCK_SIZE, change.Offset, change.Data)
		offset := change.Offset
		// fmt.Printf("Dump() before copy: c: %v\n", c)
		copy(c[offset:], change.Data[:])
		// fmt.Printf("Dump()  after copy: c: %v\n", c)
	}
	return c
}
