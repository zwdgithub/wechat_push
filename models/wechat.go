package models

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	AccessToken    = ""
	pushMessageUrl = beego.AppConfig.String("push_message_url")
)

type AccessToKen struct {
	AccessToKen string
	expiresIn   int
}

type PushMsgResp struct {
	Errcode int
	Errmsg  string
	Msgid   int
}

func RefreshAccessToken() {
	log.Println("start get access_token...")
	var accessToken AccessToKen
	grantType := beego.AppConfig.String("access_token_grant_type")
	appid := beego.AppConfig.String("access_token_appid")
	secret := beego.AppConfig.String("access_token_secret")
	accessTokenUrl := beego.AppConfig.String("access_token_url")
	client := &http.Client{Timeout: 50 * time.Second}
	var param = url.Values{}
	param.Add("grant_type", grantType)
	param.Add("appid", appid)
	param.Add("secret", secret)
	request, _ := http.NewRequest("GET", accessTokenUrl, strings.NewReader(param.Encode()))
	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &accessToken)
	AccessToken = accessToken.AccessToKen
	log.Println("get access_token end...")
}

func PushMsg(msg, desc, to string) (bool, string) {
	client := &http.Client{Timeout: 10 * time.Second}
	data := map[string]interface{}{
		"touser":      to,
		"template_id": "qoiof9xKCUKjst0cS5EjJ2LFhet_Z70RH_RzKXrtVm8",
		"data": map[string]interface{}{
			"first": map[string]string{
				"value": "安全风险通知推送",
				"color": "#173177",
			},
			"keyword1": map[string]interface{}{
				"value": time.Now(),
				"color": "#173177",
			},
			"keyword2": map[string]string{
				"value": desc,
				"color": "#173177",
			},
			"keyword3": map[string]string{
				"value": msg,
				"color": "#173177",
			},
			"remark": map[string]string{
				"value": "",
				"color": "#173177",
			},
		},
	}
	param, _ := json.Marshal(data)
	reader := bytes.NewReader(param)
	request, _ := http.NewRequest("POST", pushMessageUrl+AccessToken, reader)
	response, err := client.Do(request)
	if err != nil {
		log.Printf("message push failed : %s", err)
		log.Printf("faild message info -> msg: %s, desc: %s, to: %s", msg, desc, to)
		return false, "推送失败"
	}
	var resp PushMsgResp
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &resp)
	if resp.Errcode != 0 {
		log.Printf("push message failed -> msg: %s, desc: %s, to: %s", msg, desc, to)
		return false, "推送失败"
	}
	log.Printf("push message success -> msg: %s, desc: %s, to: %s", msg, desc, to)
	return true, ""
}
