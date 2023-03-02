package sys

import (
	"time"

	host "github.com/shirou/gopsutil/host"
)

type Host struct {
	Hostname string
	Os       string
	Platform string
	Kernel   string
	Uptime   time.Duration
	Users    []string
}

func GetHostInfo() *Host {

	info, _ := host.Info()
	uptime, _ := host.Uptime()

	return &Host{
		Hostname: info.Hostname,
		Os:       info.OS,
		Platform: info.Platform,
		Kernel:   info.KernelVersion,
		Uptime:   time.Duration(uptime),
		Users:    getUserNames(),
	}
}

func getUserNames() []string {
	u, _ := host.Users()
	arr := make([]string, len(u))
	for i, v := range u {
		arr[i] = v.User
	}
	return arr
}
