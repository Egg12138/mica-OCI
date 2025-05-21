/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package micashim

import (
	"context"
	"fmt"
	"io"
	"os"

	taskAPI "github.com/containerd/containerd/api/runtime/task/v2"
	apitypes "github.com/containerd/containerd/api/types"
	ptypes "github.com/containerd/containerd/v2/pkg/protobuf/types"
	"github.com/containerd/containerd/v2/pkg/shim"
	"github.com/containerd/containerd/v2/pkg/shutdown"
	"github.com/containerd/containerd/v2/plugins"
	"github.com/containerd/errdefs"
	"github.com/containerd/plugin"
	"github.com/containerd/plugin/registry"
	"github.com/containerd/ttrpc"
)

const (
	MicaShim = "io.containerd.mica.v1"
)

func init() {
	registry.Register(&plugin.Registration{
		Type: plugins.TTRPCPlugin,
		ID:   "task",
		Requires: []plugin.Type{
			plugins.EventPlugin,
			plugins.InternalPlugin,
		},
		InitFn: func(ic *plugin.InitContext) (interface{}, error) {
			pp, err := ic.GetByID(plugins.EventPlugin, "publisher")
			if err != nil {
				return nil, err
			}
			ss, err := ic.GetByID(plugins.InternalPlugin, "shutdown")
			if err != nil {
				return nil, err
			}
			return newTaskService(ic.Context, pp.(shim.Publisher), ss.(shutdown.Service))
		},
	})
}

func NewManager(name string) shim.Manager {
	return manager{name: name}
}

type manager struct {
	name string
}

func (m manager) Name() string {
	return m.name
}

func (m manager) Start(ctx context.Context, id string, opts shim.StartOpts) (shim.BootstrapParams, error) {
	fmt.Printf("hello Start\n")
	return shim.BootstrapParams{}, errdefs.ErrNotImplemented
}

func (m manager) Stop(ctx context.Context, id string) (shim.StopStatus, error) {
	fmt.Printf("hello Stop\n")
	return shim.StopStatus{}, errdefs.ErrNotImplemented
}

func (m manager) Info(ctx context.Context, optionsR io.Reader) (*apitypes.RuntimeInfo, error) {
	fmt.Printf("hello Info\n")
	info := &apitypes.RuntimeInfo{
		Name: MicaShim,
		Version: &apitypes.RuntimeVersion{
			Version: "v1.0.0",
		},
	}
	return info, nil
}

func newTaskService(ctx context.Context, publisher shim.Publisher, sd shutdown.Service) (taskAPI.TaskService, error) {
	return &micaTaskService{}, nil
}

var (
	_ = shim.TTRPCService(&micaTaskService{})
)

type micaTaskService struct {
}

// RegisterTTRPC allows TTRPC services to be registered with the underlying server
func (s *micaTaskService) RegisterTTRPC(server *ttrpc.Server) error {
	taskAPI.RegisterTaskService(server, s)
	return nil
}

// Create a new container
func (s *micaTaskService) Create(ctx context.Context, r *taskAPI.CreateTaskRequest) (_ *taskAPI.CreateTaskResponse, err error) {
	fmt.Printf("hello Create\n")
	return nil, errdefs.ErrNotImplemented
}

// Start the primary user process inside the container
func (s *micaTaskService) Start(ctx context.Context, r *taskAPI.StartRequest) (*taskAPI.StartResponse, error) {
	fmt.Printf("hello Start\n")
	return nil, errdefs.ErrNotImplemented
}

// Delete a process or container
func (s *micaTaskService) Delete(ctx context.Context, r *taskAPI.DeleteRequest) (*taskAPI.DeleteResponse, error) {
	fmt.Printf("hello Delete\n")
	return nil, errdefs.ErrNotImplemented
}

// Exec an additional process inside the container
func (s *micaTaskService) Exec(ctx context.Context, r *taskAPI.ExecProcessRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello Exec\n")
	return nil, errdefs.ErrNotImplemented
}

// ResizePty of a process
func (s *micaTaskService) ResizePty(ctx context.Context, r *taskAPI.ResizePtyRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello ResizePty\n")
	return nil, errdefs.ErrNotImplemented
}

// State returns runtime state of a process
func (s *micaTaskService) State(ctx context.Context, r *taskAPI.StateRequest) (*taskAPI.StateResponse, error) {
	fmt.Printf("hello State\n")
	return nil, errdefs.ErrNotImplemented
}

// Pause the container
func (s *micaTaskService) Pause(ctx context.Context, r *taskAPI.PauseRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello Pause\n")
	return nil, errdefs.ErrNotImplemented
}

// Resume the container
func (s *micaTaskService) Resume(ctx context.Context, r *taskAPI.ResumeRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello Resume\n")
	return nil, errdefs.ErrNotImplemented
}

// Kill a process
func (s *micaTaskService) Kill(ctx context.Context, r *taskAPI.KillRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello Kill\n")
	return nil, errdefs.ErrNotImplemented
}

// Pids returns all pids inside the container
func (s *micaTaskService) Pids(ctx context.Context, r *taskAPI.PidsRequest) (*taskAPI.PidsResponse, error) {
	fmt.Printf("hello Pids\n")
	return nil, errdefs.ErrNotImplemented
}

// CloseIO of a process
func (s *micaTaskService) CloseIO(ctx context.Context, r *taskAPI.CloseIORequest) (*ptypes.Empty, error) {
	fmt.Printf("hello CloseIO\n")
	return nil, errdefs.ErrNotImplemented
}

// Checkpoint the container
func (s *micaTaskService) Checkpoint(ctx context.Context, r *taskAPI.CheckpointTaskRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello Checkpoint\n")
	return nil, errdefs.ErrNotImplemented
}

// Connect returns shim information of the underlying service
func (s *micaTaskService) Connect(ctx context.Context, r *taskAPI.ConnectRequest) (*taskAPI.ConnectResponse, error) {
	fmt.Printf("hello Connect\n")
	return nil, errdefs.ErrNotImplemented
}

// Shutdown is called after the underlying resources of the shim are cleaned up and the service can be stopped
func (s *micaTaskService) Shutdown(ctx context.Context, r *taskAPI.ShutdownRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello Shutdown\n")
	os.Exit(0)
	return &ptypes.Empty{}, nil
}

// Stats returns container level system stats for a container and its processes
func (s *micaTaskService) Stats(ctx context.Context, r *taskAPI.StatsRequest) (*taskAPI.StatsResponse, error) {
	fmt.Printf("hello Stats\n")
	return nil, errdefs.ErrNotImplemented
}

// Update the live container
func (s *micaTaskService) Update(ctx context.Context, r *taskAPI.UpdateTaskRequest) (*ptypes.Empty, error) {
	fmt.Printf("hello Update\n")
	return nil, errdefs.ErrNotImplemented
}

// Wait for a process to exit
func (s *micaTaskService) Wait(ctx context.Context, r *taskAPI.WaitRequest) (*taskAPI.WaitResponse, error) {
	fmt.Printf("hello Wait\n")
	return nil, errdefs.ErrNotImplemented
} 