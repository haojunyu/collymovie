package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "配置信息",
		Long:  `查看设置配置信息`,
		Run: func(cmd *cobra.Command, args []string) {
			configure()
		},
	}
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(configCmd)
}

// config 数据库初始化
func configure() {
	log.Printf("config....\n")
}
