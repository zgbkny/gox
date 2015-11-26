package udptunnel

/*****************************************
 * just do something about tunnel, for example
 * -- control rate
 *
 *****************************************/

import (
	"net"
	"log"
	"os"
	"sync"
//	"utils"
	"udpsession"
//	"udppacket"
)


/** 同一个目的地的集合 **/
type UDPTunnel struct {
	Reserved			int					// 数据区前预留的头部空间大小

	dst					string
	listenAddr			string
	send				chan []byte
	conn				*net.UDPConn
	addr				*net.UDPAddr
	idSessionMap		map[uint32]*udpsession.Session
	connIdMap			map[net.Conn]uint32
	onDataF				func([]byte) int
	loopRead			func(*net.Conn, uint32)


	// 统计
	sessionCount		uint32				// 当运行于客户端时用于产生session id，服务端只是用于统计
}
var ll *sync.Mutex
var count int
const MAX = 1000



/**
 * 启动客户端
 **/
func CreateClientTunnel(onDataF func([]byte) int) *UDPTunnel {
	ut := createUDPTunnel()
	ut.onDataF = onDataF
	log.Println("udptunnel Init")
	ut.dst = "192.168.80.128:9001"

	ut.initClientTunnel()	
	return ut
}
func CreateServerTunnel(onDataF func([]byte) int) *UDPTunnel {
	ut := createUDPTunnel()
	ut.onDataF = onDataF
	ut.listenAddr = ":9001"
	
	return ut
}


/***********************创建对象***********************/
func createUDPTunnel() *UDPTunnel {
	ut := new(UDPTunnel)

	ut.send = make(chan []byte)
	// 初始化环境
	ut.idSessionMap = map[uint32]*udpsession.Session{}
	ut.connIdMap = map[net.Conn]uint32{}
	ut.sessionCount = 0
	ut.Reserved = 96
	return ut
}

//======================================================

func (ut *UDPTunnel)StartServer() {

	ut.initServerTunnel()
}


func (ut *UDPTunnel)initClientTunnel() {

	// 初始化连接
	addr, err := net.ResolveUDPAddr("udp", ut.dst)
	if err != nil {
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		os.Exit(1)
	}
	ut.conn = conn
	ut.addr = addr

	go ut.tunnelWriteToServerProxy()
	go ut.tunnelReadFromServerProxy()
}

func (ut *UDPTunnel)initServerTunnel() {
	// 初始化监听
	addr, err := net.ResolveUDPAddr("udp", ut.listenAddr)
	if err != nil {
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		os.Exit(1)
	}
	ut.conn = conn
	go ut.tunnelWriteToClientProxy()
	ut.tunnelReadFromClientProxy()
}

func (ut *UDPTunnel)AddNewConnId(conn net.Conn, id uint32) {
	_, ok := ut.idSessionMap[id]
	if !ok {
		log.Println("AddNewConnId error")
		return
	}
	ut.connIdMap[conn] = id
}

/**
 *	处理关闭的连接，主要是删除对应的session
 * 	todo:通知对端
 **/
func (ut *UDPTunnel)ProcessCloseConn(conn net.Conn) {
	id, ok := ut.connIdMap[conn]
	if ok {
		delete(ut.connIdMap, conn)
		_, ok2 := ut.idSessionMap[id]
		if ok2 {
			delete(ut.idSessionMap, id)
		}
	}
}
/**
 * 处理客户端网关写原始数据	
 * rawData: 原始的数据，但是前面预留了96(ut.Reserved)字节的头部空间
 **/
func (ut *UDPTunnel)WritePacketToServerProxy(data []byte) int {
	log.Println("udptunnel WritePacketToServerProxy", string(data))
	ut.send <- data
	return 1
}

/**
 * 处理服务器网关写原始数据
 * rawData: 原始数据，但是前面预留了96(ut.Reserved)字节的头空间
 **/
func (ut *UDPTunnel)WritePacketToClientProxy(data []byte) {
	log.Println("udptunnel WritePacketToClientProxy")
	ut.send <- data
}

/**
 * 从客户端读取到数据包
 * data 包含会话信息的数据
 **/
func (ut *UDPTunnel)readPacketFromClientProxy(data []byte) {
	log.Println("udptunnel readPacketFromClientProxy")
	ut.onDataF(data)
}

/**
 * 从服务器读取数据包
 * 
 **/
func (ut *UDPTunnel)readPacketFromServerProxy(data []byte) {
	log.Println("udptunnel readPacketFromServerProxy")
	ut.onDataF(data)
}

//-------------------------------------------------------------------
func(ut *UDPTunnel)tunnelWriteToServerProxy() {
	for {
		log.Println("udptunnel tunnelWrite")
		data, ok := <-ut.send
		if !ok {
			break
		}
		log.Println("connWrite", string(data))
		ut.conn.Write(data)
	}
}
func (ut *UDPTunnel)tunnelWriteToClientProxy() {
	for {
		log.Println("udptunnel tunnelWriteToClientProxy")
		data, ok := <-ut.send
		if !ok {
			break
		}
		ut.conn.WriteToUDP(data, ut.addr)
	}
}

/**
 * 服务端往客户端写数据必须知道对方的地址
 **/
func(ut *UDPTunnel)tunnelReadFromClientProxy() {
	for {
		log.Println("udptunnel tunnelReadFromClientProxy")
		data := make([]byte, 4096)
		n, addr, err := ut.conn.ReadFromUDP(data)
		ut.addr = addr
		log.Println("after read", n)
		if err != nil {
			return
		}
		log.Println("data len", n)
		go ut.readPacketFromClientProxy(data[:n])
	}
}
/**
 * 客户端往服务端写，只需要有连接就可以了
 **/
func(ut *UDPTunnel)tunnelReadFromServerProxy() {
	for {
		log.Println("udptunnel tunnelReadFromServerProxy")
		data := make([]byte, 4096)
		n, _, err := ut.conn.ReadFromUDP(data)
		log.Println("after read", n)
		if err != nil {
			return
		}
		log.Println("data len", n)
		go ut.readPacketFromServerProxy(data[:n])
	}
}
