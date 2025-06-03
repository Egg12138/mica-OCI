// this file is inspired by the runc libcontainer
package pseudo_container

import (
	"fmt"
	"os"
	"path/filepath"

	"rmica/communication"
	"rmica/logger"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// ==================== Type Definitions ====================

// runc OOP state machine;
// TODO: replace with spec-go state machine
type ContainerState interface {
	transition(ContainerState) error
	destroy() error
	status() specs.ContainerState
}

// ==================== Error Handling ====================

type StateTransitionError struct {
	From specs.ContainerState
	To   specs.ContainerState
}

func (s *StateTransitionError) Error() string {
	errMsg := fmt.Sprintf("invalid state transition from %s to %s", s.From, s.To)
	logger.Errorf("%s", errMsg)
	return errMsg
}

func newStateTransitionError(from, to ContainerState) error {
	return &StateTransitionError{
		From: from.status(),
		To:   to.status(),
	}
}

// ==================== Switch States ====================
// Creating, Created, Stopped, Running
// in runc: Created, Runnig, Paused, Stopped
// <State>.transition(to):
// 

type StoppedState struct {
	c *Container
}

func (s *StoppedState) status() specs.ContainerState {
	return specs.StateStopped
}

func (s *StoppedState) transition(to ContainerState) error {
	switch to.(type) {
	case *RunningState, *RestoredState:
		s.c.state = to
		return nil
	case *StoppedState:
		return nil
	}
	return newStateTransitionError(s, to)
}

func (s *StoppedState) destroy() error {
	return destroy(s.c)
}

type RunningState struct {
	c *Container
}

func (r *RunningState) status() specs.ContainerState {
	return specs.StateRunning
}

func (r *RunningState) transition(to ContainerState) error {
	switch to.(type) {
	case *StoppedState:
		if r.c.hasInit() {
			return fmt.Errorf("container is running")
		}
		r.c.state = to
		return nil
	case *PausedState:
		r.c.state = to
		return nil
	case *RunningState:
		return nil
	}
	return newStateTransitionError(r, to)
}

func (r *RunningState) destroy() error {
	if r.c.hasInit() {
		return fmt.Errorf("container is running")
	}
	return destroy(r.c)
}

// Not defined in spec-go but in runc libcontainer
type PausedState struct {
	c *Container
}

func (p *PausedState) status() specs.ContainerState {
	return specs.StateRunning // Using Running state as Paused is not in specs-go
}

func (p *PausedState) transition(to ContainerState) error {
	switch to.(type) {
	case *RunningState, *StoppedState:
		p.c.state = to
		return nil
	case *PausedState:
		return nil
	}
	return newStateTransitionError(p, to)
}

func (p *PausedState) destroy() error {
	if p.c.hasInit() {
		return fmt.Errorf("container is paused")
	}
	return destroy(p.c)
}

// CreatedState represents a container in created state
type CreatedState struct {
	c *Container
}

func (c *CreatedState) status() specs.ContainerState {
	return specs.StateCreated
}

func (c *CreatedState) transition(to ContainerState) error {
	switch to.(type) {
	case *RunningState, *PausedState, *StoppedState:
		c.c.state = to
		logger.Infof("[%s] %s -> %s", c.c.Id(), c.status(), to.status())
		return nil
	case *CreatedState:
		return nil
	}
	return newStateTransitionError(c, to)
}

func (c *CreatedState) destroy() error {
	return destroy(c.c)
}

// RestoredState represents a container that has been restored from a checkpoint
type RestoredState struct {
	c *Container
}

func (r *RestoredState) status() specs.ContainerState {
	return specs.StateRunning
}

func (r *RestoredState) transition(to ContainerState) error {
	switch to.(type) {
	case *StoppedState, *RunningState:
		return nil
	}
	return newStateTransitionError(r, to)
}

func (r *RestoredState) destroy() error {
	return destroy(r.c)
}

// ==================== Helper Functions ====================

// Helper function to destroy container resources
// TODO mica destroy
// TODO preHook and postHook
func destroy(c *Container) error {
	res := communication.Send2mica("destrory")
	if res != "" {
		// TODO: handle response
	}

	// Remove container directory
	containerDir := filepath.Join(c.root, c.Id())
	if err := os.RemoveAll(containerDir); err != nil {
		return fmt.Errorf("failed to remove container directory: %w", err)
	}
	c.state = &StoppedState{c: c}
	return nil
}
