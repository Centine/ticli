package checks

import (
	"errors"
	"fmt"

	"github.com/centine/ticli/internal/config"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Platform-independent hardcoded checks go in here.
type GenericChecker struct {
	ctx config.TicliContext
}

func NewGenericChecker(ctx config.TicliContext) Checker {
	return &GenericChecker{
		ctx: ctx,
	}
}

func (c *GenericChecker) DoSetup() error {
	// Deliberate no-op for now.
	return nil
}

func (c *GenericChecker) DoCleanup() error {
	// Deliberate no-op for now.
	return nil
}

var checks_any_platform = []Check{
	{CheckName: "Check CPU", Fn: checkCPU},
	{CheckName: "Check RAM", Fn: checkRAM},
}

func (c *GenericChecker) DoCheck() ([]CheckResult, error) {
	return generic_check_iterator(checks_any_platform, c, c.ctx)
}

func (c *GenericChecker) CheckerName() string {
	return "Baseline"
}

func checkCPU(ctx config.TicliContext) (CheckStatus, string, error) {
	coreCount, err := cpu.Counts(false) // logical cores = false
	if err != nil {
		return StatusInconclusive, "", errors.New("error retrieving CPU core count: " + err.Error())
	}
	return bool2CheckStatus(coreCount >= ctx.Config.WorkstationConfig.MinCores), fmt.Sprintf("Minimum core count of %d met", ctx.Config.WorkstationConfig.MinCores), nil
}

func checkRAM(ctx config.TicliContext) (CheckStatus, string, error) {
	// Check RAM
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return StatusInconclusive, "", errors.New("error retrieving RAM information: " + err.Error())
	}

	totalRAMInGB := float64(vmem.Total) / float64(1024*1024*1024)

	return bool2CheckStatus(int(totalRAMInGB) >= ctx.Config.WorkstationConfig.MinMemory), fmt.Sprintf("Minimum physical memory of %d GiB met", ctx.Config.WorkstationConfig.MinMemory), nil
}
