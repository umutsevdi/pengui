package sys

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	proc "github.com/shirou/gopsutil/process"
)

type Resources struct {
	CPU          []float64  `json:"cpu"`
	Memory       float32    `json:"mem"`
	Disk         float32    `json:"disk"`
	TopProcesses [5]Process `json:"proc"`
}

type Process struct {
	PID    int32   `json:"pid"`
	Name   string  `json:"name"`
	CPU    float64 `json:"cpu"`
	Memory float32 `json:"mem"`
}

var resource Resources

func StreamResource() *Resources {
	return &resource
}

func init() {
	log.Println("Capturing started")
	go func() {
		for {
			c, _ := cpu.Percent(250*time.Millisecond, true)
			resource.CPU = c
			resource.fetchMemory()
			resource.fetchDisk()
			resource.fetchProcesses()
		}
	}()
}

func (r *Resources) fetchMemory() {
	////m, _ := mem.VirtualMemory()
	////return float32(m.UsedPercent)
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
			r.Memory = 0
			return
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
	r.Memory = float32(available) / float32(total)
}

func (r *Resources) fetchProcesses() {

	processes, _ := proc.Processes()
	sort.Slice(processes, func(i, j int) bool {
		c1, _ := processes[i].CPUPercent()
		c2, _ := processes[j].CPUPercent()
		return c1 < c2
	})
	processes = processes[len(processes)-5:]
	for i, v := range processes {
		n, err := v.Name()
		if err == nil {
			r.TopProcesses[i].Name = n
		} else {
			r.TopProcesses[i].Name = ""
		}
		r.TopProcesses[i].PID = v.Pid
		c, err := v.CPUPercent()
		if err == nil {
			r.TopProcesses[i].CPU = c
		} else {
			r.TopProcesses[i].CPU = 0
		}

		m, err := v.MemoryPercent()
		if err == nil {
			r.TopProcesses[i].Memory = m
		} else {
			r.TopProcesses[i].Memory = 0
		}

	}
}

func (r *Resources) fetchDisk() {
	d, err := disk.Usage("/")
	if err == nil {
		r.Disk = float32(d.UsedPercent)
	}
}
