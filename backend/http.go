package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var msgFlag = map[int]string{
	http.StatusOK:                  "OK",
	http.StatusBadRequest:          "错误请求",
	http.StatusUnauthorized:        "身份验证失败，请检查 Nexus Mods API Key 是否正确",
	http.StatusForbidden:           "拒绝访问",
	http.StatusNotFound:            "404 Not Found",
	http.StatusTooManyRequests:     "Nexus Mods API 请求超过限制，请稍后再试",
	http.StatusInternalServerError: "内部服务器错误",
	http.StatusBadGateway:          "网关错误",
	http.StatusServiceUnavailable:  "服务不可用",
	http.StatusGatewayTimeout:      "网关超时",
}

func getStatusMsg(statusCode int) string {
	msg, ok := msgFlag[statusCode]
	if ok {
		return msg
	}

	return fmt.Sprintf("网络异常: %d", statusCode)
}

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

	hourlyLimit := resp.Header.Get("X-Rl-Hourly-Limit")
	hourlyRemaining := resp.Header.Get("X-Rl-Hourly-Remaining")
	hourlyReset := resp.Header.Get("X-Rl-Hourly-Reset")
	dailyLimit := resp.Header.Get("X-Rl-Daily-Limit")
	dailyRemaining := resp.Header.Get("X-Rl-Daily-Remaining")
	dailyReset := resp.Header.Get("X-Rl-Daily-Reset")

	if resp.StatusCode != http.StatusOK {
		log.Printf("GET 请求 %s 失败：%d %s", url, resp.StatusCode, resp.Status)
		statusMsg := getStatusMsg(resp.StatusCode)
		if resp.StatusCode == http.StatusTooManyRequests {
			statusMsg = fmt.Sprintf(
				"%s，当前小时限制：%s，当前小时剩余：%s，当前小时重置时间：%s，当前天限制：%s，当前天剩余：%s，当前天重置时间：%s",
				statusMsg,
				hourlyLimit, hourlyRemaining, hourlyReset,
				dailyLimit, dailyRemaining, dailyReset,
			)
		}
		return fmt.Errorf("%s", statusMsg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应体失败：%v", err)
		return err
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	data["rate_limit"] = map[string]string{
		"hourly_limit":     hourlyLimit,
		"hourly_remaining": hourlyRemaining,
		"hourly_reset":     hourlyReset,
		"daily_limit":      dailyLimit,
		"daily_remaining":  dailyRemaining,
		"daily_reset":      dailyReset,
	}

	mergedBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	log.Printf("GET 请求 %s", url)

	return json.Unmarshal(mergedBody, target)
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("POST 请求 %s 失败：%d %s", url, resp.StatusCode, resp.Status)
		return fmt.Errorf("%s", getStatusMsg(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应体失败：%v", err)
		return err
	}

	log.Printf("POST 请求 %s", url)

	return json.Unmarshal(body, target)
}
