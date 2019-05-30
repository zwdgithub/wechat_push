package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
	"time"
	"wechat_push/models"
)

const (
	TOCKEN  = "ysuejeeeiq123"
	REPLAEY = `<xml>
			  <ToUserName><![CDATA[%s]]></ToUserName>
			  <FromUserName><![CDATA[%s]]></FromUserName>
			  <CreateTime>%d</CreateTime>
			  <MsgType><![CDATA[%s]]></MsgType>
			  <Content><![CDATA[%s]]></Content>
			</xml>`
	FROM_USER = ""
)

var (
	UserTxtFile = beego.AppConfig.String("user_txt_file")
	bindMutex   sync.Mutex
	USERS       = make(map[string]string)
)

type MainController struct {
	beego.Controller
}

type Message struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string
	MsgId        string
}

/**
微信绑定校验
*/
func (this *MainController) CheckSignature() {
	s := this.GetString("signature")
	timestamp := this.GetString("timestamp")
	nonce := this.GetString("nonce")
	echoStr := this.GetString("echostr")
	if signature(timestamp, nonce) == s {
		this.Ctx.WriteString(echoStr)
	}
	this.Ctx.WriteString("")
}

/**
接收信息处理
*/
func (this *MainController) Message() {
	var msg Message
	err := xml.Unmarshal(this.Ctx.Input.RequestBody, &msg)
	if err != nil {
		log.Fatalln("Message xml build error")
	}
	log.Printf("user openid is %s,  msg is %s", msg.FromUserName, msg.Content)
	reMsg := command(msg.Content, msg.FromUserName, "")
	replay := createReplay(msg.FromUserName, msg.ToUserName, "text", reMsg)
	this.Ctx.WriteString(replay)
}

//根据 timestamp, nonce 生成signature
func signature(timestamp, nonce string) string {
	list := []string{timestamp, nonce, TOCKEN}
	sort.Strings(list)
	return sh1(fmt.Sprintf("%s%s%s", list[0], list[1], list[2]))
}

//回复信息内容填充
func createReplay(to, from, msgType, content string) string {
	return fmt.Sprintf(REPLAEY, to, from, time.Now(), msgType, content)
}

func (this *MainController) PushMsg() {
	respMsg := "推送失败"
	key := this.GetString("key")
	msg := this.GetString("msg")
	desc := this.GetString("desc")
	if to, ok := USERS[key]; !ok {
		b, s := models.PushMsg(msg, desc, to)
		if b {
			this.Data["json"] = map[string]interface{}{"status": 1, "msg": "推送成功"}
			this.ServeJSON()
			this.StopRun()
		}
		respMsg = s
	} else {
		respMsg = "请先绑定微信号"
	}
	this.Data["json"] = map[string]interface{}{"status": -1, "msg": respMsg}
	this.ServeJSON()
}

func UserInit() {
	bytes, _ := ioutil.ReadFile(UserTxtFile)
	_ = json.Unmarshal(bytes, &USERS)
}

func command(cmd, from, msg string) string {
	switch cmd {
	case "绑定":
		bindMutex.Lock()
		defer bindMutex.Unlock()
		USERS[from] = sh1(from)
		USERS[USERS[from]] = from
		bytes, _ := json.Marshal(USERS)
		ioutil.WriteFile(UserTxtFile, bytes, os.ModeAppend)
		return "绑定成功"
	case "解绑":
		bindMutex.Lock()
		defer bindMutex.Unlock()
		delete(USERS, USERS[from])
		delete(USERS, from)
		bytes, _ := json.Marshal(USERS)
		ioutil.WriteFile(UserTxtFile, bytes, os.ModeAppend)
		return "解绑成功"
	case "获取":
		if data, ok := USERS[from]; ok {
			return data
		}
		return "请先绑定帐号"
	default:
		return "你在说啥，我听不清？"
	}

}

func sh1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
