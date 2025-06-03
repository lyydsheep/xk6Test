package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{})
	client, err := CreateClient()
	if err != nil {
		panic(err)
	}

	file, err := os.Open("/Users/lyydsheep/in.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	outFile, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := scanner.Text()
		email = strings.ReplaceAll(email, "|", "")
		email = strings.TrimSpace(email)
		senderStatisticsDetailByParamRequest := &dm20151123.SenderStatisticsDetailByParamRequest{
			//AccountName: tea.String("newsletter@newsletter.wan.video"),
			ToAddress: tea.String(email),
			// 分页
			//NextStart: tea.String("5f7b268911#203#newsletter@newsletter.wan.video-1747292007#rivas.saul@gmail.com.576461128162268096"),
			StartTime: tea.String(time.Now().Add(-time.Hour * 24 * 29).Format("2006-01-02 15:04")),
			EndTime:   tea.String(time.Now().Format("2006-01-02 15:04")),
		}
		runtime := &util.RuntimeOptions{}
		tryErr := func() (_e error) {
			// 复制代码运行请自行打印 API 的返回值
			resp, err := client.SenderStatisticsDetailByParamWithOptions(senderStatisticsDetailByParamRequest, runtime)
			if err != nil {
				panic(err)
			}
			_, err = fmt.Fprintf(outFile, "%+v\n", email)
			if err != nil {
				panic(err)
			}
			for _, detail := range resp.Body.Data.MailDetail {
				if strings.Contains(*detail.Subject, "Code") {
					msec, err := strconv.ParseInt(*detail.LastUpdateTime, 10, 64)
					if err != nil {
						panic(err)
					}
					sentTime := time.UnixMilli(msec).Format("2006-01-02 15:04")
					temp := Detail{
						AccountName:         *detail.AccountName,
						ErrorClassification: *detail.ErrorClassification,
						SentTime:            sentTime,
						Message:             *detail.Message,
						Subject:             *detail.Subject,
						ToAddress:           *detail.ToAddress,
					}
					jsonData, _ := json.Marshal(temp)
					_, err = fmt.Fprintf(outFile, "%+v\n", string(jsonData))
					if err != nil {
						panic(err)
					}
				}
			}
			var res []Result
			db.Raw(`select code_table.code, task_sub.id, task_sub.to_email, task_sub.fetch_time
from eml_redemption_codes as code_table
         right join
     (SELECT *
      FROM eml_task_subs
      WHERE to_email = ? 
        and from_email = 'system@notice.wan.video') as task_sub
     on code_table.task_sub_id = task_sub.id order by task_sub.to_email;`, email).Scan(&res)

			jsonData, err := json.Marshal(res)

			_, err = fmt.Fprintf(outFile, "%+v\n", string(jsonData)+"\n================")
			if err != nil {
				panic(err)
			}
			return nil
		}()
		if tryErr != nil {
			panic(tryErr)
		}
	}
}

type Detail struct {
	AccountName         string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	ErrorClassification string `json:"ErrorClassification,omitempty" xml:"ErrorClassification,omitempty"`
	SentTime            string `json:"LastUpdateTime,omitempty" xml:"LastUpdateTime,omitempty"`
	Message             string `json:"Message,omitempty" xml:"Message,omitempty"`
	Status              int32  `json:"Status,omitempty" xml:"Status,omitempty"`
	Subject             string `json:"Subject,omitempty" xml:"Subject,omitempty"`
	ToAddress           string `json:"ToAddress,omitempty" xml:"ToAddress,omitempty"`
}

type Result struct {
	Code      string     `gorm:"column:code"`
	TaskSubID int64      `gorm:"column:id"`
	ToEmail   string     `gorm:"column:to_email"`
	FetchTime *time.Time `gorm:"column:fetch_time"`
}

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
