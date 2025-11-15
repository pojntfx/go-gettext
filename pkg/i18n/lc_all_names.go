package i18n

import (
	"errors"
	"runtime"
)

func getLCALL() (int, error) {
	switch runtime.GOOS {
	case "darwin":
		return 0, nil // BSD, but macOS is the same: https://github.com/openbsd/src/blob/master/include/locale.h#L67

	case "linux":
		return 6, nil // See https://github.com/rofl0r/musl/blob/master/include/locale.h#L24

	case "windows":
		return 0, nil // See https://github.com/tpn/winsdk-10/blob/master/Include/10.0.16299.0/ucrt/locale.h#L18

	default:
		return -1, errors.New("this GOOS is not supported")
	}
}
