package pseudo_container

import (
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"rmica/communication"
	"rmica/constants"
	"rmica/logger"
	"rmica/mcs"
	"rmica/utils"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

// ==================== Type Definitions ====================


type Container struct {
	id      string
	root    string
	stateDir string
	// Use specs-go::Spec, State to represent, following OCI-spec
	config  *specs.Spec
	state   ContainerState
	initPid int
	created time.Time
	m 			sync.Mutex
	// TODO: MCS client manager, will defined in mcs.go
	// clientManager *clientManager
}


// ==================== Getters and Setters ====================


func (c *Container) Id() string {
	return c.id
}

func (c *Container) Root() string {
	return c.root
}

func (c *Container) StateDir() string {
	return c.stateDir
}


// TODO: handle cocurrency
// Status => specs::ContainerState, container status representation
// State => specs::State, runtime state of the container
func (c *Container) Status() specs.ContainerState {
	c.m.Lock()
	defer c.m.Unlock()
	return c.state.status()
	// return specs.StateCreating
}

// TODO: handle State(), OCIState() properly
func (c *Container) State() specs.State {
	c.m.Lock()
	defer c.m.Unlock()
	state := specs.State{
		Version:     specs.Version,
		ID:          c.Id(),
		Status:      c.state.status(),
		Pid:         c.initPid,
		Bundle:      c.Root(),
		Annotations: c.config.Annotations,
	}
	return state
}

// OCIState returns the current container's state information in OCI format
// This is duplicate with State() in rmica's implementation as we don't need
// the additional information that rmica's internal State type provides.
// just for compatibility with runc
func (c *Container) OCIState() *specs.State {
	c.m.Lock()
	defer c.m.Unlock()
	state := c.State()
	return &state
}

// ignoreCgroupError filters out cgroup-related errors that can be ignored,
// because the container is stopped and its cgroup is gone.
// TODO: in rmica, we do not need to handle cgroup error, **temporarily**.
func (c *Container) ignoreCgroupError(err error) error {
	return nil
}

// TODO: ClientProcesses returns the PIDs of the handler of client RTOS ()
// which marker should we focus on?
func (c *Container) ClientProcesses() ([]int, error) {
	return []int{c.initPid}, nil
}

// FIXME: the statistics function is 0% implemented
// Stats returns statistics for the container.
func (c *Container) Stats() (*Stats, error) {
	// c.m.Lock()
	// defer c.m.Unlock()
	stats := newEmpty()
	return &stats, nil
}



func Load(root, id string) (*Container, error) {
	containerDir := filepath.Join(root, id)
	if _, err := os.Stat(containerDir); err != nil {
		return nil, fmt.Errorf("container %s not found: %w", id, err)
	}

	return &Container{
		id:   id,
		root: root,
		state: &StoppedState{},
	}, nil
}

func (c *Container) Set(config *specs.Spec) error {
	c.m.Lock()
	defer c.m.Unlock()

	status := c.Status()
	if status == specs.StateStopped {
		return utils.ErrNotRunning
	}
	// TODO: tree-like recursive assignment for all sub items in config
	// (Set c.ClientManager 
	//   (Set c.ClientManager.Item0 (...)))

	c.config = config
	_, err := c.updateState(nil)	
	return err
}

// ==================== Container Lifecycle Operations ====================

func (c *Container) Init() error {
	// Create container directory 
	// e.g.: /run/containerd/io.containerd.runtime.task.v2/moby/<id>
	containerDir := filepath.Join(c.Root(), c.Id())
	if err := os.MkdirAll(containerDir, 0o700); err != nil {
		return fmt.Errorf("failed to create container directory: %w", err)
	}

	// c.State = &specs.State{
	// 	Version:     c.Spec.Version,
	// 	ID:          c.Id(),
	// 	Status:      specs.StateCreating,
	// 	Bundle:      c.Root,
	// 	Annotations: c.Spec.Annotations,
	// }

	return nil
}

func (c *Container) Start() error {
	c.m.Lock()
	defer c.m.Unlock()
	// return c.state.transition(&RunningState{c: c})
	return c.start()
}

func (c *Container) fakeStart() error {
	c.m.Lock()
	defer c.m.Unlock()
	logger.Infof("[pseudo-container] start called for id=%s", c.id)
  return nil
}

func (c *Container) start() error {
	c.m.Lock()
	defer c.m.Unlock()

	// 合法性检查：如 exec fifo 创建
	if err := c.createExecFifo(constants.ExecFifoFilename); err != nil {
		return fmt.Errorf("failed to create exec fifo: %w", err)
	}

	// 可以添加更多合法性检查，如状态检查等
	// if c.Status() != specs.StateCreated { ... }

	// 序列化请求并转发到 micad
	res := communication.Send2mica("start")
	if res == "" {
		return fmt.Errorf("failed to send start to micad or got empty response")
	}

	return nil
}

func (c *Container) Exec() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.exec()
}

