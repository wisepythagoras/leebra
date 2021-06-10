package system

import (
	"strings"

	"github.com/shirou/gopsutil/host"
)

// GetKernelArch returns the architecture of the host system.
func GetKernelArch() string {
	info, err := host.Info()

	if err != nil {
		return "x86"
	}

	return info.KernelArch
}

// GetOS returns the operating system (Linux, Windows, etc).
func GetOS() string {
	info, err := host.Info()

	if err != nil {
		return "Unknown"
	}

	return strings.Title(info.OS)
}

// GetFlavor should return "Ubuntu", or whatever the system is.
func GetFlavor() string {
	info, err := host.Info()

	if err != nil {
		return ""
	}

	return strings.Title(info.Platform)
}
