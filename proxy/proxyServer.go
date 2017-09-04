package proxy

import (
	"quickstart/util"
	"net"
	"time"
	"crypto/tls"
	"github.com/astaxie/beego"
	"github.com/vmihailenco/msgpack"
)

const CONF_PATH = "/zzw/test/"

type OnConnectFunc func(net.Conn)

func Listen(port string, onConnect OnConnectFunc) error {
	cert, err := tls.LoadX509KeyPair(CONF_PATH+"conf/ssl.crt", CONF_PATH+"conf/ssl.key")
	if err != nil {
		beego.Error("RSA File Missing:", err)
		return err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	server, err := tls.Listen("tcp", net.JoinHostPort("0.0.0.0", port), config)
	if err != nil {
		beego.Error("Can't Listen:", err)
		return err
	}
	beego.Info("Proxy Server Running on :", port)
	go func() {
		defer server.Close()
		for {
			conn, err := server.Accept() //循环接受客户端和设备的连接请求
			if err != nil {
				beego.Error("Can't Accept: ", err)
				continue
			}
			go onConnect(conn)
		}
	}()
	return nil
}

//处理各种反向连接请求
func onCommandConnect(conn net.Conn) {

	strConn := util.Conn2Str(conn)
	beego.Debug("ProxyClient Connect:", strConn)

	conn.SetReadDeadline(time.Now().Add(15 * time.Second))

	ticket, t, err := util.ReadMsg(conn)
	if err != nil {
		beego.Info("From:", conn.RemoteAddr().String(), "t Identity: ", err)
		conn.Close()
		return
	}

	beego.Info("t:%s", t);
	beego.Info("ticket:%s", ticket);


	var req util.Command

	//读取请求信息
	buf, _, err := util.ReadMsg(conn) 

	conn.SetReadDeadline(time.Time{})

	if err != nil {
		beego.Error("Read Err: ", err)
		conn.Close()
		return
	}

	err = msgpack.Unmarshal(buf, &req)
	if err != nil {
		beego.Error("Unmarshal  Err:", err)
	}

}

//启动内网代理服务
func StartServer(port string) error {
	err := Listen(port, onCommandConnect)
	if nil != err {
		beego.Error("StartServer:", err)
		return err
	}
	return nil
}