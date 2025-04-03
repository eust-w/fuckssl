package provider

import (
	"fmt"
	"sync"
)

// Certificate 表示SSL证书信息
type Certificate struct {
	Domain      string
	Certificate []byte
	PrivateKey  []byte
}

// Provider 定义证书提供者接口
type Provider interface {
	// Apply 申请证书
	Apply(domain string) (*Certificate, error)
	// Name 返回提供者名称
	Name() string
}

var (
	providers = make(map[string]Provider)
	mu        sync.RWMutex
)

// Register 注册新的证书提供者
func Register(p Provider) {
	mu.Lock()
	defer mu.Unlock()
	providers[p.Name()] = p
}

// GetProvider 获取指定名称的证书提供者
func GetProvider(name string) (Provider, error) {
	mu.RLock()
	defer mu.RUnlock()

	p, exists := providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return p, nil
}

// ListProviders 列出所有可用的证书提供者
func ListProviders() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}
