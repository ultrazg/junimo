package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("缺少 Nexus Mods API Key 配置，请在设置中添加")
	}

	return &Client{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
		apiKey: apiKey,
	}, nil
}

func (c *Client) GetJSON(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("apiKey", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		log.Printf("GET 请求 %s 失败: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应体失败: %v", err)
		return err
	}

	log.Printf("GET 请求 %s", url)

	return json.Unmarshal(body, target)
}

func (c *Client) PostJSON(url string, data, target any) error {
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("序列化请求体失败: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("apiKey", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		log.Printf("POST 请求 %s 失败: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应体失败: %v", err)
		return err
	}

	log.Printf("POST 请求 %s", url)

	return json.Unmarshal(body, target)
}
