package config

import (
	"encoding/json"
	"io"
	"log"
	"ova-checklist-api/internal/utils"
)

type SimpleConfig struct {
	Environment string `json:"environment,omitempty"`
}

// jsonReadingVisitor implements utils.FileReadingVisitor
type jsonReadingVisitor struct {
	Config SimpleConfig
}

func (j *jsonReadingVisitor) OnOpenFail(filename string, err error) {
	log.Panicf("Unable to open configuration file: %s due to an error: %v", filename, err)
}

func (j *jsonReadingVisitor) OnOpenSuccess(filename string, reader io.Reader) {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&j.Config); err != nil {
		log.Panicf("Unable to decode SimpleConfig: %s due to an error: %v", filename, err)
	}
	log.Printf("Configuration file: %s was decoded successfully", filename)
}

func (j *jsonReadingVisitor) OnCloseFail(filename string, err error) {
	log.Panicf("Unable to close configuration file: %s due to an error: %v", filename, err)
}

func (j *jsonReadingVisitor) OnCloseSuccess(filename string) {
	log.Printf("Configuration file: %s was processed successfully", filename)
}

func OpenSimpleConfig(filename string) *SimpleConfig {
	var fs utils.FileSystemOpeningStrategy
	var visitor jsonReadingVisitor
	utils.ReadFiles(&fs, &visitor, filename)
	return &visitor.Config
}
