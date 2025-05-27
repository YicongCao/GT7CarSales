package fetcher

import (
	"errors"
	"io"
	"net/http"
)

// FetchJSONFromURL 从指定 URL 获取 JSON 响应内容并返回字节切片
func FetchJSONFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fetch failed: " + resp.Status)
	}
	return io.ReadAll(resp.Body)
}
