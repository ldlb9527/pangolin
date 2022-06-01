package cmd

import (
	"context"
	"fmt"
	"net"
	"pangolin/model"

	"github.com/spf13/cobra"
)

var ServerPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务端启动命令",
	Long:  `快速在公网服务器启动服务端程序`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server start")
		ServerExec()
		fmt.Println("server end")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().IntVarP(&ServerPort, "serverPort", "s", 30000, "服务端端口")
}

func ServerExec() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Print("recover,")
			fmt.Println(err)
		}
	}()

	clientListener, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerPort))
	if err != nil {
		panic(err)
	}
	fmt.Printf("server端监听:%d端口, 提供给client使用 \n", ServerPort)

	userListener, err := net.Listen("tcp", fmt.Sprintf(":%d", LocalPort))
	if err != nil {
		panic(err)
	}
	fmt.Printf("server端监听:%d端口, 提供给用户使用 \n", LocalPort)

	for {
		// 有Client来连接了
		clientConn, err := clientListener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Printf("有Client连接: %s \n", clientConn.RemoteAddr())

		client := &model.Client{
			Conn:      clientConn,
			ReadChan:  make(chan []byte),
			WriteChan: make(chan []byte),
			Exit:      make(chan error),
			ReConn:    make(chan bool),
		}

		userConnChan := make(chan net.Conn)
		go AcceptUserConn(userListener, userConnChan)

		go HandleClient(client, userConnChan)

		<-client.ReConn
		fmt.Println("重新等待新的client连接..")
	}

}

func AcceptUserConn(userListener net.Listener, connChan chan net.Conn) {
	userConn, err := userListener.Accept()
	if err != nil {
		panic(err)
	}
	fmt.Printf("user connect: %s \n", userConn.RemoteAddr())
	connChan <- userConn
}

func HandleClient(client *model.Client, userConnChan chan net.Conn) {
	ctx, cancel := context.WithCancel(context.Background())

	go client.Read(ctx)
	go client.Write(ctx)

	user := &model.User{
		ReadChan:  make(chan []byte),
		WriteChan: make(chan []byte),
		Exit:      make(chan error),
	}

	defer func() {
		_ = client.Conn.Close()
		_ = user.Conn.Close()
		client.ReConn <- true
	}()

	for {
		select {
		case userConn := <-userConnChan:
			user.Conn = userConn
			go handle(ctx, client, user)
		case err := <-client.Exit:
			fmt.Println("client出现错误, 关闭连接", err.Error())
			cancel()
			return
		case err := <-user.Exit:
			fmt.Println("user出现错误，关闭连接", err.Error())
			cancel()
			return
		}
	}
}

// 将两个Socket通道链接
// 1. 将从user收到的信息发给client
// 2. 将从client收到信息发给user
func handle(ctx context.Context, client *model.Client, user *model.User) {

	go user.Read(ctx)
	go user.Write(ctx)

	for {
		select {
		case userRecv := <-user.ReadChan:
			// 收到从user发来的信息
			client.WriteChan <- userRecv
		case clientRecv := <-client.ReadChan:
			// 收到从client发来的信息
			user.WriteChan <- clientRecv

		case <-ctx.Done():
			return
		}
	}
}
