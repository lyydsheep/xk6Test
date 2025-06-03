package request

type CreateRTEventReq struct {
	Cid               int64  `json:"cid" validate:"required,gt=0"`
	EmlSysUserId      int64  `json:"eml_sys_user_id" validate:"required,gt=0"`
	EmlTemplateId     int64  `json:"eml_template_id" validate:"required,gt=0"`
	UsrUserId         string `json:"usr_user_id"`
	ToEmail           string `json:"to_email" validate:"required,email"`
	TemplateVariables string `json:"template_variables" validate:"required,json"`
	Timeout           int    `json:"timeout" validate:"required,gt=0"`
}
