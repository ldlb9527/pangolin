package model

import (
	"context"
	"net"
)

type Local struct {
	Conn net.Conn
	// 数据传输通道
	ReadChan  chan []byte
	WriteChan chan []byte
	// 有异常退出通道
	Exit chan error
}

func (l *Local) Read(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := make([]byte, 10240)
			n, err := l.Conn.Read(data)
			if err != nil {
				l.Exit <- err
				return
			}
			l.ReadChan <- data[:n]
		}
	}
}

func (l *Local) Write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-l.WriteChan:
			_, err := l.Conn.Write(data)
			if err != nil {
				l.Exit <- err
				return
			}
		}
	}
}
