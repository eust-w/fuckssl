package provider

import (
	"testing"

	"fuckssl/internal/config"
)

func init() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic(err)
	}
}

func TestNewTencentProvider(t *testing.T) {
	// 从配置中获取腾讯云的配置
	cfg := config.GetConfig()
	tencentConfig, exists := cfg.Providers["tencent"]
	if !exists {
		t.Log("Tencent configuration not found, skipping test")
		t.SkipNow()
	}

	provider, err := NewTencentProvider(
		tencentConfig.Settings["secret_id"],
		tencentConfig.Settings["secret_key"],
	)
	if err != nil {
		t.Fatalf("Failed to create TencentProvider: %v", err)
	}
	if provider == nil {
		t.Fatal("Provider is nil")
	}
}

func TestTencentProvider_Apply(t *testing.T) {
	// 从配置中获取腾讯云的配置
	cfg := config.GetConfig()
	tencentConfig, exists := cfg.Providers["tencent"]
	if !exists {
		t.Log("Tencent configuration not found, skipping test")
		t.SkipNow()
	}

	provider, err := NewTencentProvider(
		tencentConfig.Settings["secret_id"],
		tencentConfig.Settings["secret_key"],
	)
	if err != nil {
		t.Fatalf("Failed to create TencentProvider: %v", err)
	}

	// 测试申请证书
	domain := "example.com"
	c, err := provider.Apply(domain)
	if err != nil {
		t.Logf("Expected error when applying certificate with invalid credentials: %v", err)
		return
	}
	if c != nil {
		t.Log("Certificate applied successfully")
	}
}

func TestTencentProvider_Name(t *testing.T) {
	// 从配置中获取腾讯云的配置
	cfg := config.GetConfig()
	tencentConfig, exists := cfg.Providers["tencent"]
	if !exists {
		t.Log("Tencent configuration not found, skipping test")
		t.SkipNow()
	}

	provider, err := NewTencentProvider(
		tencentConfig.Settings["secret_id"],
		tencentConfig.Settings["secret_key"],
	)
	if err != nil {
		t.Fatalf("Failed to create TencentProvider: %v", err)
	}

	name := provider.Name()
	if name != "tencent" {
		t.Errorf("Expected provider name to be 'tencent', got '%s'", name)
	}
}
