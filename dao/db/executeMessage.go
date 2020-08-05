package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"schoolserver/dao/search"
	"schoolserver/logger"
	"time"
)

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/olivere/elastic"
//	"schoolserver/common/utils"
//
//	"schoolserver/dao/search"
//	"schoolserver/http/handles"
//	"schoolserver/logger"
//	"time"
//)
//
var (
	//channel
	TmessageInsertChan         chan TMessage
	TmessageReceiverInsertChan chan TMessageReceive
	TmessageReceiverUpdateChan chan TMessageReceiverCh
	//数据库表名
	DbMessage        = "t_message"
	DbMessageReceive = "t_message_receive"
)
//
/*
 异步向 channel 新增消息
 用做向 t_message_receive 及 t_message 表的新增
*/
func AddChannelInsert(d interface{}) error {
	logger.Debug("异步向 channel 新增消息,用做向 t_message_receive 及 t_message 表的新增")

	switch v := d.(type) {
	case map[int32][]TMessage:

		for _, val := range v {
			for _, item := range val {
				tm := item
				go func() {
					idleDuration := 10 * time.Minute
					idleDelay := time.NewTimer(idleDuration)
					defer idleDelay.Stop()
					for {
						select {
						case TmessageInsertChan <- tm:
							logger.Debug("放入了TmessageInsertChan中")
							return
						case <-idleDelay.C:

							logger.Error("Asynchronization Insert Message:异步向 TmessageInsertChan 新增数据时等待超时")
							return
						}
					}
				}()
			}
		}
		return nil

	case map[int32][]TMessageReceive:
		for _, val := range v {
			for _, item := range val {
				tm := item
				go func() {
					idleDuration := 10 * time.Minute
					idleDelay := time.NewTimer(idleDuration)
					defer idleDelay.Stop()
					for {
						select {

						case TmessageReceiverInsertChan <- tm:
							logger.Debug("放入了TmessageInsertChan中")
							return
						case <-idleDelay.C:
							logger.Error("Asynchronization Insert Message:异步向 TmessageReceiverInsertChan 新增数据时等待超时")
							return
						}
					}
				}()
			}
		}
		return nil
	case []TMessage:

		for _, item := range v {
			tm := item
			go func() {
				idleDuration := 10 * time.Minute
				idleDelay := time.NewTimer(idleDuration)
				defer idleDelay.Stop()
				for {
					select {

					case TmessageInsertChan <- tm:
						logger.Debug("实际放入的item：")
						logger.Debug(tm)
						logger.Debug("放入了TmessageReceiverInsertChan中")
						return
					case <-idleDelay.C:
						logger.Error("Asynchronization Insert Message:异步向 TmessageInsertChan 新增数据时等待超时")
						return
					}
				}
			}()
		}
		return nil
	case []TMessageReceive:
		for _, item := range v {
			tm := item
			go func() {
				idleDuration := 10 * time.Minute
				idleDelay := time.NewTimer(idleDuration)
				defer idleDelay.Stop()
				for {
					idleDelay.Reset(idleDuration)
					select {

					case TmessageReceiverInsertChan <- tm:
						logger.Debug("实际放入的item：")
						logger.Debug(tm)
						logger.Debug("放入了TmessageReceiverInsertChan中")
						return
					case <-idleDelay.C:
						logger.Error("Asynchronization Insert Message:异步向 TmessageReceiverInsertChan 新增数据时等待超时")
						return
					}
				}
			}()
		}
		return nil
	default:
		return fmt.Errorf("Unknown Type:%T", v)
	}
}

