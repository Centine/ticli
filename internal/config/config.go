package config

import (
	"encoding/json"
	"log"
	"os"
)

type TicliContext struct {
	ConfigOrigin     string
	Config           ConfigStruct
	TicliDir         string
	SelectedPlatform string
}

type ConfigStruct struct {
	WorkstationConfig WorkstationConfig `json:"workstation_config"`
	K8sConfigs        []K8sConfig       `json:"k8s_configs"`
	Rules             []Rule            `json:"rules"`
}

type Rule struct {
	Type          string          `json:"type"`
	Specification json.RawMessage `json:"specification"`
}

type WorkstationConfig struct {
	MinCores  int `json:"min_cores"`
	MinMemory int `json:"min_memory_gb"`
}

type K8sConfig struct {
	Name          string `json:"name"`
	KubeconfigUrl string `json:"kubeconfig_url"`
}

func LoadConfig() (string, ConfigStruct) {
	origin := "ticli.config.json"
	file, err := os.ReadFile(origin)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	var ticliStructConfig ConfigStruct
	if err = json.Unmarshal(file, &ticliStructConfig); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	log.Println("Loaded config file")

	return origin, ticliStructConfig
}
