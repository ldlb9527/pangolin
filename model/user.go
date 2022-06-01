package model

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type User struct {
	Conn net.Conn
	// 数据传输通道
	ReadChan  chan []byte
	WriteChan chan []byte
	// 异常退出通道
	Exit chan error
}

// 从User端读取数据
func (u *User) Read(ctx context.Context) {
	_ = u.Conn.SetReadDeadline(time.Now().Add(time.Second * 200))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("从User端读取数据,执行了return")
			return
		default:
			data := make([]byte, 10240)
			n, err := u.Conn.Read(data)
			if err != nil && err != io.EOF {
				fmt.Println("err != nil && err != io.EOF 为真")
				u.Exit <- err
				return
			}
			fmt.Println("从user端读取的数据+" + string(data))
			u.ReadChan <- data[:n]
		}
	}
}

// 将数据写给User端
func (u *User) Write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-u.WriteChan:
			_, err := u.Conn.Write(data)
			if err != nil && err != io.EOF {
				u.Exit <- err
				return
			}
		}
	}
}
