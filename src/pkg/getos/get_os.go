package kubeswitch

import (
	"github.com/matishsiao/goInfo"
)

type GoInfoObject struct {
	Kernel string
	OS     string
}

func getos() {
	gi := goInfo.GetInfo()
	gi.VarDump()
}
