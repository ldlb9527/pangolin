package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"pangolin/model"
)

var ClientPort int
var Ip string

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "客户端启动命令",
	Long:  `在内网环境快速启动客户端`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client start")
		ClientExec()
		fmt.Println("client end")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().IntVarP(&ClientPort, "clientPort", "c", 30000, "服务端端口")
	clientCmd.PersistentFlags().StringVarP(&Ip, "ip", "i", "127.0.0.1", "服务端地址")
}

func ClientExec() {
	target := net.JoinHostPort(Ip, fmt.Sprintf("%d", ClientPort))

	for {
		serverConn, err := net.Dial("tcp", target)
		if err != nil {
			panic(err)
		}

		fmt.Printf("已连接server: %s \n", serverConn.RemoteAddr())
		server := &model.Server{
			Conn:      serverConn,
			ReadChan:  make(chan []byte),
			WriteChan: make(chan []byte),
			Exit:      make(chan error),
			ReConn:    make(chan bool),
		}

		go handleServer(server)
		<-server.ReConn
		//_ = server.conn.Close()
	}
}

func handleServer(server *model.Server) {
	// 等待server端发来的信息，也就是说user来请求server了
	ctx, cancel := context.WithCancel(context.Background())

	go server.Read(ctx)
	go server.Write(ctx)

	localConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", LocalPort))
	if err != nil {
		panic(err)
	}

	local := &model.Local{
		Conn:      localConn,
		ReadChan:  make(chan []byte),
		WriteChan: make(chan []byte),
		Exit:      make(chan error),
	}

	go local.Read(ctx)
	go local.Write(ctx)

	defer func() {
		_ = server.Conn.Close()
		_ = local.Conn.Close()
		server.ReConn <- true
	}()

	for {
		select {
		case data := <-server.ReadChan:
			local.WriteChan <- data

		case data := <-local.ReadChan:
			server.WriteChan <- data

		case err := <-server.Exit:
			fmt.Printf("server have err: %s", err.Error())
			cancel()
			return
		case err := <-local.Exit:
			fmt.Printf("local have err: %s", err.Error())
			cancel()
			return
		}
	}
}
