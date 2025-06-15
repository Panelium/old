package main

import (
	"panelium/daemon/global"
)

func main() {
	err := global.Init()
	if err != nil {
		return
	}
}
