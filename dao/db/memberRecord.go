package db
//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/gofrs/uuid"
//	"schoolserver/common/ecode"
//	"schoolserver/dao"
//	"schoolserver/dao/search"
//	"schoolserver/logger"
//	"schoolserver/models"
//	"time"
//	"xorm.io/xorm"
//)
//
//var RecordTableName = "t_memberRecord"
//var Process23TableName = "t_process23"
//var ProcessStep23TableName = "t_process_step23"
//
//const RecordPageNumPre int = 10
//
//type (
//	TMemberRecord struct {
//		// 表 t_biz_memberrecode
//		Id                        string `json:"id" xorm:"id"`
//		ProcessId                 string `json:"processId" xorm:"processId"`
//		ProcessType               int32  `json:"processType" xorm:"processType"`
//		UserId                    string `json:"userId" xorm:"userid"`
//		Name                      string `json:"name" xorm:"name"`
//		UserCode                  string `json:"userCode" xorm:"userCode"`
//		IdentityCardNo            string `json:"identityCardNo" xorm:"identityCardNo"`
//		JoinLeagueTime            string `json:"joinLeagueTime" xorm:"joinLeagueTime"` //团员入团时间
//		Mobile                    string `json:"mobile" xorm:"mobile"`                 //团员手机号
//		FileId                    string `json:"fileId" xorm:"fileId"`                 //用户档案id
//		LeagueId                  string `json:"leagueId" xorm:"leagueId"`
//		LeagueName                string `json:"leagueName" xorm:"leagueName"`
//		LeagueFullName            string `json:"leagueFullName" xorm:"leagueFullName"`
//		BeginLeagueId             string `json:"beginLeagueId" xorm:"beginLeagueId"`
//		BeginLeagueName           string `json:"beginLeagueName" xorm:"beginLeagueName"`
//		BeginLeagueFullName       string `json:"beginLeagueFullName" xorm:"beginLeagueFullName"`
//		BeginLeagueParentId       string `json:"beginLeagueParentId" xorm:"beginLeagueParentId"`
//		BeginLeagueParentName     string `json:"beginLeagueParentName" xorm:"beginLeagueParentName"`
//		BeginLeagueParentFullName string `json:"beginLeagueParentFullName" xorm:"beginLeagueParentFullName"`
//		AuditLeagueId             string `json:"auditLeagueId" xorm:"auditLeagueId"`
//		AuditLeagueName           string `json:"auditLeagueName" xorm:"auditLeagueName"`
//		AuditLeagueFullName       string `json:"auditLeagueFullName" xorm:"auditLeagueFullName"`
//		CreateTime                string `json:"createTime" xorm:"createTime"`
//		Status                    int32  `json:"status" xorm:"status"`
//		UpdateTime                string `json:"updateTime" xorm:"updateTime"`
//		Source                    string `json:"source" xorm:"source"`
//		Remark                    string `json:"remark" xorm:"-"`
//	}
//
//	TProcess23 struct {
//		Id             string `xorm:"id"`
//		ProcessName    string `xorm:"processName"`
//		ProcessType    int8   `xorm:"processType"`
//		CreaterType    int8   `xorm:"createrType"` // 1：组织创建；0：个人创建
//		UserId         string `xorm:"userId"`
//		Name           string `xorm:"name"`
//		LeagueId       string `xorm:"leagueId"`
//		LeagueName     string `xorm:"leagueName"`
//		LeagueFullName string `xorm:"leagueFullName"`
//		CreateTime     string `xorm:"createTime"`
//		Status         int8   `xorm:"status"`     // 流程状态,1进行中，,2已完成，3已终止
//		UpdateTime     string `xorm:"updateTime"` // 最后操作更新时间
//		Version        int8   `xorm:"version"`    // 当前走的流程步骤版本号
//		Source         string `xorm:"source"`
//
//		Steps []*TProcessStep23 `xorm:"-"`
//	}
//	TProcessStep23 struct {
//		Id             string                              `json:"id" xorm:"id"`
//		ProcessId      string                              `json:"processid" xorm:"processid"`
//		StepDefineId   string                              `json:"stepDefineId" xorm:"stepDefineId"`
//		StepName       string                              `json:"stepName" xorm:"stepName"`
//		OrderNo        int32                               `json:"orderNo" xorm:"orderNo"`
//		HandlerType    int8                                `json:"handlerType" xorm:"handlerType"` // 处理者类型,0：个人操作,1：组织操作
//		UserId        string                      `json:"userId" xorm:"userId"`
//		Name           string                              `json:"name" xorm:"-"`
//		LeagueId       string                              `json:"leagueId" xorm:"leagueId"`
//		LeagueName     string                              `json:"leagueName" xorm:"leagueName"`
//		LeagueFullName string                              `json:"leagueFullName" xorm:"leagueFullName"`
//		Status         int8                                `json:"status" xorm:"status"` // 步骤状态 1进行中，2已完成，3已终止
//		CreateTime     string                              `json:"createTime" xorm:"createTime"`
//		Content        *TProcessStep23Content              `json:"-" xorm:"-"`
//		Content2       *TProcessStep23Content2             `json:"-" xorm:"-"`
//		TContent       *TProcessStep23ContentSysTerminated `json:"-" xorm:"-"`
//		UpdateTime     string                              `json:"updateTime" xorm:"updateTime"`
//		Source         string                              `json:"source" xorm:"-"`
//		Operates       []*TProcessStepOperate23            `json:"-" xorm:"-"`
//		ContentStr     string                    `json:"content" xorm:"content"`
//	}
//
//	TProcessStep23Content struct {
//		UserId              string `json:"userId" xorm:"userid"`
//		Name                string `json:"name" xorm:"name"`
//		UserCode            string `json:"userCode" xorm:"userCode"`
//		IdentityCardNo      string `json:"identityCardNo" xorm:"identityCardNo"`
//		JoinLeagueTime      string `json:"joinLeagueTime" xorm:"joinLeagueTime"` //团员入团时间
//		LeagueId            string `json:"leagueId"`
//		LeagueName          string `json:"leagueName"`
//		LeagueFullName      string `json:"leagueFullName"`
//		BeginLeagueId       string `json:"beginLeagueId" xorm:"beginLeagueId"`
//		BeginLeagueName     string `json:"beginLeagueName" xorm:"beginLeagueName"`
//		BeginLeagueFullName string `json:"beginLeagueFullName" xorm:"beginLeagueFullName"`
//		AuditLeagueId       string `json:"auditLeagueId" xorm:"auditLeagueId"`
//		AuditLeagueName     string `json:"auditLeagueName" xorm:"auditLeagueName"`
//		AuditLeagueFullName string `json:"auditLeagueFullName" xorm:"auditLeagueFullName"`
//	}
//
//	TProcessStep23Content2 struct {
//		AuditLeagueId       string `json:"-"` // 审批团组织Id
//		AuditLeagueName     string `json:"auditLeagueName"`
//		AuditLeagueFullName string `json:"auditLeagueFullName"`
//		AuditResult         int    `json:"auditResult"`  // 审批结果 0不同意，1同意
//		AuditOpinion        string `json:"auditOpinion"` // 审批意见
//	}
//
//	TProcessStep23ContentSysTerminated struct {
//		StopType int32 `json:"stopType"`
//		//用于判断，业务被终止的原因//2 人被删除 ，3组织类别变更，4  组织迁移 ， 5 批量结转
//		ChangeType                int32  `json:"changeType"`
//		UserId                    string `json:"userId"`
//		Name                      string `json:"name"`
//		LeagueId                  string `json:"leagueId"`
//		LeagueName                string `json:"leagueName"`
//		LeagueFullName            string `json:"leagueFullName"`
//		BeginLeagueFullName       string `json:"beginLeagueFullName"`
//		BeginLeagueParentFullName string `json:"beginLeagueParentFullName"`
//		AuditLeagueFullName       string `json:"auditLeagueFullName"`
//		DeleteTime                string `json:"deleteTime"`
//	}
//
//	TProcessStepOperate23 struct {
//		Id             string         `xorm:"id"`
//		ProcessId      string         `xorm:"processId"`
//		StepId         string         `xorm:"stepId"`
//		HandlerType    int8           `xorm:"handlerType"`
//		UserId         string         `xorm:"userId"`
//		Name           string         `xorm:"name"`
//		LeagueId       string         `xorm:"leagueId"`
//		LeagueName     string         `xorm:"leagueName"`
//		LeagueFullName string         `xorm:"leagueFullName"`
//		CreateTime     string `xorm:"createTime"`
//	}
//)
//
////团委只创建业务，不创建流程
//func (tt *TMemberRecord) CreateProcessForMemberRecode2() error {
//	if engine == nil {
//		return errors.New("Database Handler is Nil")
//	}
//	// 开始记录到数据库（必须通过事务操作）
//	found, err := engine.Table(RecordTableName).Where("userId=? and status=?",tt.UserId,1).Get(&TMemberRecord{})
//	if found {
//		return errors.New("团员电子档案审批中,不能重复发起")
//	}
//	//
//	id, err := uuid.NewV4()
//	if err != nil {
//		logger.Error("%v", err)
//		return err
//	}
//	tt.Id = id.String()
//	tt.Status = dao.Process_STATUS_FINISH
//
//	var data_ = map[int]interface{}{
//		1: tt,
//	}
//	var table = map[int]string{
//		1: RecordTableName,
//	}
//	//tx, err := quDb.Begin()
//	tx := engine.NewSession()
//	if err != nil {
//		return fmt.Errorf("Begin Tx error:%v", err)
//	}
//	err = engineWrap(tx, func() error {
//		for k, vv := range data_ {
//			if err := tt.insertToForCreateProcess23(tx, table[k], vv); err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
////创建完整的 业务流程23
//func (tt *TMemberRecord) CreateProcessForMemberRecode() error {
//	if engine == nil {
//		return errors.New("Database Handler is Nil")
//	}
//	// 开始记录到数据库（必须通过事务操作）
//	found, err := engine.Table(RecordTableName).Where("userId=? and status=?",tt.UserId,1).Get(&TMemberRecord{})
//	if found {
//		return errors.New("团员电子档案审批中,不能重复发起")
//	}
//	//
//	id, err := uuid.NewV4()
//	if err != nil {
//		logger.Error("%v", err)
//		return err
//	}
//	tt.Id = id.String()
//	// 装配流程
//	//即 ，从 processes.yml中取出：电子档案流程
//	iprocessDefine, ok := dao.MapProcesses.Load(dao.MEMBER_RECODE)
//	if !ok {
//		return ecode.NotFoundProcess
//		//return errors.New(std.Not_Found_Process)
//	}
//
//	processDefine, ok := iprocessDefine.(*models.MemberRecord)
//	if !ok {
//		return ecode.ProcessDefineError
//		//return errors.New(std.Process_Define_Error)
//	}
//
//	processStepsDefine := processDefine.Steps
//	if len(processStepsDefine) <= 0 {
//		return ecode.ProcessStepsError
//		//return errors.New(std.Process_Steps_Error)
//	}
//
//	processId, err := uuid.NewV4()
//	if err != nil {
//		logger.Error("%v", err)
//		return err
//	}
//	process := TProcess23{
//		Id:             processId.String(),
//		ProcessType:    int8(dao.MEMBER_RECODE),
//		ProcessName:    processDefine.ProcessName,
//		CreaterType:    int8(processStepsDefine[0].HandlerType),
//		UserId:         tt.UserId,
//		Name:           tt.Name,
//		LeagueId:       tt.BeginLeagueId,
//		LeagueName:     tt.BeginLeagueName,
//		LeagueFullName: tt.BeginLeagueFullName,
//		CreateTime:     tt.CreateTime,
//		Status:         dao.Process_STATUS_ONGOING,
//		UpdateTime:     tt.CreateTime,
//		Version:        int8(processDefine.Version),
//		Source:         dao.Source_From,
//	}
//	// 冲突
//	id, err = uuid.NewV4()
//	if err != nil {
//		return err
//	}
//	tt.ProcessId = process.Id
//
//	ts4c := TProcessStep23Content{}
//
//	ts4c.UserId = tt.UserId
//	ts4c.Name = tt.Name
//	ts4c.UserCode = tt.UserCode
//	ts4c.IdentityCardNo = tt.IdentityCardNo
//	ts4c.JoinLeagueTime = tt.JoinLeagueTime
//	ts4c.LeagueId = tt.BeginLeagueId
//	ts4c.LeagueName = tt.BeginLeagueName
//	ts4c.LeagueFullName = tt.BeginLeagueFullName
//	ts4c.BeginLeagueId = tt.BeginLeagueId
//	ts4c.BeginLeagueName = tt.BeginLeagueName
//	ts4c.BeginLeagueFullName = tt.BeginLeagueFullName
//	ts4c.AuditLeagueId = tt.AuditLeagueId
//	ts4c.AuditLeagueName = tt.AuditLeagueName
//	ts4c.AuditLeagueFullName = tt.AuditLeagueFullName
//
//	// 设置第一步
//
//	ts4cStr, err := json.Marshal(ts4c)
//	if err != nil {
//		return fmt.Errorf("json Marshal Error:%v", err)
//	}
//
//	step1Id_, err := uuid.NewV4()
//	if err != nil {
//		logger.Error("%v", err)
//		return err
//	}
//	step1Id := step1Id_.String()
//	step1 := TProcessStep23{
//		Id:             step1Id,
//		ProcessId:      processId.String(),
//		StepDefineId:   processStepsDefine[0].StepId,
//		StepName:       processStepsDefine[0].StepName,
//		OrderNo:        processStepsDefine[0].OrderNo,
//		HandlerType:    int8(processStepsDefine[0].HandlerType),
//		UserId:         tt.UserId,
//		Name:           tt.Name,
//		LeagueId:       tt.BeginLeagueId,
//		LeagueName:     tt.BeginLeagueName,
//		LeagueFullName: tt.BeginLeagueFullName,
//		Status:         dao.Step_STATUS_FINISH, // 第一步标记为已完成
//		CreateTime:     tt.CreateTime,
//		ContentStr:     string(ts4cStr),
//		UpdateTime:     tt.CreateTime,
//		Source:         dao.Source_From,
//		Content:        &ts4c,
//	}
//	id_, err := uuid.NewV4()
//	if err != nil {
//		logger.Error("%v", err)
//		return err
//	}
//	step1Oper1 := TProcessStepOperate23{
//		Id:             id_.String(),
//		ProcessId:      processId.String(),
//		StepId:         step1.Id,
//		HandlerType:    step1.HandlerType,
//		UserId:         tt.UserId,
//		Name:           tt.Name,
//		LeagueId:       tt.BeginLeagueId,
//		LeagueName:     tt.BeginLeagueName,
//		LeagueFullName: tt.BeginLeagueFullName,
//		CreateTime:     tt.CreateTime,
//	}
//	sopers := make([]*TProcessStepOperate23, 0, 1)
//	sopers = append(sopers, &step1Oper1)
//	step1.Operates = sopers
//
//	steps := make([]*TProcessStep23, 0, 2)
//	steps = append(steps, &step1)
//	process.Steps = steps
//	// 下一步需要同时存下
//	nextStep, err := tt.getNextStepForTransferMemberRecord(&process, processStepsDefine)
//	if err != nil {
//		return fmt.Errorf("Get Next Step error:%v", err)
//	}
//	process.Steps = append(process.Steps, nextStep)
//	// 消息发送处理,只在第二步需要发送消息
//	stepIndex := make([]int32, 0, 2)
//	stepIndex = append(stepIndex, 1)
//	messages, messageReceivers, err := tt.getMessageInfo(&process, stepIndex)
//	if err != nil {
//		return err
//	}
//
//	stepOpers := make([]*TProcessStepOperate23, 0, len(step1.Operates)+len(nextStep.Operates)+1)
//	stepOpers = append(stepOpers, step1.Operates...)
//	if len(nextStep.Operates) > 0 {
//		stepOpers = append(stepOpers, nextStep.Operates...)
//	}
//	var data_ = map[int]interface{}{
//		1: tt,
//		2: &process,
//		3: []*TProcessStep23{&step1, nextStep},
//		6: messages,
//		7: messageReceivers,
//		8: stepOpers,
//	}
//	var table = map[int]string{
//		1: RecordTableName,
//		2: Process23TableName,
//		3: ProcessStep23TableName,
//		6: "t_message",
//		7: "t_message_receive",
//		8: "t_process_step_operate23",
//	}
//	//tx, err := quDb.Begin()
//	tx := engine.NewSession()
//
//	err = engineWrap(tx, func() error {
//		for k, vv := range data_ {
//			if err := tt.insertToForCreateProcess23(tx, table[k], vv); err != nil {
//				return err
//			}
//		}
//		// message & transaction添加es索引
//		if err := search.InitEsConn(); err != nil {
//			if err != nil {
//				return err
//			}
//		}
//		for _, messes := range messages {
//			for _, message := range messes {
//				err := message.AppendDoc()
//				if err != nil {
//					return err
//				}
//			}
//		}
//
//		for _, mrs := range messageReceivers {
//			for _, mr := range mrs {
//				err := mr.AppendDoc()
//				if err != nil {
//					return err
//				}
//			}
//		}
//		return nil
//	})
//
//	// ok，我们可以提交事务了
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// 得到团员档案审批的下一步
//func (tt *TMemberRecord) getNextStepForTransferMemberRecord(process *TProcess23,
//	stepDefines []models.MemberRecordStep) (*TProcessStep23, error) {
//
//	stepCounts := len(stepDefines)
//	curStepCounts := len(process.Steps)
//	if stepCounts < curStepCounts {
//		return nil, fmt.Errorf("Step is out Buond:%d < %d", stepCounts, curStepCounts)
//	}
//
//	if stepCounts == curStepCounts { // 走完流程
//		return nil, nil
//	}
//
//	if process == nil || len(process.Steps) == 0 { // 还没发起
//		return nil, errors.New("请先发起档案审批流程")
//	}
//	id_, err := uuid.NewV4()
//	if err != nil {
//		return nil, err
//	}
//	nextStepId := id_.String()
//
//	nextStep := TProcessStep23{}
//	nextStep.ProcessId = process.Id
//	nextStep.StepDefineId = stepDefines[curStepCounts].StepId
//	nextStep.StepName = stepDefines[curStepCounts].StepName
//	id_, err = uuid.NewV4()
//	if err != nil {
//		return nil, err
//	}
//	nextStep.Id = id_.String()
//	nextStep.OrderNo = stepDefines[curStepCounts].OrderNo
//	nextStep.HandlerType = int8(stepDefines[curStepCounts].HandlerType)
//	nextStep.Status = dao.Process_STATUS_ONGOING // 当前步状态是进行中
//	nextStep.CreateTime = time.Now().Format("2006-01-02 15:04:05")
//	nextStep.Content = nil
//	nextStep.UpdateTime = nextStep.CreateTime
//	nextStep.Source = dao.Source_From
//	nextStep.Id = nextStepId
//
//	stuff := process.Steps[0].Content // 从第一步取出发起申请时的一些信息
//	//nextStep.UserId = ""
//	nextStep.Name = ""
//	switch curStepCounts {
//	case 1: // 第二步
//		id_, err := uuid.NewV4()
//		if err != nil {
//			return nil, err
//		}
//		nextStep.LeagueId = stuff.LeagueId
//		nextStep.LeagueName = stuff.LeagueName
//		nextStep.LeagueFullName = stuff.LeagueFullName
//		tpso4Self := TProcessStepOperate23{
//			Id:             id_.String(),
//			ProcessId:      process.Id,
//			StepId:         nextStepId,
//			HandlerType:    int8(stepDefines[curStepCounts].HandlerType),
//			LeagueId:       stuff.AuditLeagueId,
//			LeagueName:     stuff.AuditLeagueName,
//			LeagueFullName: stuff.AuditLeagueFullName,
//			CreateTime:     tt.CreateTime,
//		}
//
//		tpso4s := make([]*TProcessStepOperate23, 0, 1)
//		tpso4s = append(tpso4s, &tpso4Self)
//		nextStep.Operates = tpso4s
//
//	}
//	return &nextStep, nil
//}
//
//// 得到某几步骤消息的发送者及接受者信息
//func (tt *TMemberRecord) getMessageInfo(process *TProcess23,
//	stepIndexes []int32) (map[int32][]TMessage, map[int32][]TMessageReceive, error) {
//
//	if len(stepIndexes) == 0 {
//		logger.Warn("stepIndexes lenth is 0")
//		return nil, nil, errors.New("stepIndexes length is 0")
//	}
//
//	if stepIndexes == nil {
//		return nil, nil, errors.New("Need step Index is NULL!")
//	}
//	if process == nil {
//		return nil, nil, errors.New("Process is NULL!")
//	}
//	steps := process.Steps
//	if len(steps) == 0 {
//		return nil, nil, errors.New("Not Found Needed Proces steps!")
//	}
//	stepsLen := len(steps)
//	if stepsLen < len(stepIndexes) {
//		return nil, nil, errors.New("Internal Logic Error!")
//	}
//
//	messages := make(map[int32][]TMessage, len(stepIndexes))
//	messagesReceiver := make(map[int32][]TMessageReceive, len(stepIndexes))
//	var msgReceivers []TMessageReceive
//	var messes []TMessage
//	for _, i := range stepIndexes {
//		if i > int32(stepsLen-1) {
//			logger.Error("step %d is greater than steps Max index value:%d", i, stepsLen-1)
//			continue
//		}
//		id_, err := uuid.NewV4()
//		if err != nil {
//			return nil, nil, err
//		}
//		curStep := steps[i]
//		ops := curStep.Operates
//		time_ := time.Now().Format("2006-01-02 15:04:05")
//		message := TMessage{
//			Id:          id_.String(),
//			MessageType: dao.MEMBER_RECODE,
//			CreateTime:  time_,
//			BusinessId:  tt.Id,
//			ProcessId:   process.Id,
//			StepId:      curStep.Id,
//		}
//		switch curStep.Status {
//		case dao.Step_STATUS_ONGOING:
//			messes = make([]TMessage, 0, 1)
//			message.MessageTitle = process.ProcessName
//			message.MessageContent = curStep.StepName
//
//			// 发送者为上一步的操作者
//			lastIndex := i - 1 // 只要是正在进行的状态 i一定大雨等于1
//			if lastIndex < 0 {
//				return nil, nil, errors.New("Found Error When Got Last Step,Index is less 0!")
//			}
//			lastStep := steps[lastIndex]
//			handlerType := lastStep.HandlerType
//			if handlerType == dao.PERSONAL {
//				message.SenderType = dao.PERSONAL
//				message.UserId = lastStep.UserId
//				message.Name = lastStep.Name
//			} else if handlerType == dao.LEAGUE {
//				message.SenderType = dao.LEAGUE
//				message.UserId = lastStep.UserId
//				message.Name = lastStep.Name
//				message.LeagueId = lastStep.LeagueId
//				message.LeagueName = lastStep.LeagueName
//				message.LeagueFullName = lastStep.LeagueFullName
//			} else {
//				// 内部应该出了问题
//				return nil, nil, errors.New("Error Handler Type is NOT Allow")
//			}
//			messes = append(messes, message)
//			ongoingReceiversCount := len(curStep.Operates)
//			msgReceivers = make([]TMessageReceive, 0, ongoingReceiversCount)
//			if ongoingReceiversCount > 0 {
//				for _, item := range ops {
//					id_, err := uuid.NewV4()
//					if err != nil {
//						return nil, nil, err
//					}
//					receiver := TMessageReceive{
//						Id:           id_.String(),
//						MessageId:    message.Id,
//						ReadStatus:   dao.MessageNoRead,
//						HandleStatus: dao.MessageNoHandle,
//						CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
//						HandleTime:   "0", //"0000-00-00 00:00:00",
//						ReadTime:     "0", // "0000-00-00 00:00:00",
//					}
//					handlerType := item.HandlerType
//					if handlerType == dao.PERSONAL {
//						receiver.ReceiverType = dao.PERSONAL
//						receiver.UserId = item.UserId
//						receiver.Name = item.Name
//					} else if handlerType == dao.LEAGUE {
//						receiver.ReceiverType = dao.LEAGUE
//						receiver.LeagueId = item.LeagueId
//						receiver.LeagueName = item.LeagueName
//						receiver.LeagueFullName = item.LeagueFullName
//						if process.LeagueId == dao.LEADER_LEAGUE_ID {
//							//如果流程的发起组织为团中央，就把创建流程的管理员id设置到消息的接受人中，
//							// 在消息查询时排除当前管理员给另外一个管理员审核
//							receiver.UserId = process.UserId
//							receiver.Name = process.Name
//						}
//					} else {
//						// 内部应该出了问题
//						return nil, nil, errors.New("Error Handler Type is NOT Allow.")
//					}
//					// 添加到接收者列表
//					msgReceivers = append(msgReceivers, receiver)
//				}
//			} else {
//				logger.Warn("Not Found message Receiver")
//			}
//
//			// case ongoing
//		case dao.Step_STATUS_FINISH:
//			messes = make([]TMessage, 0, 2)
//
//			message.MessageTitle = process.ProcessName
//			message.MessageContent = "您提交的团员电子档案已通过上级团组织审核。"
//			// 步骤完成时候的消息发送人,默认为当前步骤的操作者
//			handlerType := curStep.HandlerType
//			if handlerType == dao.PERSONAL {
//				message.UserId = curStep.UserId
//				message.Name = curStep.Name
//				message.SenderType = dao.PERSONAL
//			} else if handlerType == dao.LEAGUE {
//				message.UserId = curStep.UserId
//				message.Name = curStep.Name
//				message.SenderType = dao.LEAGUE
//				message.LeagueId = curStep.LeagueId
//				message.LeagueName = curStep.LeagueName
//				message.LeagueFullName = curStep.LeagueFullName
//			} else {
//				// 内部应该出了问题
//				return nil, nil, errors.New("Error Handler Type is NOT Allow..")
//			}
//			messes = append(messes, message)
//			// 消息接收者
//			var err error
//			msgReceivers, err = tt.getFinishStepMessageReceivers(process, i)
//			if err != nil {
//				return nil, nil, err
//			}
//			for idx, recever := range msgReceivers {
//				recever.HandleTime = "0" //"0000-00-00 00:00:00"
//				recever.ReadTime = "0"   // "0000-00-00 00:00:00"
//				recever.MessageId = message.Id
//				msgReceivers[idx] = recever
//			}
//
//			// case finish
//		case dao.Step_STATUS_REVOKED:
//			messes = make([]TMessage, 0, 1)
//
//			// case revoked(撤销)
//			message.MessageTitle = process.ProcessName
//			message.MessageContent = "撤销" + process.ProcessName + "申请"
//			createrType := process.CreaterType
//			// 流程撤销时的发送人,默认为当前流程的创建者
//			if createrType == 0 {
//				message.UserId = process.UserId
//				message.Name = process.Name
//				message.SenderType = dao.PERSONAL
//			} else if createrType == 1 {
//				message.UserId = process.UserId
//				message.Name = process.Name
//				message.LeagueId = process.LeagueId
//				message.LeagueName = process.LeagueName
//				message.LeagueFullName = process.LeagueFullName
//				message.SenderType = dao.LEAGUE
//			} else {
//				return nil, nil, fmt.Errorf("Error create Type:%d", createrType)
//			}
//			messes = append(messes, message)
//			var err error
//			msgReceivers, err = tt.getRevokedStepMessageReceivers(process, i)
//			if err != nil {
//				return nil, nil, err
//			}
//			for idx, recever := range msgReceivers {
//				recever.HandleTime = "0" //"0000-00-00 00:00:00"
//				recever.ReadTime = "0"   //"0000-00-00 00:00:00"
//				recever.MessageId = message.Id
//				msgReceivers[idx] = recever
//			}
//
//			// case revoked
//		case dao.Step_STATUS_TERMINATED:
//			messes = make([]TMessage, 0, 1)
//
//			message.MessageTitle = process.ProcessName
//			content2 := curStep.ContentStr
//			//审批被拒绝时发送的消息内容
//			msg := curStep.StepName + "不通过"
//			if content2 != "" {
//				tsc := TProcessStep23Content2{}
//				err := json.Unmarshal([]byte(content2), &tsc)
//				if err != nil {
//					return nil, nil, err
//				}
//				msg = fmt.Sprintf("您提交的团员电子档案未通过上级团组织审核，理由为：'%s'。请根据上级审核意见重新修改并上传新的文件。", tsc.AuditOpinion)
//			}
//			message.MessageContent = msg
//			// 步骤终结时候的消息发送人,默认为当前步骤的操作者
//			handlerType := curStep.HandlerType
//			if handlerType == dao.PERSONAL {
//				message.UserId = curStep.UserId
//				message.Name = curStep.Name
//				message.SenderType = dao.PERSONAL
//			} else if handlerType == dao.LEAGUE {
//				message.UserId = curStep.UserId
//				message.Name = curStep.Name
//				message.SenderType = dao.LEAGUE
//				message.LeagueId = curStep.LeagueId
//				message.LeagueName = curStep.LeagueName
//				message.LeagueFullName = curStep.LeagueFullName
//			} else {
//				// 内部应该出了问题
//				return nil, nil, errors.New("Error Handler Type is NOT Allow..")
//			}
//			messes = append(messes, message)
//			var err error
//			msgReceivers, err = tt.getTerminatedStepMessageReceivers(process, i)
//			if err != nil {
//				return nil, nil, err
//			}
//			for idx, recever := range msgReceivers {
//				recever.HandleTime = "0" //"0000-00-00 00:00:00"
//				recever.ReadTime = "0"   //"0000-00-00 00:00:00"
//				recever.MessageId = message.Id
//				msgReceivers[idx] = recever
//			}
//			// case terminated
//		case dao.Step_STATUS_SYS_TERMINATION:
//			messes = make([]TMessage, 0, 1)
//			message.SenderType = -1
//			tcontent := curStep.TContent
//			if tcontent == nil {
//				return nil, nil, errors.New("tcontent is NULL")
//			}
//			name := tcontent.Name
//			stopType := tcontent.StopType
//			message.MessageTitle = process.ProcessName
//			if stopType == 1 {
//				message.MessageContent = name + "已从" + tcontent.LeagueFullName + "被删除，被系统终止"
//			} else if stopType == 2 {
//				message.MessageContent = name + "已从" + tcontent.LeagueFullName + "被删除，被系统终止"
//			} else if stopType == 3 {
//				message.MessageContent = curStep.LeagueFullName + "或其上级的组织类别发生变更被系统终止"
//			} else if stopType == 4 {
//				message.MessageContent = curStep.LeagueFullName + "或其上级发生了组织迁移被系统终止"
//			} else if stopType == 5 {
//				message.MessageContent = tt.Name + "所在组织发生变化，被系统终止"
//			} else {
//				return nil, nil, nil
//			}
//			message.SenderType = dao.SYSTEM
//			messes = append(messes, message)
//			var err error
//			msgReceivers, err = tt.getSysTerminatedStepMessageReceivers(process, i)
//			if err != nil {
//				return nil, nil, err
//			}
//			for idx, recever := range msgReceivers {
//				recever.CreateTime = NowTimeStr()
//				recever.HandleTime = "0" // "0000-00-00 00:00:00"
//				recever.ReadTime = "0"   //"0000-00-00 00:00:00"
//				recever.MessageId = message.Id
//				msgReceivers[idx] = recever
//			}
//			//case sys-terminated
//		default:
//			return nil, nil, fmt.Errorf("Error Step Status:%d", curStep.Status)
//		} // switch
//		messages[i] = messes
//		messagesReceiver[i] = msgReceivers
//	} // for
//	return messages, messagesReceiver, nil
//}
//
//// 步骤不通过时的消息接受者
//func (tt *TMemberRecord) getTerminatedStepMessageReceivers(process *TProcess23, stepIndex int32) ([]TMessageReceive, error) {
//	return tt.getFinishStepMessageReceivers(process, stepIndex)
//}
//
//// 流程撤销时的消息接受者 默认为当前流程步骤的指定操作者
//func (tt *TMemberRecord) getRevokedStepMessageReceivers(process *TProcess23, stepIndex int32) ([]TMessageReceive, error) {
//	if process == nil || len(process.Steps) == 0 {
//		return nil, errors.New("Pls Give the process Params")
//	}
//	firstStep := process.Steps[0]
//	content := firstStep.Content
//	steps := process.Steps
//	if len(steps) < int(stepIndex) {
//		return nil, fmt.Errorf("Error stepIndex:%d", stepIndex)
//	}
//	curStep := steps[stepIndex]
//
//	var receivers []TMessageReceive
//	handleType := curStep.HandlerType
//
//	var receiver1 TMessageReceive
//	id_, err := uuid.NewV4()
//	if err != nil {
//
//		return nil, err
//	}
//	if handleType == dao.PERSONAL {
//		receiver1 = TMessageReceive{
//			Id:           id_.String(),
//			ReceiverType: dao.PERSONAL,
//			UserId:       curStep.UserId,
//			Name:         curStep.Name,
//			ReadStatus:   dao.MessageNoRead,
//			HandleStatus: dao.MessageNoNeed,
//		}
//	} else if handleType == dao.LEAGUE {
//		receiver1 = TMessageReceive{
//			Id:             id_.String(),
//			ReceiverType:   dao.LEAGUE,
//			UserId:         curStep.UserId,
//			Name:           curStep.Name,
//			LeagueId:       content.BeginLeagueId,
//			LeagueName:     content.BeginLeagueName,
//			LeagueFullName: content.BeginLeagueFullName,
//			ReadStatus:     dao.MessageNoRead,
//			HandleStatus:   dao.MessageNoNeed,
//		}
//	} else {
//		return nil, fmt.Errorf("Error Handler Type:%d", handleType)
//	}
//	receivers = make([]TMessageReceive, 0, 1)
//	receivers = append(receivers, receiver1)
//	return receivers, nil
//}
//
////消息完成时
//func (tt *TMemberRecord) getFinishStepMessageReceivers(process *TProcess23, stepIndex int32) ([]TMessageReceive, error) {
//	if process == nil || len(process.Steps) == 0 {
//		return nil, errors.New("Pls Give the process Params")
//	}
//	firstStep := process.Steps[0]
//	content := firstStep.Content
//	var receivers []TMessageReceive
//	switch stepIndex {
//	case 1: // 第二步
//		id_, err := uuid.NewV4()
//		if err != nil {
//			return nil, err
//		}
//		// 通知申请人审核结果
//		receiver1 := TMessageReceive{
//			Id:           id_.String(),
//			ReceiverType: dao.LEAGUE,
//			UserId:       content.UserId,
//			Name:         content.Name,
//
//			LeagueId:       tt.BeginLeagueId,
//			LeagueName:     tt.BeginLeagueName,
//			LeagueFullName: tt.BeginLeagueFullName,
//			ReadStatus:     dao.MessageNoRead,
//			HandleStatus:   dao.MessageNoNeed,
//			ReadTime:       "0",
//			HandleTime:     "0",
//			CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
//		}
//
//		receivers = make([]TMessageReceive, 0, 1)
//		receivers = append(receivers, receiver1, )
//		// case 1
//
//	default:
//		//return nil, fmt.Errorf("Error Step index:%d", stepIndex)
//	} // switch
//
//	return receivers, nil
//}
//
//// 系统终止流程
//func (tt *TMemberRecord) getSysTerminatedStepMessageReceivers(process *TProcess23, stepIndex int32) ([]TMessageReceive, error) {
//	if process == nil || len(process.Steps) == 0 {
//		return nil, errors.New("Pls Give the process Params")
//	}
//	firstStep := process.Steps[0]
//	content := firstStep.Content
//	curStep := process.Steps[stepIndex]
//	tcontent := curStep.TContent
//	if tcontent == nil {
//		return nil, errors.New("tContent is NULL")
//	}
//	var receivers []TMessageReceive
//	stopType := tcontent.StopType
//
//	switch stopType {
//
//	case 2, 3, 4, 5:
//		//2 人被删除 ，3组织类别变更，4  组织迁移 ， 5 批量结转
//		receivers = make([]TMessageReceive, 0, 2) // 最大量
//		id_1, err := uuid.NewV4()
//		if err != nil {
//			return nil, err
//		}
//		receiver1 := TMessageReceive{
//			Id:             id_1.String(),
//			ReceiverType:   dao.LEAGUE,
//			UserId:         content.UserId,
//			Name:           content.Name,
//			LeagueId:       tt.BeginLeagueId,
//			LeagueName:     tt.BeginLeagueName,
//			LeagueFullName: tt.BeginLeagueFullName,
//			ReadStatus:     dao.MessageNoRead,
//			HandleStatus:   dao.MessageNoNeed,
//			ReadTime:       "0",
//			HandleTime:     "0",
//			CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
//		}
//		receivers = append(receivers, receiver1)
//	default:
//
//	} // switch
//	return receivers, nil
//}
//
////向数据库插入数据
//func (tt *TMemberRecord) insertToForCreateProcess23(tx *xorm.Session, table string, d interface{}) error {
//
//	switch v := d.(type) {
//	case *TMemberRecord, *TProcess23:
//
//		//tMap := make(map[string]interface{})
//		//err := utils.Struct2MapByTagDb(tMap, v)
//		//if err != nil {
//		//	return err
//		//}
//		_, err:= tx.Table(table).Insert(d)
//		//_, err = insertTransfer.Exec()
//		return err
//
//	case []*TProcessStepOperate23:
//		//records := make([]goqu.Record, 0, len(v))
//		//for _, item := range v {
//		//	tMap := make(map[string]interface{})
//		//	err := utils.Struct2MapByTagDb(tMap, item)
//		//	if err != nil {
//		//		return err
//		//	}
//		//	records = append(records, goqu.Record(tMap))
//		//}
//		//if len(records) == 0 {
//		//	return nil
//		//}
//		//ins := tx.From(table).Insert(records)
//		//_, err := ins.Exec()
//		_, err  := tx.Table(table).Insert(d)
//		return err
//
//	case []*TProcessStep23:
//
//		//records := make([]goqu.Record, 0, len(v))
//		//for _, item := range v {
//		//	tMap := make(map[string]interface{})
//		//	err := utils.Struct2MapByTagDb(tMap, item)
//		//	if err != nil {
//		//		return err
//		//	}
//		//	records = append(records, goqu.Record(tMap))
//		//}
//		//ins := tx.From(table).Insert(records)
//		//_, err := ins.Exec()
//		_, err  := tx.Table(table).Insert(d)
//		return err
//
//	case map[int32][]TMessage:
//		err := AddChannelInsert(d)
//		if err != nil {
//			return err
//		}
//		return nil
//	case map[int32][]TMessageReceive:
//		err := AddChannelInsert(d)
//		if err != nil {
//			return err
//		}
//		return nil
//	default:
//		return fmt.Errorf("Unknown Type:%T", v)
//	}
//}
//
////获取团员电子档案的审批列表
////func (tt *TMemberRecord) GetTWAuditList(page, processType int, count *int, leagueId string) ([]TMemberRecord, error) {
////	sql_ := "select id,IFNULL(`processType`,23),userId,name,IFNULL(`userCode`,\"0\"),IFNULL(`leagueId`,\"0\")," +
////		"IFNULL(`leagueName`,\"0\"),IFNULL(`leagueFullName`,\"0\"),IFNULL(`createTime`,\"0\"),IFNULL(`joinLeagueTime`,\"0\")," +
////		"IFNULL(`beginLeagueId`,\"0\"),IFNULL(`beginLeagueName`,\"0\"),IFNULL(`beginLeagueFullName`,\"0\")," +
////		"status,IFNULL(`updateTime`,\"\"),tr.processId,IFNULL(`source`,\"web\") from t_biz_memberrecord tr" +
////		" where 1=1"
////	sqlCount := "select count(*) from t_biz_memberrecord tr " + " where 1=1"
////
////	if processType == 23 {
////		sql_ = sql_ + " and tr.processType=23"
////		sqlCount = sqlCount + " and tr.processType=23"
////	}
////	if len(leagueId) > 0 {
////		sql_ = sql_ + " and tr.auditLeagueId='" + leagueId + "' and tr.leagueId!='" + leagueId + "'"
////		sqlCount = sqlCount + " and tr.auditLeagueId='" + leagueId + "' and tr.leagueId!='" + leagueId + "'"
////	}
////	sql_ = sql_ + " order by tr.status, tr.createTime desc limit ?,?"
////	//sqlCount = sqlCount + " order by tr.status, tr.createTime desc"
////	logger.Debug("sql_:%s", sql_)
////	if page <= 0 {
////		page = 1
////	}
////	startRow := (page - 1) * RecordPageNumPre
////	rows, err := db_read.Query(sql_, startRow, RecordPageNumPre)
////	if err != nil {
////		logger.Error("Query error:%v", err)
////		return nil, err
////	}
////	err = db_read.QueryRow(sqlCount).Scan(count)
////	if err != nil {
////		logger.Error("QueryCount error:%v", err)
////		return nil, err
////	}
////	defer rows.Close()
////	d := make([]TMemberRecord, 0, RecordPageNumPre)
////	for rows.Next() {
////		ai := TMemberRecord{}
////		err := rows.Scan(&ai.Id, &ai.ProcessType, &ai.UserId, &ai.Name, &ai.UserCode, &ai.LeagueId,
////			&ai.LeagueName, &ai.LeagueFullName, &ai.CreateTime, &ai.JoinLeagueTime,
////			&ai.BeginLeagueId, &ai.BeginLeagueName, &ai.BeginLeagueFullName, &ai.Status,
////			&ai.UpdateTime, &ai.ProcessId, &ai.Source)
////		if err != nil {
////			logger.Error("Scan error:%v", err)
////			return nil, err
////		}
////		d = append(d, ai)
////	}
////	return d, nil
////}
//
////审批电子档案流程
//func (tt *TMemberRecord) OpersProcessForMmeberRecord(t *models.BusinessForMmeberRecordAuditReq, cb func() error) error {
//	if engine == nil {
//		return errors.New("Db IS NULL")
//	}
//	process_ := &TProcess23{Id: t.ProcessId}
//	rowCount, err := process_.GetProcessById()
//	if err != nil && rowCount < 0 {
//		return err
//	}
//	step := TProcessStep23{
//		ProcessId: t.ProcessId,
//	}
//	processSteps, err := step.GetStepsByProcessId()
//	if err != nil {
//		logger.Error("9789,GetStepsByProcessId error:%v", err)
//		return err
//	}
//	if len(processSteps[0].ContentStr) <= 0 {
//		return errors.New("Not Found Step Content for step 1")
//	}
//	stepContent := processSteps[0].ContentStr
//
//	ts1 := TProcessStep23Content{}
//	err = json.Unmarshal([]byte(stepContent), &ts1)
//	if err != nil {
//		logger.Error("9899,json.Unmarshal([]byte(stepContent),ts1) error:%v", err)
//		return err
//	}
//	processSteps[0].Content = &ts1
//	process_.Steps = make([]*TProcessStep23, 0, len(processSteps))
//	for idx, _ := range processSteps {
//		process_.Steps = append(process_.Steps, &processSteps[idx])
//	}
//
//	if process_.Status == dao.Process_STATUS_FINISH ||
//		process_.Status == dao.Process_STATUS_REVOKED ||
//		process_.Status == dao.Process_STATUS_SYS_TERMINATION ||
//		process_.Status == dao.Process_STATUS_TERMINATED {
//
//		logger.Warn("Process id:%s is Over and Status is:%d,Cant Oper Again", process_.Id, process_.Status)
//		return nil
//	}
//
//	wantStatus := t.StepStatus
//	switch wantStatus {
//	case dao.Step_STATUS_FINISH, dao.Step_STATUS_TERMINATED:
//		// 找出当前的step Id
//		var curStep *TProcessStep23
//		var stepIndex_ int
//		for stepIndex, step := range process_.Steps { // 从数据库提取是一定要 order by orderNo asc,因为slice的索引值就是第N步
//			if step.Id == t.ProcessStepId {
//				curStep = step
//				stepIndex_ = stepIndex
//				break
//			}
//		}
//		if curStep == nil {
//			return fmt.Errorf("Not Found Process Step by Id:%s", t.ProcessStepId)
//		}
//		if curStep.Status == dao.Step_STATUS_TERMINATED || curStep.Status == dao.Step_STATUS_FINISH {
//			logger.Warn("Id is %s Process Step is Had Over", t.ProcessStepId)
//			return nil
//		}
//
//		/*if t.LeagueId != tt.AuditLeagueId {
//			logger.Error("illegal非法的组织操作")
//			return errors.New("非法的组织操作")
//		}*/
//
//		operators, err := process_.GetOperatorsForStepByStepId(curStep.Id)
//		if err != nil {
//			return err
//		}
//		var curOperator TProcessStepOperate23
//		for _, op := range operators {
//			if (op.HandlerType == dao.PERSONAL && t.UserId == op.UserId) ||
//				(op.HandlerType == dao.LEAGUE && t.LeagueId == op.LeagueId) {
//				curOperator = op
//				break
//			}
//		}
//		curStep.Status = int8(wantStatus)
//		curStep.HandlerType = curOperator.HandlerType
//		curStep.UserId =  t.UserId
//		curStep.Name = t.UserName
//		curStep.LeagueId = t.LeagueId
//		curStep.LeagueName = t.LeagueName
//		curStep.LeagueFullName = t.LeagueFullName
//		curStep.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
//		curStep.ContentStr = t.ContentStr
//		curStep.Source = dao.Source_From
//
//		stepIndex__ := make([]int32, 0, 1)
//		stepIndex__ = append(stepIndex__, int32(stepIndex_))
//		message, messageReceiver, err := tt.getMessageInfo(process_, stepIndex__)
//		if err != nil {
//			logger.Error("Found Error When getMessageInfo: %v", err)
//			return err
//		}
//
//		//tMap := make(map[string]interface{})
//		//err = utils.Struct2MapByTagDb(tMap, curStep)
//		//if err != nil {
//		//	return err
//		//}
//		//tx, err := quDb.Begin()
//		tx:=engine.NewSession()
//
//
//		err = engineWrap(tx, func() error {
//			_, err := tx.Table(ProcessStep23TableName).Where("id=?",curStep.Id).Update(curStep)
//			//_, err := update.Exec()
//			if err != nil {
//				return err
//			}
//
//			if wantStatus == dao.Step_STATUS_TERMINATED {
//				process_.Status = dao.Process_STATUS_TERMINATED
//				process_.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
//				// 在事务中更新process状态、step状态、删除transaction_current的相应记录、发送消息；
//				// 同时把该流程绑定的未处理消息消息标识为无需处理状态，mysql与es中要一并更新修改
//				process_.ProcessType=23
//				_, err := tx.Table(Process23TableName).Where("id=?",process_.Id).Cols("status","updateTime","processType").Update(process_)
//
//				//_, err := update.Exec()
//				if err != nil {
//					return err
//				}
//				if cb != nil {
//					if err = cb(); err != nil {
//						return err
//					}
//				}
//				if err_ := tt.insertToForCreateProcess23(tx, "t_message", message); err_ != nil {
//					return err_
//				}
//				if err_ := tt.insertToForCreateProcess23(tx, "t_message_receive", messageReceiver); err_ != nil {
//					return err_
//				}
//
//				// es更新
//				if err := search.InitEsConn(); err != nil {
//					if err != nil {
//						return err
//					}
//				}
//				for _, messes := range message {
//					for _, mess := range messes {
//						err := mess.AppendDoc()
//						if err != nil {
//							return err
//						}
//					}
//				}
//
//				for _, mrs := range messageReceiver {
//					for _, mr := range mrs {
//						err := mr.AppendDoc()
//						if err != nil {
//							return err
//						}
//					}
//				}
//
//			} else if wantStatus == dao.Step_STATUS_FINISH {
//				// 此步已经完成， 得到下一步并标记正在进行
//				iprocessDefine, ok := dao.MapProcesses.Load(dao.MEMBER_RECODE)
//				if !ok {
//					return ecode.NotFoundProcess
//					//return errors.New(std.Not_Found_Process)
//				}
//				processDefine, ok := iprocessDefine.(*models.MemberRecord)
//				if !ok {
//					return ecode.ProcessDefineError
//					//return errors.New(std.Process_Define_Error)
//				}
//
//				processStepsDefine := processDefine.Steps
//				if len(processStepsDefine) <= 0 {
//					return ecode.ProcessStepsError
//					//return errors.New(std.Process_Steps_Error)
//				}
//				nextStep, err := tt.getNextStepForTransferMemberRecord(process_, processStepsDefine)
//				if err != nil {
//					return fmt.Errorf("Get Next Step error:%v", err)
//				}
//
//				if nextStep == nil { // 流程已经走完
//					process_.Status = dao.Process_STATUS_FINISH
//					process_.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
//					// 发送相关消息及删除冲突记录
//					stepIndex__ := make([]int32, 0, 1)
//					stepIndex__ = append(stepIndex__, int32(stepIndex_))
//					message, messageReceiver, err := tt.getMessageInfo(process_, stepIndex__)
//					if err != nil {
//						logger.Error("Found Error When getMessageInfo: %v", err)
//						return err
//					}
//					//tconfict := TTransactionCurrent{TransactionNo: 4, ProcessId: process_.Id}
//					process_.ProcessType=23
//					_, err = tx.Table(Process23TableName).Where("id=?",process_.Id).Cols("status","updateTime","processType").Update(process_)
//
//					//update := tx.Table(Process23TableName).Where("id=?",process_.Id)).Update(goqu.Record{"status": process_.Status,
//					//	"updateTime": process_.UpdateTime, "processType": 23})
//
//					//_, err = update.Exec()
//					if err != nil {
//						return err
//					}
//					if cb != nil {
//						if err = cb(); err != nil {
//							return err
//						}
//					}
//					if err_ := tt.insertToForCreateProcess23(tx, "t_message", message); err_ != nil {
//						return err_
//					}
//					if err_ := tt.insertToForCreateProcess23(tx, "t_message_receive", messageReceiver); err_ != nil {
//						return err_
//					}
//					// es更新
//					if err := search.InitEsConn(); err != nil {
//						if err != nil {
//							return err
//						}
//					}
//					for _, messes := range message {
//						for _, mess := range messes {
//							err := mess.AppendDoc()
//							if err != nil {
//								return err
//							}
//						}
//					}
//
//					for _, mrs := range messageReceiver {
//						for _, mr := range mrs {
//							err := mr.AppendDoc()
//							if err != nil {
//								return err
//							}
//						}
//					}
//
//				} else {
//					process_.Steps = append(process_.Steps, nextStep)
//					//tMap := make(map[string]interface{})
//					//err := utils.Struct2MapByTagDb(tMap, nextStep)
//					//if err != nil {
//					//	return err
//					//}
//
//					_, err := tx.Table(ProcessStep23TableName).Insert(nextStep)
//					//_, err = insert.Exec()
//					if err != nil {
//						return err
//					}
//					process_.Steps = append(process_.Steps, nextStep)
//					// 发消息
//					stepIndex__ := make([]int32, 0, 2)
//					stepIndex__ = append(stepIndex__, int32(stepIndex_), int32(stepIndex_+1))
//					message, messageReceiver, err := tt.getMessageInfo(process_, stepIndex__)
//					if err != nil {
//						return err
//					}
//
//					if err_ := tt.insertToForCreateProcess23(tx, "t_message", message); err_ != nil {
//						return err_
//					}
//					if err_ := tt.insertToForCreateProcess23(tx, "t_message_receive", messageReceiver); err_ != nil {
//						return err_
//					}
//					// es更新
//					if err := search.InitEsConn(); err != nil {
//						if err != nil {
//							return err
//						}
//					}
//					for _, messes := range message {
//						for _, mess := range messes {
//							err := mess.AppendDoc()
//							if err != nil {
//								return err
//							}
//						}
//					}
//
//					for _, mrs := range messageReceiver {
//						for _, mr := range mrs {
//							err := mr.AppendDoc()
//							if err != nil {
//								return err
//							}
//						}
//					}
//				}
//
//			} else {
//				return fmt.Errorf("Unkown wantStep Status:%d", wantStatus)
//			}
//
//			if err_ := tt.handleFinish(curStep.Id, process_, wantStatus, tx); err_ != nil {
//				return err_
//			}
//			return nil
//		})
//
//		if err != nil {
//			return err
//		}
//	case dao.Step_STATUS_SYS_TERMINATION: // 系统终止，比如用户已经被删除的情况下，系统需要值此状态
//		// 找出正在进行中的步骤
//		var stepOngoing *TProcessStep23
//		var stepIndex_ int
//		for stepIndex, step := range process_.Steps { // 从数据库提取是一定要 order by orderNo asc,因为slice的索引值就是第N步
//			if step.Status == dao.Step_STATUS_ONGOING {
//				stepOngoing = step
//				stepIndex_ = stepIndex
//				break
//			}
//		}
//		if stepOngoing == nil {
//			return errors.New("Not Found Goings Process Step")
//		}
//		stepOngoing.Status = dao.Step_STATUS_SYS_TERMINATION
//		stepOngoing.ProcessId = process_.Id
//		stepOngoing.TContent = &TProcessStep23ContentSysTerminated{}
//		stepOngoing.ContentStr = t.ContentStr
//		stepOngoing.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
//		err := json.Unmarshal([]byte(t.ContentStr), stepOngoing.TContent)
//		if err != nil {
//			logger.Error("json.Unmarshal(%s) error:%v", t.ContentStr, err)
//			return err
//		}
//		process_.Status = dao.Step_STATUS_SYS_TERMINATION // 标记此流程已经终止
//		// 发送相关消息及删除冲突记录
//		stepIndex__ := make([]int32, 0, 1)
//		stepIndex__ = append(stepIndex__, int32(stepIndex_))
//		message, messageReceiver, err := tt.getMessageInfo(process_, stepIndex__)
//		if err != nil {
//			logger.Error("Found Error When getMessageInfo: %v", err)
//			return err
//		}
//
//		// 在事务中更新process状态、step状态、删除transaction_current的相应记录、发送消息；
//		// 同时把该流程绑定的未处理消息消息标识为无需处理状态，mysql与es中要一并更新修改
//		//tx, err := quDb.Begin()
//		tx := engine.NewSession()
//		//if err != nil {
//		//	return err
//		//}
//		err = engineWrap(tx, func() error {
//			process_.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
//			process_.ProcessType=1
//			_, err := tx.Table(Process23TableName).Where("id=?",process_.Id).Cols("status","updateTime","processType").Update(process_)
//
//			//_, err := update.Exec()
//			if err != nil {
//				return err
//			}
//
//			//tMap := make(map[string]interface{})
//			//err = utils.Struct2MapByTagDb(tMap, stepOngoing)
//			//if err != nil {
//			//	return err
//			//}
//			_, err  = tx.Table(ProcessStep23TableName).Where("id=?",stepOngoing.Id).Update(stepOngoing)
//			//_, err = update.Exec()
//			if err != nil {
//				return err
//			}
//
//			if cb != nil {
//				if err = cb(); err != nil {
//					return err
//				}
//			}
//			//确认组织变更成功，才异步更新消息
//			if err_ := tt.insertToForCreateProcess23(tx, "t_message", message); err_ != nil {
//				return err_
//			}
//			if err_ := tt.insertToForCreateProcess23(tx, "t_message_receive", messageReceiver); err_ != nil {
//				return err_
//			}
//			// es更新
//			if err := search.InitEsConn(); err != nil {
//				if err != nil {
//					return err
//				}
//			}
//			for _, messes := range message {
//				for _, mess := range messes {
//					err := mess.AppendDoc()
//					if err != nil {
//						return err
//					}
//				}
//			}
//
//			for _, mrs := range messageReceiver {
//				for _, mr := range mrs {
//					err := mr.AppendDoc()
//					if err != nil {
//						return err
//					}
//				}
//			}
//
//			if err_ := tt.handleFinish(stepOngoing.Id, process_, wantStatus, tx); err_ != nil {
//				return err_
//			}
//			return nil
//		})
//
//		if err != nil {
//			return err
//		}
//
//	default:
//		logger.Error("Error Step Status:%d", wantStatus)
//		return fmt.Errorf("Error Step Want Status:%d", wantStatus)
//	}
//
//	return nil
//}
//
////通过userld 的指定流程
////func (tt *TMemberRecord) GetOne() (int, error) {
////	if db_ == nil {
////		return -1, errors.New("Database Handler is Nil")
////	}
////	sql_ := "select count(`id`) as totalRows from `t_biz_memberrecord` where `processType`=? and `userid`=? and `status`=?"
////	var totalRows int
////	err := db_read.QueryRow(sql_, tt.ProcessType, tt.UserId, tt.Status).Scan(&totalRows)
////	switch {
////	case err == sql.ErrNoRows:
////		return 0, nil
////	case err != nil:
////		return -1, err
////	default:
////		return totalRows, nil
////
////	}
////}
//
////通过id获取档案流程
////func (tt *TMemberRecord) GetMemberRecordById() (int, error) {
////	if db_ == nil {
////		return -1, errors.New("Database Handler is Nil")
////	}
////
////	if len(tt.Id) <= 0 {
////		return -1, errors.New("Error transfer Id")
////	}
////	ok, err := quDbRead.From(RecordTableName).Where(
////		goqu.I("id").Eq(tt.Id)).ScanStruct(tt)
////	switch {
////	case !ok:
////		return 0, nil
////	case err != nil:
////		return -1, err
////	default:
////		return 1, nil
////
////	}
////}
//
////通过 userId获取档案流程
////func (tt *TMemberRecord) GetUnfinishedByUserId() (int, error) {
////	if db_ == nil {
////		return -1, errors.New("Database Handler is Nil")
////	}
////
////	if len(tt.UserId) <= 0 {
////		return -1, errors.New("error_user_id")
////	}
////	ok, err := quDbRead.From(RecordTableName).Where(
////		goqu.I("userId").Eq(tt.UserId),
////		//goqu.I("status").Eq(1),
////	).Order(goqu.I("createTime").Desc()).Limit(1).ScanStruct(tt)
////	switch {
////	case !ok:
////		return 0, nil
////	case err != nil:
////		return -1, err
////	default:
////		return 1, nil
////
////	}
////}
//
//// 获取团员的档案审批记录
////func (tt *TMemberRecord) GetMemberProcess() (int, error) {
////	if db_ == nil {
////		return -1, errors.New("database Handler is Nil")
////	}
////
////	if len(tt.UserId) <= 0 {
////		return -1, errors.New("error_user_id")
////	}
////	ok, err := quDbRead.From(RecordTableName).Where(
////		goqu.I("userId").Eq(tt.UserId),
////	).Order(goqu.I("createTime").Desc()).Limit(1).ScanStruct(tt)
////	switch {
////	case !ok:
////		return 0, nil
////	case err != nil:
////		return -1, err
////	default:
////		return 1, nil
////
////	}
////}
//
////通过id获取 Process
//func (tp *TProcess23) GetProcessById() (int, error) {
//	if engine == nil {
//		return -1, errors.New("Database Handler is Nil")
//	}
//
//	if len(tp.Id) <= 0 {
//		return -1, errors.New("Error transfer Id")
//	}
//
//	ok, err := engine.Table(Process23TableName).Where("id=?",tp.Id).Get(tp)
//	switch {
//	case !ok:
//		return 0, nil
//	case err != nil:
//		return -1, err
//	default:
//		return 1, nil
//
//	}
//}
////
////通过id获取进行到的步骤
//func (ts *TProcessStep23) GetStepsByProcessId() ([]TProcessStep23, error) {
//	if len(ts.ProcessId) <= 0 {
//		return nil, errors.New("错误的参数：流程Id")
//	}
//	d := make([]TProcessStep23, 0, 3)
//	err := engine.Table(ProcessStep23TableName).Where("processId=?",ts.ProcessId).Asc("orderNo").Find(&d)
//		//Order(goqu.I("orderNo").Asc()).ScanStructs(&d)
//	if err != nil {
//		return nil, err
//	}
//
//	return d, nil
//}
//
////通过stepId获取Operators信息
//func (tp *TProcess23) GetOperatorsForStepByStepId(stepId string) ([]TProcessStepOperate23, error) {
//	if engine == nil {
//		return nil, errors.New("Db is NULL")
//	}
//
//	if len(stepId) <= 0 {
//		return nil, errors.New("Error step Id")
//	}
//	d := make([]TProcessStepOperate23, 0, 4)
//	_,err := engine.Table("t_process_step_operate23").Where("stepId=?",stepId).Get(&d)
//	if err != nil {
//		return nil, err
//	}
//
//	return d, nil
//}
//
////业务全部走完
//func (tt *TMemberRecord) handleFinish(stepId string, process_ *TProcess23, businessStatus int, tx *xorm.Session) error {
//	var messageIds []string
//	messageIds, err := SelectMessage(tt.Id, process_.Id, 23)
//	if err != nil {
//		return err
//	}
//
//	//将要修改的消息 逐个放入Channel中处理
//	tm := TMessageReceiverCh{
//		HandleStatus: dao.MessageHandled,
//		HandleTime:   process_.UpdateTime,
//		ReadStatus:   dao.MessageRead,
//		ReadTime:     process_.UpdateTime,
//	}
//	for _, messId := range messageIds {
//		tm.MessageId = messId
//		err = tm.AddChannelUpdate() //异步向Channel 中存放数据更新消息状态
//		if err != nil {
//			return err
//		}
//	}
//
//	// es更新
//	mr := TMessageReceive{}
//	for _, messId := range messageIds {
//		mr.MessageId = messId
//		err := mr.UpdateDocForHandleStatusMuti(&tm)
//		if err != nil {
//			return err
//		}
//	}
//	// 更新TRegister表状态
//	if err_ := tt.UpdateProcessStatus(tx, businessStatus); err_ != nil {
//		return err_
//	}
//	return nil
//}
//
////更新TMemberRecord表状态
//func (tt *TMemberRecord) UpdateProcessStatus(tx *xorm.Session, status int) error {
//
//	_, err  := tx.Table(RecordTableName).Where("id=?",tt.Id).Update(&TMemberRecord{ Status:int32(status),UpdateTime:time.Now().Format("2006-01-02 15:04:05")})
//	//_, err := update.Exec()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//////组织迁移时查询 ，所有 被迁移组织及其下级 正在进行中的流程。参数是被迁移的组织
////func SelectMemberRecordForOrgmove(moveLeagueId string, d *[]TMemberRecord) error {
////	if db_ == nil {
////		return errors.New("Database Handler is Nil")
////	}
////
////	if len(moveLeagueId) <= 0 {
////		return errors.New("Error transfer Id")
////	}
////	err := quDbRead.From(RecordTableName).Where(
////		goqu.Ex{
////			"status": 1,
////		},
////		goqu.Or(
////			goqu.I("leagueId").Eq(moveLeagueId),
////			goqu.I("beginLeagueId").Eq(moveLeagueId),
////			goqu.I("beginLeagueParentId").Eq(moveLeagueId),
////		), ).ScanStructs(d)
////	if err != nil {
////		return err
////	}
////	return nil
////
////}
////
//////批量结转时 查询，接转前所在组织 正在进行的业务,参数是：接转前的组织id
////func SelectMemberRecordForTransfer(transferLeagueId string, userId string, d *[]TMemberRecord) error {
////	if db_ == nil {
////		return errors.New("Database Handler is Nil")
////	}
////
////	if len(transferLeagueId) <= 0 {
////		return errors.New("Error MemberRecord leagueId")
////	}
////	if len(userId) <= 0 {
////		return errors.New("Error MemberRecord userId")
////	}
////	err := quDbRead.From(RecordTableName).Where(
////		goqu.Ex{
////			"status":   1,
////			"userId":   userId,
////			"leagueId": transferLeagueId,
////		}).ScanStructs(d)
////	if err != nil {
////		return err
////	}
////	return nil
////
////}
////
/////*
////	组织列别变更：团委变成团总支或团支部时，查询 由该组织审核的业务。
////	参数是：被变更组织类别的 组织id
////*/
////func SelectMemberRecordForLeagutypeChange(changeLeagueId string, d *[]TMemberRecord) error {
////	if db_ == nil {
////		return errors.New("Database Handler is Nil")
////	}
////	if len(changeLeagueId) <= 0 {
////		return errors.New("Error MemberRecord leagueId")
////	}
////	err := quDbRead.From(RecordTableName).Where(
////		goqu.Ex{
////			"status":        1,
////			"auditLeagueId": changeLeagueId,
////		}).ScanStructs(d)
////	if err != nil {
////		return err
////	}
////	return nil
////
////}
//
///*
//	组织列别变更：团总支或团支部变成团委时，查询 由该组织发起的业务。
//	参数是：被变更组织类别的 组织id
//*/
////func SelectMemberRecordByBeginLeagueId(BeginLeagueId string, d *[]TMemberRecord) error {
////	if db_ == nil {
////		return errors.New("Database Handler is Nil")
////	}
////
////	if len(BeginLeagueId) <= 0 {
////		return errors.New("Error MemberRecord leagueId")
////	}
////	err := quDbRead.From(RecordTableName).Where(
////		goqu.Ex{
////			"status":        1,
////			"beginLeagueId": BeginLeagueId,
////		}, ).ScanStructs(d)
////	if err != nil {
////		return err
////	}
////	return nil
////
////}
//
///*
//	组织列别变更：团总支或团支部变成团委时，查询 由该组织下级发起的业务。
//	参数是：被变更组织类别的 组织id
//*/
////func SelectMemberRecordByBeginLeagueParentId(BeginLeagueParentId string, d *[]TMemberRecord) error {
////	if db_ == nil {
////		return errors.New("Database Handler is Nil")
////	}
////
////	if len(BeginLeagueParentId) <= 0 {
////		return errors.New("Error MemberRecord leagueId")
////	}
////	err := quDbRead.From(RecordTableName).Where(
////		goqu.Ex{
////			"status":              1,
////			"beginLeagueParentId": BeginLeagueParentId,
////		}, ).ScanStructs(d)
////	if err != nil {
////		return err
////	}
////	return nil
////
////}
//
////团员被删除时 查询，该团员 正在进行的  档案审批业务
//func SelectMemberRecordForDelMember(leagueId string, userId string, d *[]TMemberRecord) error {
//	if engine == nil {
//		return errors.New("Database Handler is Nil")
//	}
//
//	if len(leagueId) <= 0 {
//		return errors.New("Error MemberRecord leagueId")
//	}
//	if len(userId) <= 0 {
//		return errors.New("Error MemberRecord userId")
//	}
//	err := engine.Table(RecordTableName).Where("status=? and userId =? and leagueId=? ", 1,userId,leagueId).Find(d)
//	if err != nil {
//		return err
//	}
//	return nil
//
//}
//
//////通过团员身份证号 获取团员你正在进行中的流程
////func SelectMemberRecordByIdcord(identityCardNo string, d *[]TMemberRecord) error {
////	if db_ == nil {
////		return errors.New("Database Handler is Nil")
////	}
////
////	if len(identityCardNo) <= 0 {
////		return errors.New("Error MemberRecord idcord")
////	}
////	err := quDbRead.From(RecordTableName).Where(
////		goqu.Ex{
////			"status":   1,
////			"identityCardNo": identityCardNo,
////
////		}).ScanStructs(d)
////	if err != nil {
////		return err
////	}
////	return nil
////
////}
