package model

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Server struct {
	Conn      net.Conn
	ReadChan  chan []byte
	WriteChan chan []byte
	Exit      chan error
	ReConn    chan bool
}

// 从Server端读取数据
func (s *Server) Read(ctx context.Context) {
	// 如果10秒钟内没有消息传输，则Read函数会返回一个timeout的错误
	_ = s.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := make([]byte, 10240)
			n, err := s.Conn.Read(data)
			if err != nil && err != io.EOF {
				// 读取超时，发送一个心跳包过去
				if strings.Contains(err.Error(), "timeout") {
					// 3秒发一次心跳
					_ = s.Conn.SetReadDeadline(time.Now().Add(time.Second * 3))
					s.Conn.Write([]byte("pi"))
					continue
				}
				fmt.Println("从server读取数据失败, ", err.Error())
				s.Exit <- err
				return
			}

			// 如果收到心跳包, 则跳过
			if data[0] == 'p' && data[1] == 'i' {
				fmt.Println("client收到心跳包")
				continue
			}
			s.ReadChan <- data[:n]
		}
	}
}

// 将数据写入到Server端
func (s *Server) Write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-s.WriteChan:
			_, err := s.Conn.Write(data)
			if err != nil && err != io.EOF {
				s.Exit <- err
				return
			}
		}
	}
}
