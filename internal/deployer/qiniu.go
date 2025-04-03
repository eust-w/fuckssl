package deployer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fuckssl/internal/provider"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
)

type QiniuDeployer struct {
	accessKey string
	secretKey string
	bucket    string
	domain    string
}

func NewQiniuDeployer(accessKey, secretKey, bucket, domain string) (*QiniuDeployer, error) {
	return &QiniuDeployer{
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		domain:    domain,
	}, nil
}

func (d *QiniuDeployer) Name() string {
	return "qiniu"
}

func (d *QiniuDeployer) Deploy(cert *provider.Certificate) error {
	// 1. 上传证书
	certID, err := d.uploadSSL(cert)
	if err != nil {
		return fmt.Errorf("failed to upload SSL certificate: %v", err)
	}

	// 2. 为域名配置证书
	err = d.replaceSSL(certID)
	if err != nil {
		return fmt.Errorf("failed to configure SSL for domain: %v", err)
	}

	return nil
}

func (d *QiniuDeployer) uploadSSL(cert *provider.Certificate) (string, error) {
	// 创建认证信息
	mac := qbox.NewMac(d.accessKey, d.secretKey)

	// 准备请求数据
	url := "https://api.qiniu.com/sslcert"
	reqData := map[string]interface{}{
		"name":        fmt.Sprintf("ssl%snew", time.Now().Format("20060102")),
		"common_name": d.domain,
		"ca":          string(cert.Certificate),
		"pri":         string(cert.PrivateKey),
	}

	// 发送请求
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "api.qiniu.com")

	// 获取认证信息
	token, err := mac.SignRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %v", err)
	}
	req.Header.Set("Authorization", "QBox "+token)

	// 创建 HTTP 客户端，支持代理设置
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// 解析响应
	var result struct {
		CertID string `json:"certID"`
		Code   int    `json:"code"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Code != 200 {
		return "", fmt.Errorf("upload failed with code: %d", result.Code)
	}

	return result.CertID, nil
}

func (d *QiniuDeployer) replaceSSL(certID string) error {
	// 创建认证信息
	mac := qbox.NewMac(d.accessKey, d.secretKey)

	// 准备请求数据
	url := fmt.Sprintf("https://api.qiniu.com/domain/%s/httpsconf", d.domain)
	reqData := map[string]interface{}{
		"certid":      certID,
		"forceHttps":  true,
		"http2Enable": true,
	}

	// 发送请求
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "api.qiniu.com")

	// 获取认证信息
	token, err := mac.SignRequest(req)
	if err != nil {
		return fmt.Errorf("failed to sign request: %v", err)
	}
	req.Header.Set("Authorization", "QBox "+token)

	// 创建 HTTP 客户端，支持代理设置
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// 解析响应
	var result struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Code != 200 {
		return fmt.Errorf("configuration failed with code: %d", result.Code)
	}

	return nil
}
