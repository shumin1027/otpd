package cmd

import (
	"os"

	"github.com/shumin1027/otpd/pkg/logger"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/shumin1027/otpd/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var conf = koanf.New(".")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: app.BuildInfo(),
	Use:     app.Name,
	Short:   "otp server",
	Long:    `otp server`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flags := cmd.PersistentFlags()
		BindPflags(conf, flags)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.L().Panic(err.Error())
		os.Exit(1)
	}
}

func BindPflags(conf *koanf.Koanf, flags *pflag.FlagSet) {
	provider := posflag.Provider(flags, ".", conf)
	if err := conf.Load(provider, nil); err != nil {
		logger.L().Fatal("error loading config: %v", zap.Error(err))
	}
}
