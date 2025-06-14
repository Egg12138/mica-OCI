package utils

import "errors"

var (
	ErrEmptyID = errors.New("container id cannot be empty")
	ErrExist          = errors.New("container with given ID already exists")
	ErrInvalidID      = errors.New("invalid container ID format")
	ErrNotExist       = errors.New("container does not exist")
	ErrPaused         = errors.New("container paused")
	ErrRunning        = errors.New("container still running")
	ErrNotRunning     = errors.New("container not running")
	ErrNotPaused      = errors.New("container not paused")
	ErrCgroupNotExist = errors.New("cgroup not exist")
	ErrDebug		   		= errors.New("debug mode")
	ErrNotImplemented = errors.New("not implemented")
)

