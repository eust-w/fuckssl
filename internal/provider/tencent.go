package provider

import (
	"fmt"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type TencentProvider struct {
	client *ssl.Client
}

func NewTencentProvider(secretId, secretKey string) (*TencentProvider, error) {
	credential := common.NewCredential(secretId, secretKey)
	prof := profile.NewClientProfile()

	client, err := ssl.NewClient(credential, "ap-guangzhou", prof)
	if err != nil {
		return nil, fmt.Errorf("failed to create tencent client: %v", err)
	}

	return &TencentProvider{client: client}, nil
}

func (p *TencentProvider) Name() string {
	return "tencent"
}

func (p *TencentProvider) Apply(domain string) (*Certificate, error) {
	// 创建证书申请请求
	request := ssl.NewApplyCertificateRequest()
	request.DvAuthMethod = common.StringPtr("DNS_AUTO") // 使用自动 DNS 验证
	request.DomainName = common.StringPtr(domain)
	request.PackageType = common.StringPtr("83")   // TrustAsia C1 DV Free
	request.ValidityPeriod = common.StringPtr("3") // 3个月有效期
	request.CsrEncryptAlgo = common.StringPtr("RSA")
	request.CsrKeyParameter = common.StringPtr("2048")
	request.ContactEmail = common.StringPtr("ssl@tencent.com") // 默认邮箱
	request.ContactPhone = common.StringPtr("18888888888")     // 默认手机号
	request.Alias = common.StringPtr(domain)                   // 使用域名作为别名

	// 申请证书
	response, err := p.client.ApplyCertificate(request)
	if err != nil {
		return nil, fmt.Errorf("failed to apply certificate: %v", err)
	}

	// 获取证书详情
	certId := response.Response.CertificateId
	detailRequest := ssl.NewDescribeCertificateDetailRequest()
	detailRequest.CertificateId = certId

	// 等待证书签发完成
	for {
		detailResponse, err := p.client.DescribeCertificateDetail(detailRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to get certificate detail: %v", err)
		}

		if *detailResponse.Response.Status == 1 {
			return &Certificate{
				Domain:      domain,
				Certificate: []byte(*detailResponse.Response.CertificatePublicKey),
				PrivateKey:  []byte(*detailResponse.Response.CertificatePrivateKey),
			}, nil
		}

		// 等待一段时间后重试
		time.Sleep(5 * time.Second)
	}
}
