package test

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
	"wechat_push/controllers"
	_ "wechat_push/routers"

	"github.com/astaxie/beego"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}


// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	s := `<xml>
  <ToUserName><![CDATA[toUser]]></ToUserName>
  <FromUserName><![CDATA[fromUser]]></FromUserName>
  <CreateTime>1348831860</CreateTime>
  <MsgType><![CDATA[text]]></MsgType>
  <Content><![CDATA[this is a test]]></Content>
  <MsgId>1234567890123456</MsgId>
</xml>`

	var msg controllers.Message
	err := xml.Unmarshal([]byte(s), &msg)
	if err != nil{
		fmt.Print("error")
	}

	timestamp, nonce := "1558974331", "1478497222"
	list := []string{timestamp, nonce, controllers.TOCKEN}
	sort.Strings(list)
	h := sha1.New()
	fmt.Println(fmt.Sprintf("%s%s%s", list[0], list[1], list[2]))
	h.Write([]byte(fmt.Sprintf("%s%s%s", list[0], list[1], list[2])))
	bs := h.Sum(nil)
	fmt.Println(hex.EncodeToString(bs))
}

