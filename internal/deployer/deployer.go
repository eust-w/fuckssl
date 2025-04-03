package deployer

import (
	"fmt"
	"fuckssl/internal/provider"
	"sync"
)

// Deployer 定义证书部署者接口
type Deployer interface {
	// Deploy 部署证书
	Deploy(cert *provider.Certificate) error
	// Name 返回部署者名称
	Name() string
}

var (
	deployers = make(map[string]Deployer)
	mu        sync.RWMutex
)

// Register 注册新的证书部署者
func Register(d Deployer) {
	mu.Lock()
	defer mu.Unlock()
	deployers[d.Name()] = d
}

// GetDeployer 获取指定名称的证书部署者
func GetDeployer(name string) (Deployer, error) {
	mu.RLock()
	defer mu.RUnlock()

	d, exists := deployers[name]
	if !exists {
		return nil, fmt.Errorf("deployer %s not found", name)
	}
	return d, nil
}

// ListDeployers 列出所有可用的证书部署者
func ListDeployers() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(deployers))
	for name := range deployers {
		names = append(names, name)
	}
	return names
}
