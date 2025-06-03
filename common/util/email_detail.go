package util

import (
	"email/common/enum"
	"errors"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"time"
)

// FetchResp accountName 和 toAddress 不能同时为空或同时不为空
func FetchResp(accountName, toAddress, nextStart, startTime, endTime string, client *dm20151123.Client) (*dm20151123.SenderStatisticsDetailByParamResponse, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	// 获取邮件详情
	senderStatisticsDetailByParamRequest := &dm20151123.SenderStatisticsDetailByParamRequest{}
	senderStatisticsDetailByParamRequest.StartTime = tea.String(time.Now().Add(-29 * time.Hour * 24).Format(enum.TimeFormatHyphenedYMDHI))
	senderStatisticsDetailByParamRequest.EndTime = tea.String(time.Now().Format(enum.TimeFormatHyphenedYMDHI))

	if accountName != "" {
		senderStatisticsDetailByParamRequest.AccountName = tea.String(accountName)
	}
	if toAddress != "" {
		senderStatisticsDetailByParamRequest.ToAddress = tea.String(toAddress)
	}
	if nextStart != "" {
		senderStatisticsDetailByParamRequest.NextStart = tea.String(nextStart)
	}
	if startTime != "" {
		senderStatisticsDetailByParamRequest.StartTime = tea.String(startTime)
	}
	if endTime != "" {
		senderStatisticsDetailByParamRequest.EndTime = tea.String(endTime)
	}

	runtime := &util.RuntimeOptions{}
	res, tryErr := func() (*dm20151123.SenderStatisticsDetailByParamResponse, error) {
		// 复制代码运行请自行打印 API 的返回值
		resp, err := client.SenderStatisticsDetailByParamWithOptions(senderStatisticsDetailByParamRequest, runtime)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}()

	if tryErr != nil {
		return nil, tryErr
	}

	return res, nil
}
