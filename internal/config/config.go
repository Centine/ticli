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
	K8sConfigs      []K8sConfig      `json:"k8s_configs"`
	ToolVersions    []ToolVersion    `json:"tool_versions"`
	RuleBasedChecks []RuleBasedCheck `json:"rule_based_checks"`
}

type K8sConfig struct {
	Name          string `json:"name"`
	KubeconfigUrl string `json:"kubeconfig_url"`
}

type ToolVersion struct {
	Tool       string `json:"tool"`
	Platform   string `json:"platform"`
	MinVersion string `json:"min_version"`
}

type RuleBasedCheck struct {
	Type                  string   `json:"type"`
	PlatformCompatibility string   `json:"platform_compatibility"`
	Prerequisites         []string `json:"prereqs"`
	Command               string   `json:"command"`
	OutputRegex           string   `json:"output_regex"`
	AllowableExitCodes    []int    `json:"allowable_exit_codes"`
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

	// Now you can use the `config` variable throughout your application.
	// To access fields, you would need to do type assertions.

	log.Println("Loaded config file")

	return origin, ticliStructConfig
}
