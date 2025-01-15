package main

import (
	"sync"
	"time"
)

const PROC_PATH string = "/proc"

type KernelJson struct {
	Uptime      string `json:"uptime"`
	IdleProcess string `json:"idle_process"`
	Version     string `json:"version"`
	CpuName     string `json:"cpu_name"`
	CpuCores    string `json:"cpu_cores"`
	MemTotal    string `json:"mem_total"`
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

func (m *Master) getRunner(name string) *Runner {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.runners[name]
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

func StartRunners() {
	global = Master{runners: make(map[string]*Runner), mu: sync.RWMutex{}}

	r1 := Runner{data: KernelJson{}, mu: sync.Mutex{}, ticker: *time.NewTicker(5 * time.Second)}

	commands := make(map[string]func(*KernelJson, string))
	commands["/proc/uptime"] = readUptime
	commands["/proc/version"] = readVersion
	commands["/proc/cpuinfo"] = readCpu
	commands["/proc/meminfo"] = readMem

	go r1.KernelRunner(commands)

	global.mu.Lock()
	global.runners["kernel"] = &r1
	global.mu.Unlock()
}
