package main

import (
	"github.com/matishsiao/goInfo"
)

func main() {
	gi := goInfo.GetInfo()
	gi.VarDump()
}
