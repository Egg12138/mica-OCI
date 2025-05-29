package communication

// communication with MICAD
import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/Egg12138/mica-OCI/rmica/constants"
	"github.com/Egg12138/mica-OCI/rmica/logger"
)

func Send2mica(data string) string {
	res, err := send2socket(data, constants.DefaultMicaSocket)
	if err != nil {
		return ""
	}
	// fmt.Printf("valid response: %s\n", res)
	logger.Infof("valid response: %s", res)
	return res
}

func send2socket(data, path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to stat socket file: %w", err)
	}

	if fileInfo.Mode()&os.ModeSocket == 0 {
		return "", fmt.Errorf("%s is not a socket file", path)
	}

	conn, err := net.Dial("unix", path)
	if err != nil {
		return "", fmt.Errorf("failed to connect to socket: %w", err)
	}
	defer conn.Close()

	// 3. send data
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

