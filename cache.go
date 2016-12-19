package wurflgo

import (
	"os"
	"encoding/gob"
	"bufio"
)

func Read(gobFile string) *Repository {
	decodeFile, err := os.Open(gobFile)
	if err != nil {
		return nil
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Place to decode into
	devices := make(map[string]*Device)

	// Decode -- We need to pass a pointer otherwise accounts2 isn't modified
	decoder.Decode(&devices)

	if devices == nil || len(devices) <= 10 {
		return nil
	}

	repo := &Repository {
		devices: devices,
	}
	repo.Initialize()
	return repo
}

func (r *Repository) Save(gobFile string) error {
	// Create a file for IO
	encodeFile, err := os.Create(gobFile)
	if err != nil {
		return err
	}
	defer encodeFile.Close()
	encoder := gob.NewEncoder(bufio.NewWriter(encodeFile))

	// Write to the file
	if err := encoder.Encode(r.devices); err != nil {
		os.Remove(gobFile)
		return err
	}
	return nil
}