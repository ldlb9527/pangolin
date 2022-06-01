package model

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Client struct {
	Conn net.Conn
	// 数据传输通道
	ReadChan  chan []byte
	WriteChan chan []byte
	// 异常退出通道
	Exit chan error
	// 重连通道
	ReConn chan bool
}

// 从Client端读取数据
func (c *Client) Read(ctx context.Context) {
	// 如果10秒钟内没有消息传输，则Read函数会返回一个timeout的错误
	_ = c.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := make([]byte, 10240)
			n, err := c.Conn.Read(data)
			if err != nil && err != io.EOF {
				if strings.Contains(err.Error(), "timeout") {
					// 设置读取时间为3秒，3秒后若读取不到, 则err会抛出timeout,然后发送心跳
					_ = c.Conn.SetReadDeadline(time.Now().Add(time.Second * 3))
					c.Conn.Write([]byte("pi"))
					continue
				}
				fmt.Println("读取出现错误...")
				c.Exit <- err
				return
			}

			// 收到心跳包,则跳过
			if data[0] == 'p' && data[1] == 'i' {
				fmt.Println("server收到心跳包")
				continue
			}
			c.ReadChan <- data[:n]
		}
	}
}

// 将数据写入到Client端
func (c *Client) Write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-c.WriteChan:
			_, err := c.Conn.Write(data)
			if err != nil && err != io.EOF {
				c.Exit <- err
				return
			}
		}
	}
}
