package esx

import "github.com/vmware/govmomi/vim25/types"

// HardwareVersion is the default VirtualMachine.Config.Version
var HardwareVersion = "vmx-13"

// Setting is captured from ESX's HostSystem.configManager.advancedOption
// Capture method:
//   govc object.collect -s -dump $(govc object.collect -s HostSystem:ha-host configManager.advancedOption) setting
var Setting = []types.BaseOptionValue{
	// This list is currently pruned to include a single option for testing
	&types.OptionValue{
		Key:   "Config.HostAgent.log.level",
		Value: "info",
	},
}