func (c *Container) exec() error {
	logger.Infof("[container] exec called for id=%s", c.id)
	res := communication.Send2mica("exec")
	if res == "" {
		logger.Errorf("[container] exec failed for id=%s: empty response from micad", c.id)
		return fmt.Errorf("failed to send exec to micad or got empty response")
	}
	logger.Infof("[container] exec succeeded for id=%s, response=%s", c.id, res)
	return nil
}

func (c *Container) Stop() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.stop()
}

func (c *Container) stop() error {
	logger.Infof("[container] stop called for id=%s", c.id)
	res := communication.Send2mica("stop")
	if res == "" {
		logger.Errorf("[container] stop failed for id=%s: empty response from micad", c.id)
		return fmt.Errorf("failed to send stop to micad or got empty response")
	}
	logger.Infof("[container] stop succeeded for id=%s, response=%s", c.id, res)
	return nil
}

func (c *Container) Pause() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.pause()
}

func (c *Container) pause() error {
	logger.Infof("[container] pause called for id=%s", c.id)
	res := communication.Send2mica("pause")
	if res == "" {
		logger.Errorf("[container] pause failed for id=%s: empty response from micad", c.id)
		return fmt.Errorf("failed to send pause to micad or got empty response")
	}
	logger.Infof("[container] pause succeeded for id=%s, response=%s", c.id, res)
	return nil
}

func (c *Container) Resume() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.resume()
}

func (c *Container) resume() error {
	logger.Infof("[container] resume called for id=%s", c.id)
	res := communication.Send2mica("resume")
	if res == "" {
		logger.Errorf("[container] resume failed for id=%s: empty response from micad", c.id)
		return fmt.Errorf("failed to send resume to micad or got empty response")
	}
	logger.Infof("[container] resume succeeded for id=%s, response=%s", c.id, res)
	return nil
}

func (c *Container) Destroy() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.destroy()
}

func (c *Container) destroy() error {
	logger.Infof("[container] destroy called for id=%s", c.id)
	res := communication.Send2mica("destroy")
	if res == "" {
		logger.Errorf("[container] destroy failed for id=%s: empty response from micad", c.id)
		return fmt.Errorf("failed to send destroy to micad or got empty response")
	}
	logger.Infof("[container] destroy succeeded for id=%s, response=%s", c.id, res)
	return nil
}




// ==================== Helper Functions ====================

// HostRootUID returns the root uid for the process on host (always 0 for rmica, no user namespace)
func (c *Container) HostRootUID() (int, error) {
	return 0, nil
}

// HostRootGID returns the root gid for the process on host (always 0 for rmica, no user namespace)
func (c *Container) HostRootGID() (int, error) {
	return 0, nil
}

// createExecFifo creates a FIFO file for container exec, owned by root
func (c *Container) createExecFifo(fifoName string) error {
	fifoPath := filepath.Join(c.StateDir(), fifoName)
	if err := os.RemoveAll(fifoPath); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(fifoPath), 0o755); err != nil {
		return err
	}
	if err := syscall.Mkfifo(fifoPath, 0o622); err != nil {
		return err
	}
	// Assign with runc libcontainer style
	uid, _ := c.HostRootUID()
	gid, _ := c.HostRootGID()
	if err := os.Chown(fifoPath, uid, gid); err != nil {
		return err
	}
	if err := os.Chmod(fifoPath, 0o622); err != nil {
		return err
	}
	return nil
}