/*
 异步向 channel 新增消息
 用做向 t_message_receive  表的修改数据
*/
func (tm TMessageReceiverCh) AddChannelUpdate() error {
	//修改时，通过 id 将其他非空项 全部修改
	logger.Debug("向 channel 新增消息用做向 t_message_receive  表的修改数据")
	go func() {
		idleDuration := 10 * time.Minute
		idleDelay := time.NewTimer(idleDuration)
		defer idleDelay.Stop()
		for {
			select {
			case TmessageReceiverUpdateChan <- tm:
				logger.Debug("放入了TmessageReceiverUpdateChan中")
				logger.Debug(tm)
				return
			case <-idleDelay.C:
				logger.Error("Asynchronization Insert Message:异步向 TmessageReceiverUpdateChan 新增数据时等待超时")
				return
			}
		}
	}()
	return nil
}
//
///*
//  异步向 t_message 表新增消息
//*/
//func (tm *TMessage) PInsertMessage(goroutine_cnt chan int) error {
//	defer handles.deleteGoroutine(goroutine_cnt)
//	logger.Debug("异步向 t_message 表新增消息")
//
//	tMap := make(map[string]interface{})
//	err := utils.Struct2MapByTagDb(tMap, tm)
//	if err != nil {
//		logger.Error("Asynchronization Insert Message:", err)
//		return err
//	}
//	tx:= engineMessage.NewSession()
//	defer tx.Close()
//	//if err != nil {
//	//	logger.Error("Asynchronization Insert Message:新增Message失败，将重试三次", err)
//	//	MessageFailureRetry(tm)
//	//	return err
//	//}
//	_, err = tx.Table(DbMessage).Insert(tMap)
//
//	if err != nil {
//		logger.Error("Asynchronization Insert Message:", err)
//		return tx.Rollback()
//	}
//	logger.Debug("异步向 t_message 表新增消息，新增成功")
//	return tx.Commit()
//}
//
///*
//  异步向t_message_receive 表新增消息
//*/
//func (tm *TMessageReceive) PInsertMessageReceive(goroutine_cnt chan int) error {
//	defer deleteGoroutine(goroutine_cnt)
//	logger.Debug("异步向 t_message_receive  表新增消息")
//
//	//tMap := make(map[string]interface{})
//	//err := utils.Struct2MapByTagDb(tMap, tm)
//	//if err != nil {
//	//	logger.Error("Asynchronization Insert Message:", err)
//	//	return err
//	//}
//
//	//if err != nil {
//	//	logger.Error("Asynchronization Insert Message:新增MessageReceiver失败，将重试三次", err)
//	//	MessageFailureRetry(tm)
//	//	return err
//	//}
//	tx := engineMessage.NewSession()
//	defer tx.Close()
//	_,err := tx.Table(DbMessageReceive).Insert(tm)
//	//logger.Debug(insertTransfer.Sql)
//	//_, err = insertTransfer.Exec()
//	if err != nil {
//		logger.Error("Asynchronization Insert Message:", err)
//		return tx.Rollback()
//	}
//	logger.Debug("异步向 t_message_receive  表新增消息，新增成功")
//	return tx.Commit()
//}
//
///*
//  异步更新t_message_receive 表 消息
//*/
//func (tm *TMessageReceiverCh) PUpdateMessageReceive(goroutine_cnt chan int) error {
//	defer handles.deleteGoroutine(goroutine_cnt)
//	logger.Debug("异步向 t_message_receive  表修改消息")
//	//将数据msg结构体 转换为map ，过滤空字段。通过id 将非空数据全部更新
//	tMap := make(map[string]interface{})
//	nMap := make(map[string]interface{}) //最终修改的字段集合
//	err := utils.Struct2MapByTagDb(tMap, tm)
//	if err != nil {
//		logger.Error("Asynchronization Insert Message:", err)
//		return err
//	}
//	for key, value := range tMap {
//		//要修改的信息中过滤 id 和 messageId
//		if !utils.IsEmptys(value) && "id" != key && "messageId" != key {
//			nMap[key] = value
//		}
//	}
//	tx := engineMessage.NewSession()
//	defer tx.Close()
//	if !utils.IsEmptys(tm.Id) { //如果id 不为空，则通过id 修改其他字段的非空数据
//		_,err = tx.Table(DbMessageReceive).Where("id=?",tm.Id).And("readStatus=?",0).Update(nMap)
//		//logger.Debug(update.Sql)
//		//_, err = update.Exec()
//	} else { //如果id 为空，则通过messageId 修改其他字段的非空数据. 2019.7.8  只允许修改操作状态为未审批的，消息提醒无需修改
//		_,err = tx.Table(DbMessageReceive).Where("messageId=?",tm.MessageId).And("handleStatus=?",0).Update(nMap)
//		//update := tx.From(DbMessageReceive).Where(goqu.I("messageId").Eq(tm.MessageId), goqu.I("handleStatus").Eq(0)).Update(nMap)
//		//logger.Debug(update.Sql)
//		//_, err = update.Exec()
//	}
//	if err != nil {
//		logger.Error("Asynchronization Insert Message:", err)
//		return tx.Rollback()
//	}
//	logger.Debug("异步向 t_message_receive  表修改消息，修改成功")
//	return tx.Commit()
//}
//
//通过业务流程 获取的messageID集合
func SelectMessage(businessId string, processId string, messageType int) ([]string, error) {
	var messageIds []string

	//先从es中获取当前业务所对应的消息
	boolQuery := elastic.NewBoolQuery().Must(elastic.NewTermQuery("processId", processId), elastic.NewTermQuery("businessId", businessId),
		elastic.NewTermQuery("messageType", messageType))
	res, err := search.EsClient.Search(search.GetCoreIdx0()).Type(search.GetCoreType0()).From(0).Size(100).Query(boolQuery).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if res.Hits.TotalHits > 0 {
		r := make([]TMessage, 0, res.Hits.TotalHits)
		for _, hit := range res.Hits.Hits {
			t := TMessage{}
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				return nil, err
			}
			r = append(r, t)
		}
		var messIds []string
		for _, item := range r {
			messIds = append(messIds, item.Id)
		}
		logger.Debug("es得到的", messIds)
		messageIds = messIds
	}
	//如果es中获取不到，则从数据库来获取
	if len(messageIds) == 0 {
		err := engineMessage.Table(DbMessage).Select("id").Where("businessId=?",businessId).And("processId=?",processId).And("messageType=?",messageType).Find(&messageIds)
		if err != nil {
			return nil, err
		}
		logger.Debug("数据库得到的", messageIds)
	}
	return messageIds, nil
}



