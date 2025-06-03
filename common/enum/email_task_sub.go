package enum

const (
	// 邮件状态
	EmailTaskSubStatusUnsent              = "init"    // 未发送
	EmailTaskSubStatusProcess             = "ing"     // 发送中
	EmailTaskSubStatusFailure             = "fail"    // 发送失败
	EmailTaskSubStatusSent                = "done"    // 已发送
	EmailTaskSubStatusPreDone             = "predone" // 已发送
	EmailTaskSubStatusDoneSmtpRetOvertime = "doneSmtpRetOvertime"
	EmailTaskSubStatusNotFound            = "resultNotFound"
	EmailTaskSubStatusOpened              = "opened"  // 已打开
	EmailTaskSubStatusClicked             = "clicked" // 已点击
)
