package handler

const (
	Success      = "100000"
	ServerBusy   = "100001"
	ServerErr    = "100002"
	ParameterErr = "100003"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    `json:"data"`
}

type Data struct {
	TaskSubId      int64  `json:"task_sub_id"`
	Classification string `json:"classification"`
	Description    string `json:"description"`
}
