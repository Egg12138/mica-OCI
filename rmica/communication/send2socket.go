package communication

// communication with MICAD
import (
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"rmica/defs"
	"rmica/logger"
)

func Send2micaCmd(cmd string, args ...string) string {
	return Send2mica(strings.Join(append([]string{cmd}, args...), " "))
}

func Send2mica(data string) string {
	logger.Debugf("sending data [%s] to [%s]", data, defs.DefaultMicaSocket)
	res, err := send2socket(data, defs.DefaultMicaSocket)
	if err != nil {
		logger.Errorf("failed to send data to micad %v", err)
		logger.Fprintf("failed to send data to micad %v", err)
		return ""
	}
	logger.Infof("received response from micad: %s", res)
	return res
}

func send2socket(data, path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		logger.Fprintf("failed to stat socket file: %v", err)
		logger.Debugf("failed to stat socket file: %v", err)
		return "", fmt.Errorf("failed to stat socket file: %w", err)
	}

	if fileInfo.Mode()&os.ModeSocket == 0 {
		logger.Fprintf("%s is not a socket file", path)
		logger.Debugf("%s is not a socket file", path)
		return "", fmt.Errorf("%s is not a socket file", path)
	}

	conn, err := net.Dial("unix", path)
	if err != nil {
		logger.Fprintf("failed to connect to socket: %v", err)
		logger.Debugf("failed to connect to socket: %v", err)
		return "", fmt.Errorf("failed to connect to socket: %w", err)
	}
	defer conn.Close()

	logger.Debugf("trying to write byte sequence: %s \n [%x]", data, []byte(data))
	_, err = conn.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("failed to write to socket: %w", err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err != syscall.EAGAIN {
			return "", fmt.Errorf("failed to read response: %w", err)
		}
		// EAGAIN means no data available, which is fine
	}

	response := string(buf[:n])

	return response, nil
}

