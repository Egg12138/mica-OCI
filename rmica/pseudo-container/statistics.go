package pseudo_container

import "rmica/mcs"

type Stats struct {
	// NetworkInterface []*types.NetworkInterface
	ClientStates *mcs.ClientStats
}

func NewEmpty() Stats {
	return Stats{}
}
