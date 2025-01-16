package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const PROC_PATH string = "/proc"
const PROC_UPTIME string = "/proc/uptime"
const PROC_VERSION string = "/proc/version"
const PROC_CPU string = "/proc/cpuinfo"
const PROC_MEM string = "/proc/meminfo"

type Process struct {
	State string `json:"state"`
	Name  string `json:"name"`
	Pid   string `json:"pid"`
}

type ProcessesCache struct {
	processes map[string]Process
}

type KernelJson struct {
	Uptime      string `json:"uptime"`
	IdleProcess string `json:"idle_process"`
	Version     string `json:"version"`
	CpuName     string `json:"cpu_name"`
	CpuCores    string `json:"cpu_cores"`
	MemTotal    string `json:"mem_total"`
}

type ProcessJson struct {
	ProcessTotal int       `json:"process_total"`
	TopProcesses []Process `json:"top_processes"`
}

type Runner struct {
	data   interface{}
	mu     sync.Mutex
	ticker time.Ticker
}

type Master struct {
	runners map[string]*Runner
	mu      sync.RWMutex
	cache   ProcessesCache
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

func (m *Master) resetCache() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache = ProcessesCache{processes: make(map[string]Process)}
}

func (m *Master) AddProcess(name string, p Process) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache.processes[name] = p
}

func (r *Runner) SetProcesses(rang int) {
	r.mu.Lock()
	pStruct, _ := r.data.(ProcessJson)
	global.mu.Lock()
	for _, v := range global.cache.processes {
		if len(pStruct.TopProcesses) == rang {
			break
		}
		pStruct.TopProcesses = append(pStruct.TopProcesses, v)
	}
	global.mu.Unlock()

	r.data = pStruct
	r.mu.Unlock()
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
		for file_path, function := range commands {
			// evil stuff?
			kernelStruct, _ := r.data.(KernelJson)
			function(&kernelStruct, file_path)
			r.data = kernelStruct
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
		global.resetCache()

		if err != nil {
			log.Fatalf("Error on reading /proc: %v", err)
		}

		for _, file := range processes {
			file_name := file.Name()
			if _, err := strconv.Atoi(file_name); err == nil {
				count++
			}
			if !file.IsDir() {
				continue
			}

			file_path := strings.Builder{}
			file_path.WriteString(PROC_PATH + "/" + file_name + "/status")
			content, err := os.ReadFile(file_path.String())

			if err != nil {
				continue
			}

			lines := strings.Split(string(content), TOKEN_NEW_LINE)
			name := strings.TrimSpace(strings.Split(lines[0], TOKEN_COMMA)[1])
			status := strings.TrimSpace(strings.Split(lines[2], TOKEN_COMMA)[1])
			p := Process{Name: name, State: status, Pid: file_name}

			r.data = ProcessJson{ProcessTotal: count}
			global.AddProcess(name, p)

		}
		r.mu.Unlock()
	}
}

func StartRunners() {
	global = Master{runners: make(map[string]*Runner), mu: sync.RWMutex{}, cache: ProcessesCache{
		processes: make(map[string]Process),
	}}

	r1 := Runner{data: KernelJson{}, mu: sync.Mutex{}, ticker: *time.NewTicker(5 * time.Second)}

	// the idea is that we can create functions to parse especific files
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
