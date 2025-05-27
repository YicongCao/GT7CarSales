package wxwork

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

// SendBotMarkdown 向企业微信自建群机器人推送 markdown 消息
func SendBotMarkdown(apikey string, markdown string) error {
	const maxLen = 4096

	// 按字节切割，保证不拆分 UTF-8 字符
	slices := splitByByteLen(markdown, maxLen)

	for i, part := range slices {
		url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + apikey

		payload := map[string]interface{}{
			"msgtype": "text",
			"text":    map[string]string{"content": strings.TrimSpace(part)},
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

		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("wxwork resp status: %s\n", resp.Status)
		fmt.Printf("wxwork resp body: %s\n", string(respBody))

		if resp.StatusCode != http.StatusOK {
			return errors.New("wxwork push failed: " + resp.Status)
		}

		// 如果还有下一条，等待 500ms
		if i < len(slices)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return nil
}

// splitByByteLen 按最大字节长度切割字符串，保证不拆分 UTF-8 字符
func splitByByteLen(s string, maxBytes int) []string {
	var result []string
	runes := []rune(s)
	start := 0
	for start < len(runes) {
		byteCount := 0
		end := start
		for end < len(runes) {
			runeLen := utf8.RuneLen(runes[end])
			if byteCount+runeLen > maxBytes {
				break
			}
			byteCount += runeLen
			end++
		}
		result = append(result, string(runes[start:end]))
		start = end
	}
	return result
}
