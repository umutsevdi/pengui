package sys

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	cpu "github.com/shirou/gopsutil/cpu"
	_ "github.com/shirou/gopsutil/disk"
	proc "github.com/shirou/gopsutil/process"
)

type Resources struct {
	CPU          []float64
	Memory       float32
	Disk         float32
	TopProcesses [5]Process
}

type Process struct {
	PID    int32
	Name   string
	CPU    float64
	Memory float32
}

func Capture() *Resources {
	return &Resources{
		CPU:          getCPU(),
		Memory:       getMemory(),
		Disk:         getDisk(),
		TopProcesses: getProcesses(),
	}

}

func getMemory() float32 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bufio.NewScanner(file)
	available := 0
	total := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()
		arr := strings.Split(strings.ReplaceAll(raw[:len(raw)-2], " ", ""), ":")
		v, err := strconv.Atoi(arr[1])
		if err != nil {
			return 0
		}

		if arr[0] == "MemTotal" {
			total = v
		} else if arr[0] == "MemAvailable" {
			available = v
		}
		if available != 0 && total != 0 {
			break
		}
	}
	return float32(available) / float32(total)
}

func getCPU() []float64 {
	stat, _ := cpu.Percent(time.Millisecond, true)
	return stat
}

func getProcesses() [5]Process {

	processes, _ := proc.Processes()
	sort.Slice(processes, func(i, j int) bool {
		c1, _ := processes[i].CPUPercent()
		c2, _ := processes[j].CPUPercent()
		return c1 > c2
	})
	processes = processes[0:5]
	p := new([5]Process)
	for i, v := range processes {
		p[i].Name, _ = v.Name()
		p[i].PID = v.Pid
		p[i].CPU, _ = v.CPUPercent()
		p[i].Memory, _ = v.MemoryPercent()
	}
	return *p
}

func getDisk() float32 {
	return 0.0
}
