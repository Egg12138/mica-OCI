// this file is inspired by the runc libcontainer
package pseudo_container

import (
	"fmt"
	"os"

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
		s.c.cstate = to
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
		r.c.cstate = to
		return nil
	case *PausedState:
		r.c.cstate = to
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
		p.c.cstate = to
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
		c.c.cstate = to
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
// 1. remove host container dir
// 2. remove the task or shut down the clientOS
// TODO: we have to wrap task config into container spec files 
func destroy(c *Container) error {
	// BUG: container informatino and task config is separated, here is just for demo
	res := communication.Send2micaCmd("rm", c.Id())
	if res != "" {
		// TODO: handle response
		logger.Debugf("destroy container %s: %s", c.Id(), res)
		logger.Fprintf("destroy container %s: %s", c.Id(), res)
	}

	// Remove container directory
	if err := os.RemoveAll(c.StateDir()); err != nil {
		return fmt.Errorf("failed to remove container directory: %w", err)
	}

	if c.config.Hooks != nil {
		s := c.OCIState()
		logger.Fprintf("get OCI state: %s [%s]", s.Status, s.Bundle)
		logger.Fprintf("we do not run hook in demo")
		s.Status = specs.StateStopped
	}
	c.cstate = &StoppedState{c: c}
	return nil
}


func runHook(hooks *specs.Hooks, name string, state *specs.State) error {
	logger.Fprintf("run hook %s: %s", name, state)
	logger.Debugf("run hook %s: %s", name, state)
	return nil
}