///*
//  获取 消息的channel空闲数 ，用于界面跟踪
//*/
////func GetMessageChannel() utils.MessageChannel {
////
////	m1 := utils.ChannelDetails{
////		InUseNum:  len(TmessageInsertChan),
////		UnusedNum: cap(TmessageInsertChan) - len(TmessageInsertChan),
////		MaxNum:    cap(TmessageInsertChan),
////	}
////	m2 := utils.ChannelDetails{
////		InUseNum:  len(TmessageReceiverInsertChan),
////		UnusedNum: cap(TmessageReceiverInsertChan) - len(TmessageReceiverInsertChan),
////		MaxNum:    cap(TmessageReceiverInsertChan),
////	}
////	m3 := utils.ChannelDetails{
////		InUseNum:  len(TmessageReceiverUpdateChan),
////		UnusedNum: cap(TmessageReceiverUpdateChan) - len(TmessageReceiverUpdateChan),
////		MaxNum:    cap(TmessageReceiverUpdateChan),
////	}
////	mc := utils.MessageChannel{
////		TmessageInsertChanNum:  m1,
////		TmessageRInsertChanNum: m2,
////		TmessageRUpdateChanNum: m3,
////	}
////
////	return mc
////}
////
////func deleteGoroutine(goroutine_cnt chan int) {
////	<-goroutine_cnt
////}
//
///*
//  异步插入消息失败时时，重试三次（每次间隔1秒）
//*/
////func MessageFailureRetry(d interface{}) {
////	logger.Debug("进入异步插入消息失败重试！")
////	switch v := d.(type) {
////	case *TMessage:
////		tMap := make(map[string]interface{})
////		err := utils.Struct2MapByTagDb(tMap, v)
////		if err != nil {
////			logger.Error("message Failure to retry:", err)
////			return
////		}
////
////		for i := 0; i < 3; i++ {
////			//每次重试时先等待一秒
////			time.Sleep(1 * time.Second)
////			tx:= engineMessage.NewSession()
////			defer tx.Close()
////			//if err != nil {
////			//	logger.Error("message Failure to retry:", err)
////			//	continue
////			//}
////			_, err= tx.Table(DbMessage).Insert(tMap)
////			//_, err = insertTransfer.Exec()
////			if err != nil {
////				logger.Error("message Failure to retry:", err)
////				tx.Rollback()
////				continue
////			}
////			logger.Debug("异步失败重试向 t_message 表新增消息，新增成功")
////			tx.Commit()
////			return
////		}
////
////		return
////	case *TMessageReceive:
////		tMap := make(map[string]interface{})
////		err := utils.Struct2MapByTagDb(tMap, v)
////		if err != nil {
////			logger.Error("message Failure to retry:", err)
////			return
////		}
////		for i := 0; i < 3; i++ {
////			time.Sleep(1 * time.Second)
////			tx:= engineMessage.NewSession()
////			defer tx.Close()
////			//if err != nil {
////			//	logger.Error("message Failure to retry:", err)
////			//	continue
////			//}
////			xorm.Record
////			_, err= tx.Table(DbMessageReceive).Insert(xorm.Record(tMap))
////			//logger.Debug(insertTransfer.Sql)
////			//_, err = insertTransfer.Exec()
////			if err != nil {
////				logger.Error("message Failure to retry:", err)
////				tx.Rollback()
////				continue
////			}
////			logger.Debug("异步失败重试向 t_message_receive  表新增消息，新增成功")
////			tx.Commit()
////			return
////		}
////		return
////	case *TMessageReceiverCh:
////
////		//将数据msg结构体 转换为map ，过滤空字段。通过id 将非空数据全部更新
////		tMap := make(map[string]interface{})
////		nMap := make(map[string]interface{}) //最终修改的字段集合
////		err := utils.Struct2MapByTagDb(tMap, v)
////		if err != nil {
////			logger.Error("message Failure to retry:", err)
////			return
////		}
////		for key, value := range tMap {
////			//要修改的信息中过滤 id 和 messageId
////			if !utils.IsEmptys(value) && "id" != key && "messageId" != key {
////				nMap[key] = value
////			}
////		}
////		for i := 0; i < 3; i++ {
////			time.Sleep(1 * time.Second)
////			tx, err := quDb2.Begin()
////			if err != nil {
////				logger.Error("message Failure to retry:", err)
////				continue
////			}
////			if !utils.IsEmptys(v.Id) { //如果id 不为空，则通过id 修改其他字段的非空数据
////				update := tx.From(DbMessageReceive).Where(goqu.I("id").Eq(v.Id)).Update(nMap)
////				logger.Debug(update.Sql)
////				_, err = update.Exec()
////			} else { //如果id 为空，则通过messageId 修改其他字段的非空数据
////				update := tx.From(DbMessageReceive).Where(goqu.I("messageId").Eq(v.MessageId)).Update(nMap)
////				logger.Debug(update.Sql)
////				_, err = update.Exec()
////			}
////			if err != nil {
////				logger.Error("message Failure to retry:", err)
////				tx.Rollback()
////				continue
////			}
////			logger.Debug("异步失败重试向 t_message_receive  表修改消息，修改成功")
////			tx.Commit()
////			return
////		}
////
////		return
////	default:
////		logger.Error("message Failure to retry:类型匹配错误！")
////		return
////	}
////}
