package util

import (	
	"io"
	"os"
	"net"
	"fmt"
	"flag"
	"time"
	"errors"
	"crypto/md5"
	"encoding/hex"	
	"encoding/binary"
)

const MAX_STRING = 10240

type Command struct {
	CmdType string
	HWID    string
	VMID    string
	Port    string
	Data    string
	Host    string
	ServerType string
}

//消息类型
const (
	MSG_TYPE_CMD       = iota //普通命令
	MSG_TYPE_INFO             //RSA验证信息
	MSG_TYPE_RES              //动态资源信息
	MSG_TYPE_LOG              //告警信息
	MSG_TYPE_VERSION          //版本信息
	MSG_TYPE_KEEPALIVE        //心跳
	MSG_TYPE_DATA 				//数据统计信息
	MSG_TYPE_INVALID          //无效类型
)

//连接类型
const (
	TOKEN_LEN       = 4
	C2P_CONNECT     = "C2P0"
	C2P_SESSION     = "C2P1"
	C2P_KEEP_ALIVE  = "C2P2"
	P2C_NEW_SESSION = "P2C1"
	C2P_REPORT      = "C2P3"
	SEPS            = "\n"
)

//功能：获取随机盐值，御用加强密码存储安全性
//参数：
//*
//返回值：
//*
//说明：
func GetSalt() string {
	h := md5.New()
	h.Write([]byte(time.Now().String())) //

	salt := hex.EncodeToString(h.Sum(nil)) //gen_hw_id()
	return salt
}
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

//打印代理连接联通的两端
func Conn2Str(conn net.Conn) string {
	return conn.LocalAddr().String() + " <-> " + conn.RemoteAddr().String()
}

//功能：网络解包传输函数
//参数：
//*r      :读IO对象
//返回值：
//*[]byte ：消息数据
//*int32  : 消息类型
//*error  ：错误
//说明：包格式（32位长度|消息类型|数据）
func ReadMsg(r io.Reader) ([]byte, int32, error) {
	var size int32
	var t int32
	err := binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return nil, MSG_TYPE_INVALID, err
	}
	err = binary.Read(r, binary.LittleEndian, &t)
	if err != nil {
		return nil, MSG_TYPE_INVALID, err
	}
	//beego.Debug("Recive ", size, "bytes data")
	if size > MAX_STRING {
		return nil, MSG_TYPE_INVALID, errors.New("Too Long String")
	}
	buff := make([]byte, size)
	n, err := r.Read(buff)
	if err != nil {
		return nil, MSG_TYPE_INVALID, err
	}
	if int32(n) != size {
		return nil, MSG_TYPE_INVALID, errors.New("Invalid String Size")
	}
	return buff, t, nil
}

//功能：网络封包传输函数
//参数：
//*w   :写IO对象
//*buf ：消息数据
//*t   : 消息类型
//返回值：
//*int  ：传输数据字节数
//*error ：错误

//说明：包格式（32位长度|消息类型|数据）
func WriteMsg(w io.Writer, buf []byte, t int32) (int, error) {
	//beego.Debug("Send ", len(buf), "bytes data")
	binary.Write(w, binary.LittleEndian, int32(len(buf)))
	binary.Write(w, binary.LittleEndian, int32(t))
	return w.Write(buf)
}