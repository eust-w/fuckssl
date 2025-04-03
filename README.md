# FuckSSL

FuckSSL 是一个命令行工具，用于自动从多个云服务商申请 SSL 证书并部署到不同的平台。

## 功能特点

- 支持多个证书提供商（腾讯云、阿里云等）
- 支持多个部署平台（七牛云等）
- 配置管理功能
- 易于扩展的插件架构

## 安装

```bash
go install github.com/yourusername/fuckssl@latest
```

## 配置

1. 创建配置文件：

```bash
cp config.yaml.example ~/.fuckssl/config.yaml
```

2. 编辑配置文件，填入相应的密钥信息。

## 使用方法

### 申请并部署证书

```bash
fuckssl apply-deploy example.com --provider tencent --deployer qiniu
```

### 配置管理

列出所有配置：
```bash
fuckssl config list
```

设置提供商配置：
```bash
fuckssl config set provider tencent --provider-type tencent --settings secret_id=xxx,secret_key=xxx
```

设置部署者配置：
```bash
fuckssl config set deployer qiniu --deployer-type qiniu --settings access_key=xxx,secret_key=xxx,bucket=xxx,domain=xxx
```

查看配置：
```bash
fuckssl config get provider tencent
fuckssl config get deployer qiniu
```

## 支持的提供商

- 腾讯云
- 阿里云

## 支持的部署平台

- 七牛云

## 开发

### 添加新的证书提供商

1. 在 `internal/provider` 目录下创建新的提供商实现
2. 实现 `Provider` 接口
3. 在 `init()` 函数中注册提供商

### 添加新的部署平台

1. 在 `internal/deployer` 目录下创建新的部署平台实现
2. 实现 `Deployer` 接口
3. 在 `init()` 函数中注册部署平台

## 许可证

MIT 