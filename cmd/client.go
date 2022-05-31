package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clientPort int

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "客户端启动命令",
	Long: `在内网环境快速启动客户端`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	rootCmd.Flags().IntVarP(&clientPort,"clientPort","cp",30000,"服务端端口")
}
