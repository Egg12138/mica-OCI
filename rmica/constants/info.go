package constants

const (
	RuntimeName = "rmica"
	SpecConfig = "config.json"
	Usage      = `Simple Pseudo-Container Runtime

A simple drop-in replacement for runc that implements basic container lifecycle management APIs
but does not actually handling any containers following the OCI specification.
`
  Root = "/run/rmica"
)
