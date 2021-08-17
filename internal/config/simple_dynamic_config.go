package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type SimpleDynamicConfig struct {
	EndpointSet []string `json:"endpoint_set"`
}

type SimpleDynamicConfigReader struct {
	filename string
	checkPeriod time.Duration
	config SimpleDynamicConfig
}

func NewSimpleDynamicConfigReader(filename string, checkPeriod time.Duration) SimpleDynamicConfigReader {
	return SimpleDynamicConfigReader{
		filename: filename,
		checkPeriod: checkPeriod,
	}
}

func closeFile(file *os.File, filename string) {
	if err := file.Close(); err != nil {
		log.Printf("Unable to close file: %v properly due to an error: %v", filename, err)
	}
}

func (r *SimpleDynamicConfigReader) Update() {
	for {
		func() {
			file, err := os.Open(r.filename)
			if err != nil {
				log.Printf("Unable to open file: %v due to an error: %v", r.filename, err)
				return
			}
			defer closeFile(file, r.filename)
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&r.config); err != nil {
				log.Printf("Unable to parse config: %v due to an error: %v", r.filename, err)
			}
		}()
		time.Sleep(r.checkPeriod)
	}
}
