package app

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

const (
	Name        = "optd"
	Version     = "1.0.0"
	Description = "OTP Server"
	Repository  = "https://github.com/shumin1027/otpd"
)

var (
	GitTag    string
	GitCommit string
	GitBranch string
	BuildTime string
)

func BuildInfo() string {
	info := map[string]string{}
	info["Version"] = Version
	info["BuildTime"] = BuildTime
	info["GitCommit"] = GitCommit
	info["GitBranch"] = GitBranch
	info["GitTag"] = GitTag
	j, err := json.Marshal(info)
	if err != nil {
		log.Panic(err)
	}
	return string(j)
}

const templ = `
	{{ .AnsiColor.BrightGreen }}{{ .Title "%s" "standard" 0 }}{{ .AnsiColor.BrightCyan }}
	GOOS: {{ .GOOS }}
	GOARCH: {{ .GOARCH }}
	GoVersion: {{ .GoVersion }}
	Compiler: {{ .Compiler }}
	NumCPU: {{ .NumCPU }}
	Now: {{ .Now "Monday, 2 Jan 2006" }}
	{{ .AnsiColor.Default }}
	`

func PrintBanner() {
	banner.InitString(colorable.NewColorableStdout(), true, true, fmt.Sprintf(templ, Name))
	fmt.Print("\n")
}
