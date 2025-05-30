package mica

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

const (
	MicaConfigPath = "/tmp/mica"
	SocketPath     = "/tmp/mica"
	SocketName 		 = "mica-create.socket"
)

type CreateMsg struct {
	CPU     uint32
	Name    [32]byte
	Path    [128]byte
	Ped     [32]byte
	PedCfg  [128]byte
	Debug   bool
}

// Pack serializes the CreateMsg into a byte slice
func (m *CreateMsg) Pack() []byte {
	buf := make([]byte, 0, 325) // 4 + 32 + 128 + 32 + 128 + 1 Bytes

	// Pack CPU (uint32)
	cpuBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(cpuBuf, m.CPU)
	buf = append(buf, cpuBuf...)

	// Pack Name (32 bytes)
	buf = append(buf, m.Name[:]...)

	// Pack Path (128 bytes)
	buf = append(buf, m.Path[:]...)

	// Pack Ped (32 bytes)
	buf = append(buf, m.Ped[:]...)

	// Pack PedCfg (128 bytes)
	buf = append(buf, m.PedCfg[:]...)

	// Pack Debug (bool)
	if m.Debug {
		buf = append(buf, 1)
	} else {
		buf = append(buf, 0)
	}

	return buf
}

type Socket struct {
	conn  net.Conn
	debug bool
}

func NewSocket(socketPath string, debug bool) (*Socket, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", socketPath, err)
	}
	return &Socket{conn: conn, debug: debug}, nil
}

func (s *Socket) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *Socket) SendMsg(msg []byte) error {
	if s.conn == nil {
		return fmt.Errorf("socket not connected")
	}
	if s.debug {
		fmt.Printf("Sending data to %s: %v\n", s.conn.RemoteAddr(), msg)
	}
	_, err := s.conn.Write(msg)
	return err
}

func (s *Socket) Recv(bufferSize int, timeout time.Duration) (string, error) {
	if s.conn == nil {
		return "", fmt.Errorf("socket not connected")
	}

	s.conn.SetReadDeadline(time.Now().Add(timeout))
	buf := make([]byte, bufferSize)
	var response strings.Builder

	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return "", fmt.Errorf("timeout while waiting for micad response")
			}
			return "", err
		}

		if n == 0 {
			break
		}

		response.Write(buf[:n])
		respStr := response.String()

		if s.debug {
			fmt.Printf("Received response: %s\n", respStr)
		}

		if strings.Contains(respStr, "MICA-FAILED") {
			parts := strings.Split(respStr, "MICA-FAILED")
			msg := strings.TrimSpace(parts[0])
			if msg != "" {
				fmt.Println(msg)
			}
			fmt.Println("Error occurred!")
			fmt.Println("Please see system log ('cat /var/log/messages' or 'journalctl -u micad') for details.")
			return "MICA-FAILED", nil
		}

		if strings.Contains(respStr, "MICA-SUCCESS") {
			parts := strings.Split(respStr, "MICA-SUCCESS")
			msg := strings.TrimSpace(parts[0])
			if msg != "" {
				fmt.Println(msg)
			}
			return "MICA-SUCCESS", nil
		}
	}

	return response.String(), nil
}

func ParseConfig(configFile string) (*CreateMsg, error) {
	fmt.Println("Parsing config file:", configFile)
	cfg, err := ini.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v", err)
	}

	section := cfg.Section("Mica")
	if section == nil {
		return nil, fmt.Errorf("section 'Mica' not found in config file")
	}

	msg := &CreateMsg{
		Debug: true, 
	}

	if section.HasKey("CPU") {
		cpu, err := section.Key("CPU").Uint()
		if err != nil {
			return nil, fmt.Errorf("invalid CPU value: %v", err)
		}
		msg.CPU = uint32(cpu)
	}

	if section.HasKey("Name") {
		name := section.Key("Name").String()
		copy(msg.Name[:], name)
	}

	if section.HasKey("ClientPath") {
		path := section.Key("ClientPath").String()
		copy(msg.Path[:], path)
	}

	if section.HasKey("Pedestal") {
		ped := section.Key("Pedestal").String()
		copy(msg.Ped[:], ped)
	}
	if section.HasKey("PedestalConf") {
		pedCfg := section.Key("PedestalConf").String()
		copy(msg.PedCfg[:], pedCfg)
	}

	if section.HasKey("Debug") {
		msg.Debug = section.Key("Debug").MustBool(true)
	}

	return msg, nil
}

// TODO: 重复逻辑 configFile应该考虑缺省的embedded content
func SendCreateMsg(configFile string) error {
	micaConfig := configFile
	if !fileExists(micaConfig) {
		micaConfig = filepath.Join(MicaConfigPath, configFile)
		if !fileExists(micaConfig) {
			return fmt.Errorf("configuration file '%s' not found", configFile)
		}
	}

	target := filepath.Join(SocketPath, SocketName)
	if !fileExists(target) {
		return fmt.Errorf("error occurred! Please check if %s is running", target)
	}

	msg, err := ParseConfig(micaConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	fmt.Printf("Creating %s...\n", strings.TrimRight(string(msg.Name[:]), "\x00"))

	socket, err := NewSocket(filepath.Join(SocketPath, SocketName), msg.Debug)
	if err != nil {
		return err
	}
	defer socket.Close()

	if err := socket.SendMsg(msg.Pack()); err != nil {
		return err
	}

	response, err := socket.Recv(512, 5*time.Second)
	if err != nil {
		return err
	}

	if response == "MICA-SUCCESS" {
		fmt.Printf("Successfully created %s!\n", strings.TrimRight(string(msg.Name[:]), "\x00"))
	} else if response == "MICA-FAILED" {
		fmt.Printf("Create %s failed!\n", strings.TrimRight(string(msg.Name[:]), "\x00"))
	}

	return nil
}

func SendCtrlMsg(command, client string) error {
	ctrlSocket := filepath.Join(SocketPath, client+".socket")
	if !fileExists(ctrlSocket) {
		return fmt.Errorf("cannot find %s. Please run 'mica create <config>' to create it", client)
	}

	socket, err := NewSocket(ctrlSocket, true) // Enable debug by default for control messages
	if err != nil {
		return err
	}
	defer socket.Close()

	if err := socket.SendMsg([]byte(command)); err != nil {
		return err
	}

	response, err := socket.Recv(512, 5*time.Second)
	if err != nil {
		return err
	}

	if response == "MICA-SUCCESS" {
		fmt.Printf("%s %s successfully!\n", command, client)
	} else if response == "MICA-FAILED" {
		fmt.Printf("%s %s failed!\n", command, client)
	}

	return nil
}

func QueryStatus() error {
	if !fileExists(filepath.Join(SocketPath, "mica-create.socket")) {
		return fmt.Errorf("error occurred! Please check if micad is running")
	}

	fmt.Printf("%-30s%-20s%-20s%s\n", "Name", "Assigned CPU", "State", "Service")

	files, err := os.ReadDir(SocketPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == "mica-create.socket" || !strings.HasSuffix(file.Name(), ".socket") {
			continue
		}

		socket, err := NewSocket(filepath.Join(SocketPath, file.Name()), true)
		if err != nil {
			continue
		}

		if err := socket.SendMsg([]byte("status")); err != nil {
			socket.Close()
			continue
		}

		response, err := socket.Recv(512, 5*time.Second)
		socket.Close()

		if err != nil || response == "MICA-FAILED" {
			name := strings.TrimSuffix(file.Name(), ".socket")
			fmt.Printf("Query %s status failed!\n", name)
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
} 