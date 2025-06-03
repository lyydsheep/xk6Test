package test_test

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSplit(t *testing.T) {
	res := strings.Split("550  5.7.1  https://support.google.com/mail/?p=UnsolicitedMessageError d2e1a72fcca58-74237727d56si10058339b3a.77 - gsmtp", " ")
	fmt.Println(res[2])
}

func TestPointer(t *testing.T) {
	type A struct {
		Name string
	}
	a := new(A)
	a.Name = "123"
	fmt.Printf("%v", a)
}

func TestAPI(t *testing.T) {
	client, err := CreateClient()
	if err != nil {
		panic(err)
	}

	senderStatisticsDetailByParamRequest := &dm20151123.SenderStatisticsDetailByParamRequest{
		//AccountName: tea.String("newsletter@newsletter.wan.video"),
		ToAddress: tea.String("elderl.el7947@gmail.com"),
		// 分页
		//NextStart: tea.String("5f7b268911#203#newsletter@newsletter.wan.video-1747292007#rivas.saul@gmail.com.576461128162268096"),
		StartTime: tea.String(time.Now().Add(-time.Hour * 24 * 29).Format("2006-01-02 15:04")),
		EndTime:   tea.String(time.Now().Format("2006-01-02 15:04")),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, err := client.SenderStatisticsDetailByParamWithOptions(senderStatisticsDetailByParamRequest, runtime)
		if err != nil {
			panic(err)
		}
		var (
			sentTime time.Time
		)
		for _, detail := range resp.Body.Data.MailDetail {
			msec, err := strconv.ParseInt(*detail.LastUpdateTime, 10, 64)
			if err != nil {
				panic(err)
			}
			sentTime = time.UnixMilli(msec)
			fmt.Println(sentTime, *detail.Message)
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, err = util.AssertAsString(error.Message)
		if err != nil {
			panic(err)
		}
	}
}

// Description:
//
// 使用凭据初始化账号Client
//
// @return Client
//
// @throws Exception
func CreateClient() (_result *dm20151123.Client, err error) {
	// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见：https://help.aliyun.com/document_detail/378661.html。
	credential, err := credential.NewCredential(&credential.Config{
		Type:            tea.String("access_key"),
		AccessKeyId:     tea.String(os.Getenv("ACCESS_KEY_ID")),
		AccessKeySecret: tea.String(os.Getenv("ACCESS_KEY_SECRET")),
	})
	if err != nil {
		return _result, err
	}

	config := &openapi.Config{
		Credential: credential,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dm
	config.Endpoint = tea.String("dm.ap-southeast-1.aliyuncs.com")
	_result = &dm20151123.Client{}
	_result, err = dm20151123.NewClient(config)
	return _result, err
}

func TestSortTime(t *testing.T) {
	times := []time.Time{
		time.Now(), time.Now().Add(-time.Hour)}
	slices.SortFunc(times, func(i, j time.Time) int {
		if i.After(j) {
			return -1
		}
		return 1
	})
	fmt.Println(times)
}

func TestTime(t *testing.T) {
	tmp := 2 * time.Second
	fmt.Println(int(tmp / time.Second))
}
