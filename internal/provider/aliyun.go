package provider

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
)

type AliyunProvider struct {
	client *cas.Client
}

func NewAliyunProvider(accessKeyId, accessKeySecret string) (*AliyunProvider, error) {
	client, err := cas.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create aliyun client: %v", err)
	}

	return &AliyunProvider{client: client}, nil
}

func (p *AliyunProvider) Name() string {
	return "aliyun"
}

func (p *AliyunProvider) Apply(domain string) (*Certificate, error) {
	// 创建证书申请请求
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "cas.aliyuncs.com"
	request.Version = "2020-04-07"
	request.ApiName = "CreateCertificateForPackageRequest"

	request.QueryParams["Domain"] = domain
	request.QueryParams["ValidateType"] = "DNS"

	// 申请证书
	response, err := p.client.ProcessCommonRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to apply certificate: %v", err)
	}

	orderId := response.GetHttpContentString()

	// 获取证书详情
	detailRequest := requests.NewCommonRequest()
	detailRequest.Method = "POST"
	detailRequest.Scheme = "https"
	detailRequest.Domain = "cas.aliyuncs.com"
	detailRequest.Version = "2020-04-07"
	detailRequest.ApiName = "DescribeCertificateDetail"
	detailRequest.QueryParams["OrderId"] = orderId

	// 等待证书签发完成
	for {
		detailResponse, err := p.client.ProcessCommonRequest(detailRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to get certificate detail: %v", err)
		}

		status := detailResponse.GetHttpContentString()
		if status == "issued" {
			return &Certificate{
				Domain:      domain,
				Certificate: []byte(detailResponse.GetHttpContentString()),
				PrivateKey:  []byte(detailResponse.GetHttpContentString()),
			}, nil
		}

		// 等待一段时间后重试
		time.Sleep(5 * time.Second)
	}
}
