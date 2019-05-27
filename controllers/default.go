package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"sort"
	"time"
)

const (
	TOCKEN    = "ysuejeeeiq123"
	FROM_USER = ""
	REPLAEY   = `<xml>
			  <ToUserName><![CDATA[%s]]></ToUserName>
			  <FromUserName><![CDATA[%s]]></FromUserName>
			  <CreateTime>%d</CreateTime>
			  <MsgType><![CDATA[%s]]></MsgType>
			  <Content><![CDATA[%s]]></Content>
			</xml>`
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

func (this *MainController) Message() {
	var msg Message
	err := xml.Unmarshal(this.Ctx.Input.RequestBody, &msg)
	if err != nil {
		log.Fatalln("Message xml build error")
	}
	log.Println("user openid is %s,  msg is %s", msg.FromUserName, msg.Content)
	replay := createReplay(msg.FromUserName, msg.ToUserName, "text", "哈哈")
	this.Ctx.WriteString(replay)
}

func signature(timestamp, nonce string) string {
	list := []string{timestamp, nonce, TOCKEN}
	sort.Strings(list)
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s%s%s", list[0], list[1], list[2])))
	bs := h.Sum(nil)
	fmt.Println(string(bs))
	return hex.EncodeToString(bs)
}

func createReplay(to, from, msgType, content string) string {
	return fmt.Sprintf(REPLAEY, to, from, time.Now(), msgType, content)
}
