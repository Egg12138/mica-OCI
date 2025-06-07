package pseudo_container

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"rmica/communication"
	"rmica/defs"
	"rmica/logger"
	"rmica/mcs"
	"rmica/utils"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"

	securejoin "github.com/cyphar/filepath-securejoin"
)

// ==================== Type Definitions ====================

type Container struct {
	id      string
	// from global flag --root: 
	root    string
	// stateDir = root/id
	// stateDir string
	// Use specs-go::Spec, State to represent, following OCI-spec
	config  *specs.Spec
	state   ContainerState
	initPid int
	created time.Time
	m 			sync.Mutex
	// TODO: MCS client manager, will defined in mcs.go
	// clientManager *clientManager
}

// TODO: add more members
type runner struct {
	init          bool
	shouldDestory bool
	detach        bool
	pidFile       string
	container     *Container
	action        defs.CtAct
	notifySocket  *notifySocket
	consoleSocket string
	criuOpts      *libcontainer.CriuOpts
}


// ==================== Getters and Setters ====================

func (c *Container) Id() string {
	return c.id
}

func (c *Container) Root() string {
	return c.root
}

func (c *Container) StateDir() string {
	// return c.stateDir
	return filepath.Join(c.root, c.id)
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
	stats := NewEmpty()
	return &stats, nil
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
	stateDir := c.StateDir()
	if err := os.MkdirAll(stateDir, defs.ContainerDirPerm); err != nil {
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
	if err := c.createExecFifo(defs.ExecFifoFilename); err != nil {
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

func (c *Container) Run() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.run()
}

func (c *Container) run() error {
	logger.Infof("[container] run called for id=%s", c.id)
	res := communication.Send2mica("run")
	if res == "" {
		logger.Errorf("[container] run failed for id=%s: empty response from micad", c.id)
		return fmt.Errorf("failed to send run to micad or got empty response")
	}
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

func (c *Container) Restore() error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.restore()
}

func (c *Container) restore() error {
	logger.Infof("[container] restore called for id=%s", c.id)
	res := communication.Send2mica("restore")
	if res == "" {
		logger.Errorf("[container] restore failed for id=%s: empty response from micad", c.id)
		return fmt.Errorf("failed to send restore to micad or got empty response")
	}
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
	tmpFile, err := os.CreateTemp(c.StateDir(), "state-")
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

		stateFilePath := filepath.Join(c.StateDir(), defs.StateFilename)
		return os.Rename(tmpFile.Name(), stateFilePath)

}

func (c *Container) updateState(clientProcess *mcs.ClientTask) (*specs.State, error) {
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

// ==================== Notify Socket Operations ====================

// From runc:
// As a container monitor
type notifySocket struct {
	// ok, as a UDS listenr
	socket     *net.UnixConn
	// should not be Nil
	host       string
	// ok, mark a listner
	socketPath string
}

func newNotifySocket(context *cli.Context, notifySocketHost string, id string) *notifySocket {
	// Basically, notifier does not matter. Hence we just return when host is empty.
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

// If systemd is supporting sd_notify protocol, this function will add support
// for sd_notify protocol from within the container.
func (s *notifySocket) setupSpec(spec *specs.Spec) {
	pathInContainer := filepath.Join("/run/notify", path.Base(s.socketPath))
	mount := specs.Mount{
		Destination: path.Dir(pathInContainer),
		Source:      path.Dir(s.socketPath),
		Options:     []string{"bind", "nosuid", "noexec", "nodev", "ro"},
	}
	spec.Mounts = append(spec.Mounts, mount)
	spec.Process.Env = append(spec.Process.Env, "NOTIFY_SOCKET="+pathInContainer)
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
	// select {} 
	return nil
}

// ==================== Verification Utilities ====================

func verifyContainerDir(container *Container) {
	root := container.Root()
	id := container.Id()
	statePath := filepath.Join(container.StateDir(), defs.StateFilename)
	stateFile, err := os.OpenFile(statePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		logger.Errorf(
			"[rmica] verifyContainerDir: failed to open state file %s id=%s: %v", 
			statePath, id, err)
	}
	logger.Infof("[rmica] verifyContainerDir called for id=%s", id)
	logger.Infof("[rmica] root=%s, state.json=%s", root, statePath)
	defer stateFile.Close()
	// recursively travel root dir and print like tree:
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		logger.Infof("[rmica] %s", path)
		return nil
	})
}

// ==================== Runner Operations ====================
// TODO: For compatibility, use runc libcontainer.CriuOpts as criuOpts.
// LEARN: 后续参考 kata runtime 的思路
func StartContainer(context *cli.Context, action defs.CtAct, criuOpts *libcontainer.CriuOpts) (int, error) {
	if err := utils.RevisePidFile(context); err != nil {
		return -1, err
	}
	spec, err := utils.SetupSpec(context)
	if err != nil {
		return -2, fmt.Errorf("failed to load spec: %w", err)
	}

	cntrId := context.Args().First()
	if cntrId == "" {
		return -3, utils.ErrEmptyID
	}

	// IDEA: StartContainer 进程（runc/rmica本身)与子进程（容器1号进程）进行通信同步；
	// 而这个通信同步管理，混部的runtime部分需要承担多少? ——哪些events需要通知？
	// NOTICE: 整合后，runtime部分就属于 mica daemon, 所以是同一个进程；(****)
	// rmica 本身会 wait for "ready" 或其他, 
	notifySocket := newNotifySocket(context, os.Getenv("NOTIFY_SOCKET"), cntrId)
	if notifySocket != nil {
		// update ENV and Mount information to spec
		// IDEA: does running mica client needs NOTIFY_SOCKET=... ENV and Mount information?
		notifySocket.setupSpec(spec)
	}

	cntr, err := createContainer(context, cntrId, spec)
	if err != nil {
		return -4, fmt.Errorf("failed to create container: %w", err)
	}

	if notifySocket != nil {
		if err := notifySocket.setupSocketDirectory(); err != nil {
			return -5, fmt.Errorf("failed to setup socket directory: %w", err)
		}
		if action == defs.CT_ACT_RUN {
			if err := notifySocket.bindSocket(); err != nil {
				return -6, fmt.Errorf("failed to bind socket: %w", err)
			}
		}
	}

	ct := &mcs.ClientTask{
		Terminal: false,
		Name: "test",
		Tty: "/dev/micatty",
	}

	r := &runner{
		init:  true,
		shouldDestory: !context.Bool("keep"),
		detach:  context.Bool("detach"),
		pidFile: context.String("pid-file"),
		container: cntr,
		action: action,
		notifySocket: notifySocket,
		criuOpts: criuOpts,
	}

	if pidFile := context.String("pid-file"); pidFile != "" {
		if err := utils.CreatePidFile(pidFile, os.Getpid()); err != nil {
			return -7, fmt.Errorf("failed to create pid file: %w", err)
		}
	}

	return r.runTask(ct)
}

// IDEA: what do we need to run a task?
// TODO: update ClientProcess
// TODO: wrap many members in `runner`, as action, terminal, etc.
// LEARN: analyze why LISTEN_FDS is used in runc?
func (r *runner) runTask(taskConfig *mcs.ClientTask) (int, error) {
	var err error = nil
	defer func() {
		if err != nil {
			r.destroy()
		}
	}()

	if err = r.checkTerminal(taskConfig); err != nil {
		return -1, err
	}

	// task, err := newTaskFromConfig(taskConfig)
	logger.Infof("[rmica] runTask called for id=%s", r.container.Id())
	
	switch r.action {
	case defs.CT_ACT_RUN:
		err = r.container.Run()
	case defs.CT_ACT_CREATE:
		err = r.container.Start()
	case defs.CT_ACT_RESTORE:
		err = r.container.Restore()
	}
	
	verifyContainerDir(r.container)

	return 0, err
}

func (r *runner) destroy() {
	if r.shouldDestory {
		if err := r.container.Destroy(); err != nil {
			logger.Warnf("[rmica] failed to destroy container<%s>: %v", r.container.Id(), err)
		}
	}
}

// TODO: mica console 
func (r *runner) checkTerminal(taskConfig *mcs.ClientTask) error {
	detach := r.detach || (r.action == defs.CT_ACT_CREATE)
  if detach && taskConfig.Terminal && r.consoleSocket == "" {
		return errors.New("cannot allocate tty if rmica will detach without setting a console socket")
	}
	if (!detach || !taskConfig.Terminal) && r.consoleSocket != "" {
		return errors.New("cannot allocate tty if rmica will not detach")
	}
	if taskConfig.Terminal && r.consoleSocket != "" {
		return errors.New("cannot allocate tty if rmica will not detach or allocate tty")
	}
	return nil
}

// NOTICE: in runc, newTaskFromConfig() converts spec to 
// but in rmica, it is just a dummy
func newTaskFromConfig(taskConfig *mcs.ClientTask) (*mcs.ClientTask, error) {
	return taskConfig, nil
}

// ==================== Container Utilities ====================

func createContainer(context *cli.Context, id string, spec *specs.Spec) (*Container, error) {
	return Create(context.GlobalString("root"), id, spec)
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

// NOTICE We create state dir in host for container engine
func Create(root, id string, config *specs.Spec) (*Container, error) {
	if root == "" {
		return nil, errors.New("root is empty")
	}

	if err := utils.ValidateID(id); err != nil {
		return nil, err
	}

	if err := utils.ValidateSpec(config); err != nil {
		return nil, err
	}

	if err := utils.ValidateTaskSpec(config); err != nil {
		return nil, err
	}


	if err := os.MkdirAll(root, 0o711); err != nil {
		// return nil, utils.DebugPrintf("failed to create root dir: %s; %w", root, err)
		return nil, fmt.Errorf("failed to create root dir: %s; %w", root, err)
	}

	stateDir, err := securejoin.SecureJoin(root, id)
	if err != nil {
		// return nil, utils.DebugPrintf("failed to create state directory: %s; %w", stateDir, err)
		return nil, fmt.Errorf("failed to create state directory: %s; %w", stateDir, err)
	}

	if _, err := os.Stat(stateDir); err == nil {
		return nil, utils.ErrExist
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to stat state directory: %s; %w", stateDir, err)
	}

	if err := os.Mkdir(stateDir, 0o711); err != nil {
		return nil, fmt.Errorf("failed to create state directory for parent: %s; %w", stateDir, err)
	}

	// TODO: create network namespace 

	// TODO: Members initialization.
	cntr := &Container{
		id: id,
		root: root,
		config: config,
	}
	cntr.state = &StoppedState{c: cntr}
	return cntr, nil
}
