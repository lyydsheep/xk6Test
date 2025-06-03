package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"email/common/enum"
	"email/common/log"
	"email/config"
	"email/event"
	"email/logic/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
	"time"
)

// TODO  测试使用
const (
	ExchangeName   = "pay"
	PaidQueueName  = "crm_email"
	PaidBindingKey = "2.v1.user.paid"
)

var (
	conn *amqp.Connection
	once sync.Once
)

func InitMQ() error {
	var err error
	conn, err = connect()
	if err != nil {
		return err
	}
	if err = bindingExchangeAndQueue(); err != nil {
		return err
	}
	return nil
}

type EmailEventService struct {
	EmailTaskSubRepo repository.EmailTaskSubRepo
	TaskConfigRepo   repository.TaskConfigRepo
	EmailTemplateRepo repository.EmailTemplateRepo
}

func NewEmailEventService(repo repository.EmailTaskSubRepo, taskConfigRepo repository.TaskConfigRepo, emailTemplateRepo repository.EmailTemplateRepo) *EmailEventService {
	return &EmailEventService{
		EmailTaskSubRepo: repo,
		TaskConfigRepo:   taskConfigRepo,
		EmailTemplateRepo: emailTemplateRepo,
	}
}

func (svc *EmailEventService) Start(ctx context.Context) {
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	for {
		err := InitMQ()
		if err != nil {
			log.New(ctx).Error("Failed to connect to RabbitMQ", "err", err)
			log.New(ctx).Info("retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
		log.New(ctx).Info("starting to consume email event")
		ch, err := conn.Channel()
		if err != nil {
			log.New(ctx).Error("Failed to create RabbitMQ channel", "err", err)
			log.New(ctx).Info("retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
		msgs, err := ch.Consume(PaidQueueName, "", false, false, false, false, nil)
		if err != nil {
			log.New(ctx).Error("Failed to consume messages", "err", err)
			log.New(ctx).Info("retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
		// 发送兑换码
		var eg errgroup.Group
		eg.Go(func() error {
			return svc.createCodeTask(ctx, msgs)
		})
		if eg.Wait() != nil {
			log.New(ctx).Error("Failed to consume messages", "err", err)
			log.New(ctx).Info("retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
		if conn != nil {
			conn.Close()
		}
	}
}

// 支付成功  --->  触发一个邮件任务   ---->   根据模板

func (svc *EmailEventService) createCodeTask(ctx context.Context, msgs <-chan amqp.Delivery) error {
	// 使用指数退避策略处理消息
	consume := func(body []byte) error {
		var paid event.PaidEvent
		err := json.Unmarshal(body, &paid)
		log.New(ctx).Info("[createCodeTask]consume message", "message", string(body))
		if err != nil {
			log.New(ctx).Error("[createCodeTask]Failed to unmarshal message", "message", string(body), "err", err)
			return err
		}
		operation := func() error {
			// 尽可能保证 eos
			priority, err := svc.EmailTemplateRepo.ReadPriority(ctx, 2)
			//taskConfig, err := svc.TaskConfigRepo.Read(ctx, paid.Account.Cid, enum.Email, enum.TaskConfigRedemption)
			//log.New(ctx).Info("[createCodeTask]taskConfig", "taskConfig", taskConfig)
			log.New(ctx).Info("[createCodeTask]priority", "priority", priority)
			if err != nil {
				log.New(ctx).Error("[createCodeTask]Fail to read task config", "cid", paid.Account.Cid, "category", enum.Email, "type", enum.TaskConfigRedemption, "error", err)
				return err
			}
			isUnknownGoods, err := svc.EmailTaskSubRepo.CreateWithEvent(ctx, &paid, priority)
			var mysqlErr *mysql.MySQLError
			if err != nil {
				log.New(ctx).Error("[createCodeTask]Failed to create email task", "message", string(body), "err", err)
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					log.New(ctx).Warn("[createCodeTask]Duplicate entry", "message", string(body), "err", err)
					return nil
				}
				if isUnknownGoods {
					dingTalk.notifyInfo(fmt.Sprintf("兑换码任务创建失败. goodsid未定义: %+v", &paid))
					return nil
				}
			}
			return err
		}

		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.MaxElapsedTime = 1 * time.Minute

		err = backoff.Retry(operation, expBackoff)
		if err != nil {
			return err
		}
		return nil
	}

	for d := range msgs {
		if err := consume(d.Body); err != nil {
			log.New(ctx).Error("Failed to consume message", "message", string(d.Body), "err", err)
			d.Nack(false, true)
		} else {
			log.New(ctx).Info("Message consumed successfully", "message", string(d.Body))
			d.Ack(false)
		}
	}

	return nil
}

func connect() (*amqp.Connection, error) {
	var conn *amqp.Connection

	operation := func() error {
		// 加载 CA 证书
		caCert, err := os.ReadFile(config.MQ.CAFilePath)
		if err != nil {
			log.New(context.TODO()).Error("Error reading CA certificate", "err", err)
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// 加载客户端证书
		clientCert, err := tls.LoadX509KeyPair(config.MQ.CerFilePath, config.MQ.KeyFilePath)
		if err != nil {
			log.New(context.TODO()).Error("Error loading client certificate", "err", err)
			return err
		}

		// 创建 TLS 配置
		tlsConfig := &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{clientCert},
		}

		// 使用 TLS 建立连接
		conn, err = amqp.DialTLS(config.MQ.Url, tlsConfig)
		return err
	}

	// 指数退避重试策略
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = time.Minute

	err := backoff.Retry(operation, expBackoff)
	if err != nil {
		return nil, err
	}

	log.New(context.TODO()).Info("Successfully connected to RabbitMQ")
	return conn, nil
}

// 已经在控制面板 binding 了队列，不需要再创建
func bindingExchangeAndQueue() error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		PaidQueueName, true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		PaidQueueName, PaidBindingKey, ExchangeName, false, nil,
	)
	if err != nil {
		return err
	}
	return nil
}
