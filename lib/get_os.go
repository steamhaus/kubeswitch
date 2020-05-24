package main

import (
	"github.com/matishsiao/goInfo"
)

type GoInfoObject struct {
	GoOS string
	Kernel string
	Core string
	Platform string
	OS string
	Hostname string
	CPUs int
}

func main() {
	gi := goInfo.GetInfo()
	gi.VarDump()
	}
}
