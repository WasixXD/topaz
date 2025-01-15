package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const PROC_PATH string = "/proc"
const PROC_UPTIME string = "/proc/uptime"
const PROC_VERSION string = "/proc/version"
const PROC_CPU string = "/proc/cpuinfo"
const PROC_MEM string = "/proc/meminfo"

type KernelJson struct {
	Uptime      string `json:"uptime"`
	IdleProcess string `json:"idle_process"`
	Version     string `json:"version"`
	CpuName     string `json:"cpu_name"`
	CpuCores    string `json:"cpu_cores"`
	MemTotal    string `json:"mem_total"`
}

type ProcessJson struct {
	ProcessTotal int `json:"process_total"`
}

type Runner struct {
	data   interface{}
	mu     sync.Mutex
	ticker time.Ticker
}

type Master struct {
	runners map[string]*Runner
	mu      sync.RWMutex
}

var global Master

func (m *Master) GetRunner(name string) *Runner {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.runners[name]
}

func (m *Master) AddRunner(name string, r *Runner) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.runners[name] = r
}

func (r *Runner) GetData() interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.data
}

func (r *Runner) KernelRunner(commands map[string]func(*KernelJson, string)) {
	for {
		<-r.ticker.C

		r.mu.Lock()
		for i, v := range commands {
			// evil stuff?
			field, _ := r.data.(KernelJson)
			v(&field, i)
			r.data = field
		}

		r.mu.Unlock()
	}
}

func (r *Runner) ProcessRunner() {
	for {
		<-r.ticker.C
		count := 0

		r.mu.Lock()

		processes, err := os.ReadDir(PROC_PATH)

		if err != nil {
			log.Fatalf("Error on reading /proc: %v", err)
		}

		for _, file := range processes {
			if _, err := strconv.Atoi(file.Name()); err == nil {
				count++
			}
		}

		r.data = ProcessJson{ProcessTotal: count}

		r.mu.Unlock()
	}
}

func StartRunners() {
	global = Master{runners: make(map[string]*Runner), mu: sync.RWMutex{}}

	r1 := Runner{data: KernelJson{}, mu: sync.Mutex{}, ticker: *time.NewTicker(5 * time.Second)}

	// TODO: make a nice interfaces for this
	commands := make(map[string]func(*KernelJson, string))
	commands[PROC_UPTIME] = readUptime
	commands[PROC_VERSION] = readVersion
	commands[PROC_CPU] = readCpu
	commands[PROC_MEM] = readMem

	r2 := Runner{data: ProcessJson{}, mu: sync.Mutex{}, ticker: *time.NewTicker(2 * time.Second)}

	go r1.KernelRunner(commands)
	go r2.ProcessRunner()

	global.AddRunner("kernel", &r1)
	global.AddRunner("process", &r2)
}
