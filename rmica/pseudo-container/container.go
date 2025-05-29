package pseudo_container

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Egg12138/mica-OCI/rmica/logger"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// ==================== Type Definitions ====================

type Container struct {
	ID      string
	Root    string
	// Use specs-go::Spec, State to represent, following OCI-spec
	Spec    *specs.Spec
	State   *specs.State
	Config  *specs.Spec
	state   ContainerState
	initPid int
}

// ==================== Container Creation and Loading ====================

func New(root, id string, spec *specs.Spec) (*Container, error) {
	return &Container{
		ID:   id,
		Root: root,
		Spec: spec,
		Config: spec,
		state: &CreatedState{},
	}, nil
}

// Load loads an existing container
func Load(root, id string) (*Container, error) {
	containerDir := filepath.Join(root, id)
	if _, err := os.Stat(containerDir); err != nil {
		return nil, fmt.Errorf("container %s not found: %w", id, err)
	}

	return &Container{
		ID:   id,
		Root: root,
		state: &StoppedState{},
	}, nil
}

// ==================== Container Lifecycle Operations ====================

func (c *Container) Init() error {
	// Create container directory 
	// e.g.: /run/containerd/io.containerd.runtime.task.v2/moby/<id>
	containerDir := filepath.Join(c.Root, c.ID)
	if err := os.MkdirAll(containerDir, 0o700); err != nil {
		return fmt.Errorf("failed to create container directory: %w", err)
	}

	c.State = &specs.State{
		Version:     c.Spec.Version,
		ID:          c.ID,
		Status:      specs.StateCreating,
		Bundle:      c.Root,
		Annotations: c.Spec.Annotations,
	}

	return nil
}

func (c *Container) Start() error {
	return c.state.transition(&RunningState{c: c})
}

func (c *Container) Stop() error {
	return c.state.transition(&StoppedState{c: c})
}

func (c *Container) Pause() error {
	logger.Info("resume and pause hasn't implemented yet")
	return c.state.transition(&PausedState{c: c})
}

func (c *Container) Resume() error {
	logger.Info("resume and pause hasn't implemented yet")
	return c.state.transition(&RunningState{c: c})
}

func (c *Container) Destroy() error {
	return c.state.destroy()
}

// ==================== Container Status and Information ====================

func (c *Container) Status() specs.ContainerState {
	return c.state.status()
}

// ==================== Helper Functions ====================

func (c *Container) hasInit() bool {
	return c.initPid != 0
} 