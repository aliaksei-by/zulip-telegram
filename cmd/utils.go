package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const FILE_CONFIG = "init/config.yaml"

// Load config from file
func readConfig() error {
	content, err := os.ReadFile(FILE_CONFIG)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
