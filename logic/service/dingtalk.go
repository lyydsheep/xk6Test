package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	DINGTALK_CBWP_WEBHOOK_URL_NORMAL    = "https://oapi.dingtalk.com/robot/send?access_token=febe8c6502398c86c20f4f8438429a0dde6694b787f5322841b9b9c2afcd983d"
	DINGTALK_CBWP_WEBHOOK_URL_IMPORTANT = "https://oapi.dingtalk.com/robot/send?access_token=4d8d178298ede7d263c38654173d85bef9e2b50cf476eb1405e0cf210e5e78c6"
	DINGTALK_LEVEL                      = int32(2)
)

type DingTalk struct {
	ch             chan DingTalkPkg
	chImportant    chan DingTalkPkg
	HttpClient     http.Client
	CustomKeywords string
}

type DingTalkPkg struct {
	Msg DingTalkMsg
}

type DingTalkMsg struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

var dingTalk *DingTalk

func init() {
	dingTalk = &DingTalk{
		ch:             make(chan DingTalkPkg, 10),
		chImportant:    make(chan DingTalkPkg, 100),
		HttpClient:     http.Client{},
		CustomKeywords: "CBWP",
	}
	go dingTalk.loopImportant()
	go dingTalk.loop()
}
func (d *DingTalk) notifyDebug(content string) {
	d.notify("DEBUG", 3, content, false)
}
func (d *DingTalk) notifyInfo(content string) {
	d.notify("INFO", 2, content, false)
}
func (d *DingTalk) notifyInfoImportant(content string) {
	d.notify("INFO", 2, content, true)
}
func (d *DingTalk) notifyError(content string, err error) {
	d.notify("ERROR", 1, fmt.Sprintf("%s: %v", content, err), true)
}
func (d *DingTalk) notifyFatal(content string, err error) {
	d.notify("FATAL", 0, fmt.Sprintf("%s: %v", content, err), true)
}

func (d *DingTalk) notify(level string, levelInt int32, content string, isImportant bool) {
	log.Printf("[DingTalk] %s", content)
	if levelInt > DINGTALK_LEVEL {
		return
	}
	_content := fmt.Sprintf("[%s][%s]%s: %s", time.Now().Format(time.RFC3339), level, d.CustomKeywords, content)
	pkg := DingTalkPkg{
		Msg: DingTalkMsg{
			MsgType: "text",
		},
	}
	pkg.Msg.Text.Content = _content
	select {
	case dingTalk.ch <- pkg:
		// 发送成功
	default:
		// 通道已满，可以执行其他逻辑或处理
		log.Printf("[DingTalk] 钉钉通知通道已满，无法发送消息 %s", content)
	}

	//重要信息，往重要通知群里单独发送一份
	if isImportant {
		pkg := DingTalkPkg{
			Msg: DingTalkMsg{
				MsgType: "text",
			},
		}
		pkg.Msg.Text.Content = _content
		select {
		case dingTalk.chImportant <- pkg:
			// 发送成功
		default:
			// 通道已满，可以执行其他逻辑或处理
			log.Printf("[DingTalk] 钉钉通知通道已满，无法发送消息 %s", content)
		}
	}
}

func (d *DingTalk) loop() {
	for {
		pkg := <-d.ch
		if DINGTALK_CBWP_WEBHOOK_URL_NORMAL == "" {
			// URL not configured, unable to send DingTalk notification
			continue
		}
		msgJSON, err := json.Marshal(pkg.Msg)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
		}
		req, err := http.NewRequest("POST", DINGTALK_CBWP_WEBHOOK_URL_NORMAL, bytes.NewBuffer(msgJSON))
		if err != nil {
			log.Printf("failed to create request: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := d.HttpClient.Do(req)
		if err != nil {
			log.Printf("failed to send request: %v", err)
			continue
		}
		defer resp.Body.Close()
		time.Sleep(time.Second * 4) //每个机器人每分钟最多发送20条消息到群里，如果超过20条，会限流10分钟。
	}
}

func (d *DingTalk) loopImportant() {
	for {
		pkg := <-d.chImportant
		if DINGTALK_CBWP_WEBHOOK_URL_IMPORTANT == "" {
			// URL not configured, unable to send DingTalk notification
			continue
		}
		msgJSON, err := json.Marshal(pkg.Msg)
		if err != nil {
			log.Printf("failed to marshal message: %v", err)
		}
		req, err := http.NewRequest("POST", DINGTALK_CBWP_WEBHOOK_URL_IMPORTANT, bytes.NewBuffer(msgJSON))
		if err != nil {
			log.Printf("failed to create request: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := d.HttpClient.Do(req)
		if err != nil {
			log.Printf("failed to send request: %v", err)
			continue
		}
		defer resp.Body.Close()
		time.Sleep(time.Second * 4) //每个机器人每分钟最多发送20条消息到群里，如果超过20条，会限流10分钟。
	}
}
