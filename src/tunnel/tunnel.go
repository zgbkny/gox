package tunnel
import "net"
import "log"
import "container/list"
import "sync"
import "io"
//import "os"
/* tunnel 配置 */
type TunnelConf struct {
	Dst string
	
	InitConns int
	MaxConns int
}

// session
type tunnel1 struct {
	id int
	* list.Element
	send chan []byte
	remote io.Writer
}

// 对应一条tcp连接的所有相关信息 
type bundle struct {
	t [maxBundleClientConns] tunnel
	* list.List
	* xsocket
	sync.Mutex
}

/* 监听端口 */
type Listener struct {
	Bundles [maxLongConns]* bundle
	Id [maxLongConns] int
}

/* 隧道监听端口 */
type TunnelListener struct {
	
}

type session struct {
	id uint32
	* list.Element		// 节点指针 如果空闲就指向自己的list节点；如果在使用，nil
	c net.Conn
}

type tunnel struct {
	state int			// tunnel的状态
	send chan []byte	// 需要往tunnel发送的数据缓冲区
	* xsocket			// tunnel 连接
}

// 封装的conn
type xsocket struct {
	net.Conn
	* sync.Mutex	// 多条客户端连接写数据时需要
}

type tunnelManager struct {
	ss [maxClientConns] session
	* list.List

	connMp map[net.Conn]*session	
	idMp map[uint32]*session

	ts [maxTunnelConns] tunnel
	tunnelCount int
	sessionIdIndex uint32			// 循环递增id, id > 0
	* sync.Mutex					// 多条客户端连接同时使用
}


const maxBundleClientConns = 0x10
const maxLongConns = 10
const serverAddr = "localhost:9001"

const T_OPEN = 0
const T_CLOSE = 1
const bufferSize = 4096
const maxClientConns = 0x1000
const maxTunnelConns = 0x100

var CONF *TunnelConf 
var tm *tunnelManager 

/***
 * tunnel上的数据包格式  4 byte id + 2 byte date len + 1 byte addr + other
 ***/

/* tunnel 初始化 */
func Init() {
	CONF = new(TunnelConf)
	CONF.InitConns = 0
	tm = newTunnelManager()
}

func SelectTunnel() int {
	return 0
}

/* 发送对端关闭信息 并移除map中的item */
func RemoveSession(conn net.Conn) {
	delete(tm.connMp, conn)	
}

func ProcessData(onDataf func(net.Conn, []byte), conn net.Conn, rawData []byte, dst string) int {
	var s *session
	if tm.connMp[conn] == nil {
		s = tm.allocSession(conn)	
		tm.connMp[conn] = s
	} else {
		s = tm.connMp[conn]
	}	 
	
	// 打包数据
	onDataf(conn, rawData)
	// 发送数据
	return 0
}

func ProcessTunnel(conn net.Conn) {
	
}

/************** 发送打好包的数据  ****************/
func (t *tunnel) sendData(packData []byte, len int) {
	if t.state == T_OPEN {
		t.send <- packData[:len]
	}
}

/************** tunnel读数据 ********************/
func (t *tunnel) processRead() {
	var buf [bufferSize]byte
	for {
		// 连接断开
		_, err := t.Read(buf[:5])
		if err != nil {
			log.Println(err)
		}
		
		// 解析数据 得到 session id
	
		// 读取数据 发送给session[id]
	}
	
}

/************** tunnel写数据 ********************/
func (t *tunnel) processWrite() {
	var count int
	for {
		// 读取 channel
		buf, ok := <- t.send
		if !ok {
			
			break
		}
		// 写数据
		count = 0
		for {
			n, err := t.Write(buf[count:])
			if err != nil {
				t.Lock()
				defer t.Unlock()
				t.state = T_CLOSE
			} else if (n != len(buf)) {
				count = n
			}
		}
		
	}
}


/************** 初始化tunnel ********************/
func (t *tunnel) init() {
	addr, err1 := net.ResolveTCPAddr("tcp", serverAddr)
	if err1 != nil {
		log.Fatal(err1)
	}

	_, err2 := net.DialTCP("tcp", nil, addr)
	if err2 != nil {
		log.Fatal(err2)
	}
	//t.conn = conn
	
	go t.processRead()
	go t.processWrite()
}


/*********************
 *** tunnelManager ***
 *********************/


func newTunnelManager() *tunnelManager {
	log.Println("newTunnelManager")
	tm := new(tunnelManager)
	tm.List = list.New()
	/*for i := 0; i < maxClientConns; i++ {
		t := &tm.ts[i]
		t.id = i
		t.Element = tm.PushBack(t)
	}*/
	tm.tunnelCount = CONF.InitConns
	tm.connMp = map[net.Conn]*session{}
	tm.idMp = map[uint32]*session{}
	tm.sessionIdIndex = 1
	// 初始化tunnel
	for i := 0; i < tm.tunnelCount; i++ {
		t := &tm.ts[i]
		go t.init()
	}
	return tm
}

func (tm * tunnelManager) allocSession(conn net.Conn) *session {
	tm.Lock()
	defer tm.Unlock()
	f := tm.Front()
	if f == nil {
		return nil
	}
	s := tm.Remove(f).(*session)
	s.Element = nil
	s.id = tm.sessionIdIndex
	tm.sessionIdIndex++
	tm.idMp[s.id] = s
	s.c = conn
	return s
}

func (tm * tunnelManager) getSession(id uint32) *session {
	tm.Lock()
	defer tm.Unlock()
	
	return tm.idMp[id]
}

func (tm * tunnelManager) freeSession(id uint32) {
	tm.Lock()
	defer tm.Unlock()
	s := tm.idMp[id]
	if s != nil && s.Element == nil {
		/*t.sendClose()
		s.Element = tm.PushBack(t)
		t.close()*/
	}
}
/*
func (s xsocket) Read(buf []byte) (n int, err os.Error) {
	n, err = io.ReadFull(s.Conn, buf)
	return 
}

func (s xsocket) Write(data []byte) (n int, err os.Error) {
	log.Println("Write data len", len(data))
	
}*/
