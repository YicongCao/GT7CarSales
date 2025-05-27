package wxwork

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// SendBotMarkdown 向企业微信自建群机器人推送 markdown 消息
func SendBotMarkdown(apikey string, markdown string) error {
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + apikey

	payload := map[string]interface{}{
		"msgtype":  "markdown",
		"markdown": map[string]string{"content": markdown},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("wxwork push failed: " + resp.Status)
	}
	return nil
}
