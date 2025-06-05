package mcs

type CpuStats struct {
	MaxCpu uint32 `json:"max_cpu,omitempty"`
}

type ClientConf struct {
	Name string `json:"name"`
	CPU  uint32 `json:"cpu"`
	ClientPath string `json:"client_path"`
	AutoBoot   bool   `json:"auto_boot"`
}

type ClientTask struct {
	Name string `json:"name"`
	Terminal bool `json:"terminal,omitempty"`
	Tty string	`json:"tty"`
}