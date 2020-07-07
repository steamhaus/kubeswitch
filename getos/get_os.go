package getos

import (
	"github.com/matishsiao/goInfo"
)

type GoInfoObject struct {
	Kernel string
	OS     string
}

func GetOS() {
	gi := goInfo.GetInfo()
	gi.VarDump()
}
