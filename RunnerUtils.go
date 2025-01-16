package main

import (
	"log"
	"os"
	"strings"
)

const (
	TOKEN_MODEL_NAME       = "model name"
	TOKEN_CPU_CORES        = "cpu cores"
	TOKEN_POWER_MANAGEMENT = "power management"
	TOKEN_SPACE            = " "
	TOKEN_NEW_LINE         = "\n"
	TOKEN_COMMA            = ":"
)

func fileContent(path string) []byte {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("[!] Error on reading %s %v", path, err)
	}

	return content
}

func readUptime(k *KernelJson, path string) {
	content := fileContent(path)

	parsed := strings.Split(string(content), TOKEN_SPACE)

	k.Uptime = parsed[0]
	k.IdleProcess = parsed[1]
}

func readVersion(k *KernelJson, path string) {
	content := fileContent(path)
	k.Version = string(content)
}

func readCpu(k *KernelJson, path string) {
	content := fileContent(path)

	// TODO: USE STRING BUILDER
	parsed := string(content)
	block, _, _ := strings.Cut(parsed, TOKEN_POWER_MANAGEMENT)

	lines := strings.Split(block, TOKEN_NEW_LINE)

	for _, line := range lines {
		values := strings.Split(line, TOKEN_COMMA)
		if strings.Contains(line, TOKEN_MODEL_NAME) {
			k.CpuName = strings.TrimSpace(values[1])
		}

		if strings.Contains(line, TOKEN_CPU_CORES) {
			k.CpuCores = strings.TrimSpace(values[1])
		}

	}

}

func readMem(k *KernelJson, path string) {
	content := fileContent(path)

	parsed := string(content)

	lines := strings.Split(parsed, TOKEN_NEW_LINE)
	line := lines[0]
	values := strings.Split(line, TOKEN_COMMA)
	k.MemTotal = strings.TrimSpace(values[1])
}
