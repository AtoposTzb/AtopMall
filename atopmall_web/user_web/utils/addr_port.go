package utils

import (
	"net"
)

// GetAddrPort 获取可用的端口号
func GetAddrPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0") //ResolveTCPAddr 解析TCP地址,0表示动态获取可用端口号
	if err != nil {
		return 0, err
	}
	conn, err := net.ListenTCP("tcp", addr) //ListenTCP 监听TCP地址
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return conn.Addr().(*net.TCPAddr).Port, nil //返回可用的端口号
}
