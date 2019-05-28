package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

func loadConfigFile(cfgPath string, ptr interface{}) error {
	if ptr == nil {
		return fmt.Errorf("ptr of type (%T) is nil", ptr)
	}
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to open config %s, err: %v", cfgPath, err)
	}
	if err := json.Unmarshal(data, ptr); err != nil {
		return fmt.Errorf("failed to unmarshal config %s, err: %v", cfgPath, err)
	}
	return nil
}

func LoadConfigFileOrPanic(cfgPath string, ptr interface{}) {
	if err := loadConfigFile(cfgPath, ptr); err != nil {
		log.Fatal().Err(err).Msgf("failed to load config %s", cfgPath)
	}
}
