package main

import (
	"github.com/shumin1027/otpd/app"
	"github.com/shumin1027/otpd/cmd"
)

func init() {
	app.PrintBanner()
}

func main() {
	cmd.Execute()
}
