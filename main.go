package main

import (
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"net/http"
      u "net/url"
	"strings"
	"os"
)

func main() {
	client := &http.Client{}
	SignIn(client)
	
	// success := SignIn(client)
	// if success {
	// 	result := "签到成功"
	// 	fmt.Println(result)
	// 	dingding(result)
	// } else {
	// 	result := "签到失败"
	// 	fmt.Println(result)
	// 	dingding(result)
	// 	os.Exit(3)
	// }
}


// SignIn 签到
func SignIn(client *http.Client) bool {
    //生成要访问的url
    urlStr := "https://www.hifini.com/sg_sign.htm"
    cookie := os.Getenv("COOKIE")
    SIGN_KEY := os.Getenv("SIGN_KEY")
    fmt.Println(SIGN_KEY)
    if cookie == "" {
        fmt.Println("COOKIE不存在，请检查是否添加")
        return false
    }
    if SIGN_KEY == "" {
        fmt.Println("SIGN_KEY不存在，请检查是否添加")
        return false
    }

    //提交请求
    formData := url.Values{}
    formData.Set("sign", SIGN_KEY)

    req, err := http.NewRequest("POST", urlStr, strings.NewReader(formData.Encode()))
    if err != nil {
        panic(err)
    }

    req.Header.Add("Cookie", cookie)
    req.Header.Add("x-requested-with", "XMLHttpRequest")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    //处理返回结果
    response, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer response.Body.Close()
    buf, _ := ioutil.ReadAll(response.Body)
    fmt.Println(string(buf))

    // 钉钉推送
    dingding(string(buf))
    return strings.Contains(string(buf), "成功")
}

func dingding(result string){
	// 构造要发送的消息
	message := struct {
		MsgType string `json:"msgtype"`
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: "HiFiNi：" + result,
		},
	}

	// 将消息转换为JSON格式
	messageJson, _ := json.Marshal(message)
	DINGDING_WEBHOOK := os.Getenv("DINGDING_WEBHOOK")
	// 发送HTTP POST请求
	resp, err := http.Post(DINGDING_WEBHOOK,
		"application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
