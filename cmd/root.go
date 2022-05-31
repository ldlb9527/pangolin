/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var localPort int

var rootCmd = &cobra.Command{
	Use:   "pangolin",
	Short: "内网穿透命令行工具",
	Long: `使用该工具可快速启动服务端和客户端，达到内网穿透效果`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().IntVarP(&localPort,"localPort","lp",8080,"本地端口，服务端指用户访问端口，客户端指程序端口")
}


