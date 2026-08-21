

## Assignment

Create an in-memory block device backend in Go that manages data (bytes),
tracks changes, and allows serialization and deserialization.

## Motivation

When we start sandboxes, each sandbox gets a filesystem. When you pause the
sandbox, we don’t want to save the whole filesystem; we only want to save the
changes made by each sandbox on top of the “base” filesystem these sandboxes
share—essentially doing a copy-on-write for each sandbox.

We implemented this using Network Block Device, but the actual data is stored in
a simple structure: a block device backend that exposes methods for reading and
writing bytes at specific offsets, as well as methods for serializing the changes.

## Details

Block device backend allows data to be read and written at specific byte offsets.

The length of data to read or write and the offset must be multiples of the “block
size”. Otherwise, we would need to individually keep track of every byte written,
which would cause significant overhead. Use a block size of 4096 bytes.

We also want a way to serialize and deserialize the backend. Serialization should
create a byte slice containing only the changes made by writing to the device, and
deserialization creates a new block device backend from these changes and the
initial “base” data.

## Specification

The result is a struct that stores the data in memory and the
following functions/methods:

The struct is initialized with:

* initial data

** This is the “base” representing the shared filesystem

** It should not be modified

** Assume the initial data is divisible by the “block size”

** Assume that we won’t be reading/writing outside of the length of the data

The functions/methods we want to have:

* ReadAt

** Takes a byte offset and a data byte slice with a length corresponding to the number of bytes we want to read. The data from the device backend will be read into the data

** Assume that the offset and length to read will always be a multiple of the “block size”

* WriteAt

** It takes a byte offset and a data byte slice containing the data we want to write. The content of the data will be written into the device backend

** Assume that the offset and length to write will always be a multiple of the “block size”

* Serialize

** Return a byte slice with serialized changes made to the block device backend

** The format should have the smallest possible overhead

* Deserialize

** Takes a byte slice created by Serialize and the initial data

** Creates a new block device backend

** It will deserialize the changes and load them on top of the initial data

