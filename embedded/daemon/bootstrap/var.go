package bootstrap

import "runtime"

const (
	PATH_DRWN = "/dev/cu.usbmodem14101"
	PATH_LNX  = "/dev/ttyACM0"

	LOGS = "dev.log"
)

func SetDevice() string {
	if runtime.GOOS == "darwin" {
		return PATH_DRWN
	} else {
		return PATH_LNX
	}
}