func (c *Container) hasInit() bool {
	return c.initPid != 0
} 

func (c *Container) saveState(s *specs.State) (retErr error) {
	tmpFile, err := os.CreateTemp(c.stateDir, "state-")
		if err != nil {
			return err
		}

		defer func() {
			if retErr != nil {
				tmpFile.Close()
				os.Remove(tmpFile.Name())
			}
		}()

		err = utils.WriteJSON(tmpFile, s)
		if err != nil {
			return err
		}
		err = tmpFile.Close()
		if err != nil {
			return err
		}

		stateFilePath := filepath.Join(c.stateDir, constants.StateFilename)
		return os.Rename(tmpFile.Name(), stateFilePath)

}

func (c *Container) updateState(clientProcess *mcs.ClientProcess) (*specs.State, error) {
	state := c.State()
	if err := c.saveState(&state); err != nil {
		return nil, err
	}
	return &state, nil
}

// Alway consider the instance as a `container`
func GetContainer(context *cli.Context) (*Container, error) {
	id := context.Args().First()
	if id == "" {
		return nil, utils.ErrEmptyID
	}
	root := context.GlobalString("root")
	return Load(root, id)
}


// From runc:
// As a container monitor
type notifySocket struct {
	socket     *net.UnixConn
	host       string
	socketPath string
}

func newNotifySocket(context *cli.Context, notifySocketHost string, id string) *notifySocket {

	if notifySocketHost == "" {
		return nil
	}

	cntrRoot := filepath.Join(utils.GetRootDir(context), id)
	socketPath := filepath.Join(cntrRoot, "notify", "notify.sock")

	notifySocket := &notifySocket{
		socket:     nil,
		host:       notifySocketHost,
		socketPath: socketPath,
	}

	return notifySocket
}

func (s *notifySocket) Close() error {
	return s.socket.Close()
}


func (s *notifySocket) bindSocket() error {
	addr := net.UnixAddr{
		Name: s.socketPath,
		Net:  "unixgram",
	}

	socket, err := net.ListenUnixgram("unixgram", &addr)
	if err != nil {
		return err
	}

	err = os.Chmod(s.socketPath, 0o777)
	if err != nil {
		socket.Close()
		return err
	}

	s.socket = socket
	return nil
}

func (s *notifySocket) setupSocketDirectory() error {
	return os.Mkdir(path.Dir(s.socketPath), 0o755)
}

func NotifySocketStart(context *cli.Context, notifySocketHost, id string) (*notifySocket, error) {
	notifySocket := newNotifySocket(context, notifySocketHost, id)
	if notifySocket == nil {
		return nil, nil
	}

	if err := notifySocket.bindSocket(); err != nil {
		return nil, err
	}
	return notifySocket, nil
}


func (s *notifySocket) WaitForContainer(container *Container) error {
	state := container.State()
	return s.fakeSuccessRun(state)

}

func (s *notifySocket) fakeSuccessRun(state specs.State) error {
	logger.Info("fakeSuccessRun", "state", state)
	notifySocketHostAddr := net.UnixAddr{Name: s.host, Net: "unixgram"}
	client, err := net.DialUnix("unixgram", nil, &notifySocketHostAddr)
	if err != nil {
		return err
	}
	defer client.Close()

	ready := []byte("READY=1\n")
	if _, err := client.Write(ready); err != nil {
		return err
	}

	mainPid := "MAINPID=12345\n"
	if _, err := client.Write([]byte(mainPid)); err != nil {
		return err
	}

	// 4. 可选：实现一个简单的 barrier
	// 创建一个管道用于 barrier 同步
	pipeR, pipeW, _ := os.Pipe()
	defer pipeR.Close()
	defer pipeW.Close()

	_, err = client.Write([]byte("BARRIER=1\n"))
	if err != nil {
		return err
	}
	time.Sleep(time.Hour)
	select {} 
	return nil
}


