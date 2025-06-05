package defs

const (
	RuntimeName = "rmica"
	SpecConfig = "config.json"
	Usage      = `Simple Pseudo-Container Runtime

A simple drop-in replacement for runc that implements basic container lifecycle management APIs
but does not actually handling any containers following the OCI specification.
`
	DefaultRootDir = "/run/rmica"
	Root = DefaultRootDir
	DefaultMicaSocket = "/var/run/micad.sock"
	SysVLogPath = "/var/log/rmica"
	DefaultLogFile = "/var/log/rmica.log"

	StateFilename    = "state.json"
	ExecFifoFilename = "exec.fifo"

	ContainerDirPerm = 0o700

	// mica-related:

	MicaConfigPath = "/etc/mica"
	MicaSocketName 		 = "mica-create.socket"
)
