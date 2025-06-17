package main

import (
	"panelium/daemon/internal/global"
)

func main() {
	err := global.Init()
	if err != nil {
		return
	}
}
