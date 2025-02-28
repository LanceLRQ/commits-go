package main

import (
	"github.com/LanceLRQ/commits-go/cmd"
	"github.com/LanceLRQ/commits-go/utils"
)

func main() {
	utils.InitI18n()
	cmd.CommandEntry()
}
