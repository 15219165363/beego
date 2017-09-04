package main

import (
	"net"
	"flag"
	"time"
	"runtime"

	"crypto/tls"

	"github.com/astaxie/beego"
	"github.com/vmihailenco/msgpack"	

	"quickstart/util"
)

var(
	Trial_Internal  = 1
	connectStatus   = false
	remote          = flag.String("r", "0.0.0.0:8088", "Address of Server")
)

//功能：main函数
//参数：
//返回值：
//说明：
func main() {

	//设置日志文件
	beego.SetLogger("file", `{"filename":"/zzw/test/client.log"}`)
	beego.SetLevel(7)
	beego.SetLogFuncCall(true)
	//beego.BeeLogger.DelLogger("console")
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Usage = util.Usage
	flag.Parse()

	csc := "192.168.3.243:8088"
	remote = &csc


	//go signalListen() //监听外部信号
	for {

		connectServer()
		beego.Info("After ", Trial_Internal, "Second Try Again!")
		time.Sleep(time.Duration(Trial_Internal) * time.Second)
		//每失败之后多重试一次，时延*2，成功之后重新置为1
		if Trial_Internal < 60 {
			Trial_Internal *= 2
		}
	}
}

//功能：连接服务器处理函数
//参数：
//返回值：
//说明：
func connectServer() {
	hwid := "abcdefg"
	var cmd util.Command
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", *remote, conf)
	if err != nil {
		beego.Error("CAN'T CONNECT:", *remote, " err:", err)
		return
	}
	//发送RSA身份信息

	ticket := "123456"
	_, err = util.WriteMsg(conn, []byte(ticket), util.MSG_TYPE_INFO)
	if err != nil {
		beego.Error("Send RAS identity failed: ", err)
		return
	}

	defer conn.Close()

	command := util.Command{
		CmdType: util.C2P_CONNECT,
		HWID:    hwid,
	}
	//发送设备加入指令
	bufCmd, err := msgpack.Marshal(command)
	if err != nil {
		beego.Error("Marshal Join CMD Err: " + err.Error())
		return
	}
	_, err = util.WriteMsg(conn, bufCmd, util.MSG_TYPE_CMD)
	if err != nil {
		beego.Error("Send Join CMD Msg Err: ", err)
		return
	}


	//等待服务器版本信息反馈,设置超时20S
	conn.SetReadDeadline(time.Now().Add(20 * time.Second))

	buf, t, err := util.ReadMsg(conn)
	if err != nil { //如果超时或出错则尝试重新连接
		beego.Error("Read Server Version Err: ", err)

		return
	}
	beego.Info("bug:", buf)
	beego.Info("t:", t)

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	//连接成功，尝试间隔置为1
	Trial_Internal = 1
	//启动消息上报服务
	//go HandleMsg(conn)
	connectStatus = true

	//压力检测,短时间内突发访问量剧增，则建立更多的session
	leave := make(chan error)
	maxOccurs := 10
	stressTest := func() {
		for {
			select {
			case <-leave:
				return
			case <-time.After(time.Second * 1):
				if maxOccurs > 10 {
					maxOccurs = maxOccurs - 10
				}
			}

		}

	}
	go stressTest()
	for {

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		buf, _, err := util.ReadMsg(conn)

		if err == nil {
			err = msgpack.Unmarshal(buf, &cmd)
			if err != nil {
				beego.Error("UnMarshal Err: " + err.Error())
				continue
			}
			beego.Debug("zzw---cmd.CmdType", cmd.CmdType)
			if cmd.CmdType == util.P2C_NEW_SESSION {

				i := 1
				for i <= maxOccurs { //秒内突发访问量定
					//go session(hwid, local)
					beego.Info("signal connect")
					i = i + 1
				}
				beego.Info("maxOccurs:", maxOccurs)
				if maxOccurs < 50 {
					maxOccurs = maxOccurs + 10
				}

			}  else {

			}
		} else { //
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {

				continue
			} else {
				beego.Error("Server Close, err:", err)
				leave <- nil
				return
			}
		}

	}

}