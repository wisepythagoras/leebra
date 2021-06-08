package system

import (
	"strings"

	"github.com/shirou/gopsutil/host"
)

func GetKernelArch() string {
	info, err := host.Info()

	if err != nil {
		return "x86"
	}

	return info.KernelArch
}

func GetOS() string {
	info, err := host.Info()

	if err != nil {
		return "Unknown"
	}

	return strings.Title(info.OS)
}

func GetFlavor() string {
	info, err := host.Info()

	if err != nil {
		return ""
	}

	return strings.Title(info.Platform)
}
