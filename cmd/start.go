package cmd

import (
	"fmt"

	"github.com/shumin1027/otpd/http"
	"github.com/shumin1027/otpd/pkg/logger"
	"github.com/shumin1027/otpd/pkg/otp"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start otp server",
	Long:  `Start otp server`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.SetGlobal(logger.Config{
			Filenames:  conf.Strings("log.path"),
			MaxSize:    conf.Int("log.maxsize"),
			MaxAge:     conf.Int("log.maxage"),
			MaxBackups: conf.Int("log.maxbackups"),
			LocalTime:  conf.Bool("log.localtime"),
			Compress:   conf.Bool("log.compress"),
			LogLevel:   conf.String("log.level"),
			Encoder:    conf.String("log.format"),
		})

		otp.Init(conf.String("data.path"))

	},
	Run: func(cmd *cobra.Command, args []string) {
		bind := conf.String("bind")
		port := conf.Int("port")
		addr := fmt.Sprintf("%s:%d", bind, port)
		http.Start(addr)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	flags := startCmd.PersistentFlags()
	flags.IntP("port", "p", 18181, "web listening port")
	flags.StringP("bind", "b", "0.0.0.0", "bind ip addr")
	flags.StringP("data.path", "d", "", "data path")
	flags.StringSliceP("log.path", "", []string{"stderr"}, "log path, support stdout, stderr and file")
	flags.IntP("log.maxsize", "", 100, "log file size megabytes")
	flags.IntP("log.maxage", "", 90, "log file retain days")
	flags.IntP("log.maxbackups", "", 5, "log file retain nums")
	flags.BoolP("log.localtime", "", false, "log file using localtime")
	flags.BoolP("log.compress", "", false, "log file rotate compress")
	flags.StringP("log.level", "", "info", "log level, support debug, info, warn, error, dpanic, panic, fatal")
	flags.StringP("log.format", "", "console", "log format, support json and consolel")
}
