package deployer

import (
	"testing"

	"fuckssl/internal/config"
	"fuckssl/internal/provider"
)

func init() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic(err)
	}
}

func TestNewQiniuDeployer(t *testing.T) {
	// 从配置中获取七牛云的配置
	cfg := config.GetConfig()
	qiniuConfig, exists := cfg.Deployers["qiniu"]
	if !exists {
		t.Log("Qiniu configuration not found, skipping test")
		t.SkipNow()
	}

	// 使用配置中的值
	deployer, err := NewQiniuDeployer(
		qiniuConfig.Settings["access_key"],
		qiniuConfig.Settings["secret_key"],
		qiniuConfig.Settings["bucket"],
		qiniuConfig.Settings["domain"],
	)
	if err != nil {
		t.Fatalf("Failed to create QiniuDeployer: %v", err)
	}
	if deployer == nil {
		t.Fatal("Deployer is nil")
	}
}

func TestQiniuDeployer_Deploy(t *testing.T) {
	// 从配置中获取七牛云的配置
	cfg := config.GetConfig()
	qiniuConfig, exists := cfg.Deployers["qiniu"]
	if !exists {
		t.Log("Qiniu configuration not found, skipping test")
		t.SkipNow()
	}

	// 使用配置中的值
	deployer, err := NewQiniuDeployer(
		qiniuConfig.Settings["access_key"],
		qiniuConfig.Settings["secret_key"],
		qiniuConfig.Settings["bucket"],
		qiniuConfig.Settings["domain"],
	)
	if err != nil {
		t.Fatalf("Failed to create QiniuDeployer: %v", err)
	}

	// 测试部署证书
	cert := &provider.Certificate{
		Domain:      qiniuConfig.Settings["domain"],
		Certificate: []byte("test_certificate"),
		PrivateKey:  []byte("test_private_key"),
	}

	err = deployer.Deploy(cert)
	if err != nil {
		t.Logf("Deployment error: %v", err)
		return
	}
	t.Log("Certificate deployed successfully")
}

func TestQiniuDeployer_Name(t *testing.T) {
	// 从配置中获取七牛云的配置
	cfg := config.GetConfig()
	qiniuConfig, exists := cfg.Deployers["qiniu"]
	if !exists {
		t.Log("Qiniu configuration not found, skipping test")
		t.SkipNow()
	}

	// 使用配置中的值
	deployer, err := NewQiniuDeployer(
		qiniuConfig.Settings["access_key"],
		qiniuConfig.Settings["secret_key"],
		qiniuConfig.Settings["bucket"],
		qiniuConfig.Settings["domain"],
	)
	if err != nil {
		t.Fatalf("Failed to create QiniuDeployer: %v", err)
	}

	name := deployer.Name()
	if name != "qiniu" {
		t.Errorf("Expected deployer name to be 'qiniu', got '%s'", name)
	}
}
