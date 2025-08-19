package posix

import (
	"strings"

	"github.com/cuhsat/futils/pkg/sys"
)

func ModList(mod string) (ls []string, err error) {
	out, err := sys.StdCall("lsmod")

	if err != nil {
		return
	}

	ls = strings.Split(strings.TrimSpace(out), "\n")

	return
}

func ModLoad(args ...string) (err error) {
	_, err = sys.StdCall("modprobe", args...)

	return
}
