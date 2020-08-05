package ecode

var (
	OK        = New(1000, "ok")
	BackWrong = New(1001, "not ok")
	ApiError  = New(1003, "获取数据失败")
)

// 通用 请求格式 文件操作... [10000 - 20000]

var (
	MissingRequiredParameters = New(11001, "缺少必要参数")
	ParamsTypeError           = New(11002, "参数类型错误")
	ParamsError               = New(11003, "参数错误")
	JsonMarshalError          = New(11004, "json序列化失败")
	JsonUnmarshalError        = New(11005, "json反序列化失败")
	RENewPassNot              = New(11006, "两次输入密码不一致")
	ParamsTypeChangeError     = New(11007, "参数类型转换错误")

	RedisSaveError = New(12004, "redis保存失败")
	RedisGetError  = New(12005, "redis获取失败")

	GetDbDataFail = New(13001, "数据库未找到相关数据")
)

//业务相关
var (
	UnKnowRecordIsSuccessError = New(33001, "无法判断该团员档案是否已完成审核")
	RecordUnFinishError        = New(33002, "该团员尚未完成团员身份认证，必须认证通过后才能进行组织关系转接。")
	CheckCanTransferFail       = New(33003, "验证该人员是否满足接转条件时异常，请稍后再试。")
	NotFoundProcess            = New(33004, "获取业务流程时发生错误，请联系系统管理员")
	ProcessDefineError         = New(33005, "业务流程定义有误，请联系系统管理员")
	ProcessStepsError          = New(33006, "业务流程长度错误，请联系系统管理员")
	ProcessTransferBatchError  = New(33007, "批量接转失败，请重试")
	GetDataError               = New(33008, "获取数据发生错误，请联系系统管理员")
	CheckDataConflictError     = New(33009, "校验数据是否冲突时发生错误")
	ProcessingError            = New(33010, "业务正在处理中请稍后...")
	SystemError                = New(33011, "发生系统错误，请联系管理员")
)

// 账户相关   [20000 - 30000]

var (
	LoginExpire  = New(20000, "您的登录已失效！请重新登录")
	WrongPass    = New(20001, "登录失败 密码不正确")
	NoPermission = New(20003, "无权限访问")

	JWTGenError        = New(21006, "JWT token 生成失败")
	TokenInvalidError  = New(21007, "token无效")
	TokenMissing       = New(21008, "token 不存在")
	TokenMissingUid    = New(21009, "token 缺少 userId")
	TokenMissingLid    = New(21010, "token 缺少 leagueId")
	ContextMissingUid  = New(21011, "token 缺少 userId")
	DecryptParamsError = New(21012, "解密参数失败")
	IllegalLid         = New(21013, "非法的组织ID")
	NotAdmin           = New(21014, "该用户不是管理员")
)

// 小程序相关 [50000 - 60000]

var (
	QRCodeExpiredError = New(50001, "二维码过期")
	WxLoginFailed      = New(50002, "登录失败")
	WxLogOutFailed     = New(50003, "退出登录失败")
	WxAccountNotBound  = New(50004, "微信未绑定")
	WxUnBoundFailed    = New(50005, "微信解除绑定失败")
)
