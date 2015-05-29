package main

import (
	"errors"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

type Transfer struct {
	RemoteAddr string
	RemotePort string
	LocalAddr  string
	LocalPort  string
	IsStop     bool
	TcpListen  *net.TCPListener
}

func NewTransfer(remoteAddr, remotePort, localAddr, localPort string) *Transfer {
	transfer := &Transfer{
		RemoteAddr: remoteAddr,
		RemotePort: remotePort,
		LocalAddr:  localAddr,
		LocalPort:  localPort,
		IsStop:     false,
	}
	return transfer
}

func (this *Transfer) Start() {
	go this.TcpTransfer()
}

func (this *Transfer) Stop() {
	this.IsStop = true
	this.TcpListen.Close()
}

// 拷贝数据：local<-->
func (this *Transfer) copyData(src, dst *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	buf := make([]byte, 32*1024)
	var ErrShortWrite = errors.New("short write")
	for {
		if this.IsStop {
			return
		}
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				log.Printf("%s --> %s. bytes:%d", src.RemoteAddr(), dst.RemoteAddr(), nw)
			}
			if ew != nil {
				log.Printf("ew: %s", ew)
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}

	if err != nil {
		if err.Error() != "EOF" {

			log.Printf("Error: %s", err.Error())
		} else {
			log.Printf("%s disconnected!!!", dst.RemoteAddr())
		}
	}
	dst.CloseWrite()
	src.CloseRead()
}

// 转换线程
func (this *Transfer) trans(local *net.TCPConn, remote *net.TCPConn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go this.copyData(remote, local, &wg)
	go this.copyData(local, remote, &wg)
	wg.Wait()
}

// 连接服务端
func (this *Transfer) connect(address string) (*net.TCPConn, error) {
	var remote *net.TCPConn

	remote_conn, err := net.DialTimeout("tcp", address, time.Second*5)

	if err != nil {
		log.Printf("Failed to connect to %s", address)
	} else {
		remote = remote_conn.(*net.TCPConn)
	}

	return remote, err
}

func (this *Transfer) TcpTransfer() {
	runtime.GOMAXPROCS(4)

	localAddr := this.LocalAddr + ":" + this.LocalPort
	address := this.RemoteAddr + ":" + this.RemotePort

	local, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Printf("Failed to open listening socket: %s", err)
	}
	this.TcpListen = local.(*net.TCPListener)
	log.Printf("Start listening on %s", localAddr)
	for {
		if this.IsStop {
			log.Printf("Stop listening on %s", localAddr)
			return
		}
		conn, err := local.Accept()
		if err != nil {
			continue
		}
		log.Printf("%s Connected!!", conn.RemoteAddr())
		if remote, _ := this.connect(address); remote != nil {
			go this.trans(conn.(*net.TCPConn), remote)
		} else {
			conn.Close()
			return
		}
	}
}
