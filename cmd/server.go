package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务端启动命令",
	Long: `快速在公网服务器启动服务端程序`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.Flags().IntVarP(&serverPort,"serverPort","sp",30000,"服务端端口")
}
