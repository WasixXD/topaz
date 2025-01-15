package main

import (
	"log"
	"os"
	"strings"
)

func readUptime(k *KernelJson, path string) {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Printf("[!] Error on reading %s %v", path, err)
		return
	}

	parsed := strings.Split(string(content), " ")

	k.Uptime = parsed[0]
	k.IdleProcess = parsed[1]
}

func readVersion(k *KernelJson, path string) {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Printf("[!] Error on reading %s %v", path, err)
		return
	}

	k.Version = string(content)
}

func readCpu(k *KernelJson, path string) {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Printf("[!] Error on reading %s %v", path, err)
		return
	}

	// TODO: USE STRING BUILDER
	parsed := string(content)
	block, _, _ := strings.Cut(parsed, "power management")

	lines := strings.Split(block, "\n")

	for _, line := range lines {
		values := strings.Split(line, ":")
		if strings.Contains(line, "model name") {
			k.CpuName = strings.TrimSpace(values[1])
		}

		if strings.Contains(line, "cpu cores") {
			k.CpuCores = strings.TrimSpace(values[1])
		}

	}

}

func readMem(k *KernelJson, path string) {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Printf("[!] Error on reading %s %v", path, err)
		return
	}

	parsed := string(content)

	lines := strings.Split(parsed, "\n")

	line := lines[0]

	values := strings.Split(line, ":")

	k.MemTotal = strings.TrimSpace(values[1])
}
