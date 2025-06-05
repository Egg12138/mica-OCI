package defs

type CtAct uint8

const (
	CT_ACT_CREATE CtAct = iota + 1
	CT_ACT_RUN
	CT_ACT_RESTORE
)