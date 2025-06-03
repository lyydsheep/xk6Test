package handler

import (
	"email/common/log"
	"email/common/util"
	"email/dal/model"
	"email/logic/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserHandler struct {
	svc *service.EmailTaskSubService
}

func NewUserHandler(svc *service.EmailTaskSubService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) UnsubscribeHandler(ctx *gin.Context) {
	// TODO binding
	var req UnsubscribeReq
	if err := ctx.ShouldBind(&req); err != nil {
		log.New(ctx).Error("binding error", "err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong parameters",
		})
	}
	log.New(ctx).Debug("receive unsubscribe request", "req", req)
	key, err := util.DecodeBase64(req.Data.QueryString.Key)
	if err != nil {
		log.New(ctx).Error("fail to decode key", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to decode key",
		})
	}
	if err := h.svc.Unsubscribe(ctx, &model.EmlUnsubscribeUsrUser{
		Spm: req.Data.QueryString.Spm,
		IP:  req.Data.Headers.XTrueIP,
	}, key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "unsubscribe success",
	})
}

type UnsubscribeReq struct {
	Data struct {
		Body    string `json:"body"`
		Headers struct {
			ContentLength           string `json:"content-length"`
			XRequestID              string `json:"X-Request-ID"`
			XOriginalForwardedFor   string `json:"X-Original-Forwarded-For"`
			X5Uuid                  string `json:"x5-uuid"`
			UserAgent               string `json:"User-Agent"`
			SecFetchDest            string `json:"Sec-Fetch-Dest"`
			AcceptEncoding          string `json:"Accept-Encoding"`
			SecFetchMode            string `json:"Sec-Fetch-Mode"`
			SecChUaMobile           string `json:"sec-ch-ua-mobile"`
			EagleEyeTraceId         string `json:"EagleEye-TraceId"`
			UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
			XForwardedCluster       string `json:"X-Forwarded-Cluster"`
			WLProxyClientIP         string `json:"WL-Proxy-Client-IP"`
			SecFetchUser            string `json:"Sec-Fetch-User"`
			XRealIP                 string `json:"X-Real-IP"`
			Accept                  string `json:"Accept"`
			XForwardedHost          string `json:"X-Forwarded-Host"`
			XForwardedProto         string `json:"X-Forwarded-Proto"`
			SecFetchSite            string `json:"Sec-Fetch-Site"`
			Host                    string `json:"Host"`
			XForwardedPort          string `json:"X-Forwarded-Port"`
			XClientIP               string `json:"X-Client-IP"`
			SecChUa                 string `json:"sec-ch-ua"`
			SecChUaPlatform         string `json:"sec-ch-ua-platform"`
			XForwardedFor           string `json:"X-Forwarded-For"`
			AcceptLanguage          string `json:"Accept-Language"`
			EagleeyeRpcid           string `json:"eagleeye-rpcid"`
			WebServerType           string `json:"Web-Server-Type"`
			XSinfo                  string `json:"X-Sinfo"`
			XScheme                 string `json:"X-Scheme"`
			XTrueIP                 string `json:"X-True-IP"`
		} `json:"headers"`
		HttpMethod  string `json:"httpMethod"`
		Path        string `json:"path"`
		QueryString struct {
			Spm string `json:"spm"`
			Key string `json:"key"`
		} `json:"queryString"`
	} `json:"data"`
	Id                      string    `json:"id"`
	Source                  string    `json:"source"`
	Specversion             string    `json:"specversion"`
	Type                    string    `json:"type"`
	Datacontenttype         string    `json:"datacontenttype"`
	Time                    time.Time `json:"time"`
	Subject                 string    `json:"subject"`
	Aliyunaccountid         string    `json:"aliyunaccountid"`
	Aliyunpublishtime       time.Time `json:"aliyunpublishtime"`
	Aliyuneventbusname      string    `json:"aliyuneventbusname"`
	Aliyunregionid          string    `json:"aliyunregionid"`
	Aliyunoriginalaccountid string    `json:"aliyunoriginalaccountid"`
	Aliyunpublishaddr       string    `json:"aliyunpublishaddr"`
}
