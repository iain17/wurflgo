package wurflgo

import (
	"os"
	"encoding/gob"
	"fmt"
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

	if devices == nil {
		return nil
	}
	fmt.Println(devices)
	return &Repository {
		devices: devices,
	}
}

func (r *Repository) Save(gobFile string) error {
	// Create a file for IO
	encodeFile, err := os.Create(gobFile)
	if err != nil {
		return err
	}
	defer encodeFile.Close()
	encoder := gob.NewEncoder(encodeFile)

	// Write to the file
	if err := encoder.Encode(r.devices); err != nil {
		os.Remove(gobFile)
		return err
	}
	return nil
}