package defs

const (
	RuntimeName = "rmica"
	SpecConfig = "config.json"
	Usage      = `Simple Pseudo-Container Runtime

A simple drop-in replacement for runc that implements basic container lifecycle management APIs
but does not actually handling any containers following the OCI specification.
`
	// DefaultRootDir = "/run/rmica"
	DefaultRootDir = "/tmp/run/rmica"

	Root = DefaultRootDir
	// DefaultMicaSocket = "/var/run/micad.sock"

	DefaultMicaSocket = "/tmp/mica/mica-create.socket"
	SysVLogPath = "/var/log/rmica" // permission check
	DefaultLogFile = "/var/tmp/rmica.log"

	StateFilename    = "state.json"
	ExecFifoFilename = "exec.fifo"

	ContainerDirPerm = 0o700

	// mica-related:

	MicaConfigPath = "/etc/mica"
	MicaSocketName 		 = "mica-create.socket"
	// prefix of annotaion fields belonging to mica
	MicaAnnotationPrefix = "org.openeuler.mica."
)
