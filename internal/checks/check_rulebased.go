package checks

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/centine/ticli/internal/config"
)

// Rule-based checks go in here.

type RulebasedChecker struct {
	ctx config.TicliContext
}

func newRuleBasedChecker(ctx config.TicliContext) Checker {
	return &RulebasedChecker{
		ctx: ctx,
	}
}

func (c *RulebasedChecker) CheckerName() string {
	return "Rule"
}

func (c *RulebasedChecker) DoSetup() error {
	// Deliberate no-op for now.
	return nil
}

func (c *RulebasedChecker) DoCleanup() error {
	return nil
}

func (c *RulebasedChecker) DoCheck() ([]CheckResult, error) {
	checkResults := make([]CheckResult, 0)

	log.Println("Checking rule-based checks")

	for _, rule := range c.ctx.Config.Rules {
		switch rule.Type {
		case "tcp_connectivity":
			var check tcpConnectivityCheck
			if err := json.Unmarshal(rule.Specification, &check); err != nil {
				log.Fatalf("configuration error: %v", err)
			}
			checkResults = append(checkResults, c.checkTCPConnectivity(check))
		case "env_var":
			var check environmentVariableCheck
			if err := json.Unmarshal(rule.Specification, &check); err != nil {
				log.Fatalf("configuration error: %v", err)
			}
			checkResults = append(checkResults, c.checkEnvVar(check))
		case "http_connectivity":
			var check httpConnectivityCheck
			if err := json.Unmarshal(rule.Specification, &check); err != nil {
				log.Fatalf("configuration error: %v", err)
			}
			// TODO now use check
		case "validate_tool":
			var check toolValidationCheck
			if err := json.Unmarshal(rule.Specification, &check); err != nil {
				log.Fatalf("configuration error: %v", err)
			}
			checkResults = append(checkResults, c.checkTool(check)...)
		default:
			// unknown type
		}
	}

	return checkResults, nil // TODO: return errors
}

type toolValidationCheck struct {
	Tool                    string              `json:"tool"`
	Platform                string              `json:"platform"`
	VersionDetectionCommand string              `json:"version_detection_command"`
	VersionRegex            string              `json:"version_regex"`
	MinVersions             []map[string]string `json:"min_versions"`
}

type tcpConnectivityCheck struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type httpConnectivityCheck struct {
	URL                        string            `json:"url"`
	Method                     string            `json:"method"`
	Headers                    map[string]string `json:"headers"`
	ResponseAcceptableStatuses []int             `json:"response_acceptable_status_codes"`
}

type environmentVariableCheck struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

func (gc *RulebasedChecker) checkTCPConnectivity(check tcpConnectivityCheck) CheckResult {
	address := fmt.Sprintf("%s:%d", check.Host, check.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)

	if err != nil {
		return CheckResult{
			CheckName: "TCP Connectivity",
			Status:    StatusFail,
			Notes:     fmt.Sprintf("Failed to connect to %s, error %v", address, err),
			Source:    gc.CheckerName(),
		}
	}

	conn.Close()

	return CheckResult{
		CheckName: "TCP Connectivity",
		Status:    StatusSuccess,
		Notes:     fmt.Sprintf("Successfully connected to %s", address),
		Source:    gc.CheckerName(),
	}
}

func (gc *RulebasedChecker) checkEnvVar(check environmentVariableCheck) CheckResult {
	val, exists := os.LookupEnv(check.Name)
	checkName := "Env " + check.Name
	if !exists {
		cr := newCheckResult(checkName, StatusWarning, fmt.Sprintf("Env %s not set", check.Name), withSource(gc.CheckerName()))
		return cr
	} else if check.Prefix != "" && !strings.HasPrefix(val, check.Prefix) {
		cr := newCheckResult(checkName, StatusWarning, fmt.Sprintf("Env %s set to %s, unexpected prefix %s", check.Name, val, check.Prefix), withSource(gc.CheckerName()))
		return cr
	} else {
		cr := newCheckResult(checkName, StatusSuccess, fmt.Sprintf("Env %s exists (%s)", check.Name, val), withSource(gc.CheckerName()))
		return cr
	}
}

func (gc *RulebasedChecker) checkTool(tool toolValidationCheck) []CheckResult {
	var results []CheckResult

	// Check if tool is installed
	path, err := exec.LookPath(tool.Tool)
	if err != nil {
		cr := newCheckResult(tool.Tool, StatusFail, "Tool not found in PATH", withSource(gc.CheckerName()))
		results = append(results, cr)
		return results
	}

	foundCr := newCheckResult(tool.Tool, StatusSuccess, fmt.Sprintf("Tool found in path (%s)", path), withSource(gc.CheckerName()))
	results = append(results, foundCr)

	// Execute version detection command
	cmdParts := strings.Fields(tool.VersionDetectionCommand)
	cmd := exec.Command(path, cmdParts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil && (cmd.ProcessState.ExitCode() > 1) {
		cr := newCheckResult(tool.Tool, StatusFail, "Error executing version detection command: "+err.Error(), withSource(gc.CheckerName()))
		results = append(results, cr)
		return results
	}

	// Extract version using regex - ignore error codes
	re := regexp.MustCompile("(?i)(?m)" + tool.VersionRegex)
	matches := re.FindAllStringSubmatch(string(output), -1)
	if matches == nil {
		cr := newCheckResult(tool.Tool, StatusFail, "Unable to find any version strings in output, using regex "+tool.VersionRegex+" on output "+string(output), withSource(gc.CheckerName()))
		results = append(results, cr)
		return results
	}

	matchedVersions := make(map[string]*semver.Version)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		versionKey := match[1]
		versionValue := match[2]
		semVerVersion, err := semver.NewVersion(versionValue)
		if err != nil {
			cr := newCheckResult(tool.Tool, StatusFail, "Error parsing version: "+err.Error(), withSource(gc.CheckerName()))
			results = append(results, cr)
		}
		matchedVersions[versionKey] = semVerVersion
	}

	for _, minVersion := range tool.MinVersions {
		log.Printf("Checking version key %v for tool %s\n", minVersion, tool.Tool)

		for minVersionKey, requiredVersionStr := range minVersion {
			outputMatch := matchedVersions[minVersionKey]
			if outputMatch != nil {
				requiredVersion, err := semver.NewVersion(requiredVersionStr)
				if err != nil {
					log.Fatalf("Unable to parse version key %s for tool %s", tool.Tool, requiredVersionStr)
				}
				if outputMatch.LessThan(requiredVersion) {
					cr := newCheckResult(
						tool.Tool+" "+strings.ToLower(minVersionKey),
						StatusFail,
						fmt.Sprintf("Version %v of %v is less than minimum required %v", outputMatch, minVersionKey, requiredVersionStr),
						withSource(gc.CheckerName()))
					results = append(results, cr)
				} else {
					cr := newCheckResult(
						tool.Tool+" "+strings.ToLower(minVersionKey),
						StatusSuccess,
						"Version meets requirement",
						withSource(gc.CheckerName()))
					results = append(results, cr)
				}
			} else {
				cr := newCheckResult(
					tool.Tool+" "+strings.ToLower(minVersionKey),
					StatusFail,
					fmt.Sprintf("Unable to find version key %v in output %v", minVersionKey, string(output)),
					withSource(gc.CheckerName()))
				results = append(results, cr)
			}
		}
	}
	return results
}
