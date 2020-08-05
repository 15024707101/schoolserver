package db

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/olivere/elastic"
	"schoolserver/dao/search"
	"schoolserver/logger"
	"reflect"
)

const PageNumPreForLoad int = 10000

var fieldName = "aipid"

type (
	AppearMessageReceiver struct {
		Id             string `json:"id" xorm:"id"`
		MessageId      string `json:"messageId" xorm:"messageId"`
		ReceiverType   int32  `json:"receiverType" xorm:"receiverType"` // 接收者类型,0：个人接收,1：组织接收
		UserId         string `json:"-" xorm:"userId"`
		Name           string `json:"name" xorm:"name"`
		LeagueId       string `json:"-" xorm:"leagueId"`
		LeagueName     string `json:"leagueName" xorm:"leagueName"`
		LeagueFullName string `json:"leagueFullName" xorm:"leagueFullName"`
		ReadStatus     int32  `json:"readStatus" xorm:"readStatus"`     // 是否已经阅读过,0：未读,1：已读
		HandleStatus   int32  `json:"handleStatus" xorm:"handleStatus"` // 是否处理该条消息业务,0：未处理,2：不需要处理，3：已处理
		ReadTime       string `json:"readTime" xorm:"-"`                // layout: 2006-01-02 15:04:05
		HandleTime     string `json:"handleTime" xorm:"-"`
		CreateTime     string `json:"createTime" xorm:"createTime"` // layout: 2006-01-02 15:04:05
		//Aipid          int `json:"-" xorm:"-"`
	}

	TMessage struct {
		Id             string `json:"id" xorm:"id"`
		SenderType     int32  `json:"senderType" xorm:"senderType"` // 发送者类型 0：个人发送 1：组织发送 -1 系统发送
		UserId         string `json:"userId" xorm:"UserId"`
		Name           string `json:"name" xorm:"name"`
		LeagueId       string `json:"leagueId" xorm:"leagueId"`
		LeagueName     string `json:"leagueName" xorm:"leagueName"`
		LeagueFullName string `json:"leagueFullName" xorm:"leagueFullName"`
		MessageType    int32  `json:"messageType" xorm:"messageType"`
		MessageTitle   string `json:"messageTitle" xorm:"messageTitle"`
		MessageContent string `json:"messageContent" xorm:"messageContent"`
		CreateTime     string `json:"createTime" xorm:"createTime"` // layout: 2006-01-02 15:04:05
		BusinessId     string `json:"businessId" xorm:"businessId"`
		ProcessId      string `json:"processId" xorm:"processId"`
		StepId         string `json:"stepId" xorm:"stepId"`
		//Aipid          int `json:"-" xorm:"-"`
	}

	TMessageReceive struct {
		Id             string `xorm:"id" json:"id" xorm:"id"`
		MessageId      string `xorm:"messageId" json:"messageId" xorm:"messageId"`
		ReceiverType   int32  `xorm:"receiverType" json:"receiverType" xorm:"receiverType"` // 接收者类型,0：个人接收,1：组织接收
		UserId         string `xorm:"userId" json:"userId" xorm:"userId"`
		Name           string `xorm:"name" json:"name" xorm:"name"`
		LeagueId       string `xorm:"leagueId" json:"leagueId" xorm:"leagueId"`
		LeagueName     string `xorm:"leagueName" json:"leagueName" xorm:"leagueName"`
		LeagueFullName string `xorm:"leagueFullName" json:"leagueFullName" xorm:"leagueFullName"`
		ReadStatus     int32  `xorm:"readStatus" json:"readStatus" xorm:"readStatus"`     // 是否已经阅读过,0：未读,1：已读
		HandleStatus   int32  `xorm:"handleStatus" json:"handleStatus" xorm:"handleStatus"` // 是否处理该条消息业务,0：未处理,2：不需要处理，3：已处理
		ReadTime       string `xorm:"readTime" json:"readTime" xorm:"-"`                // layout: 2006-01-02 15:04:05
		HandleTime     string `xorm:"handleTime" json:"handleTime" xorm:"-"`
		CreateTime     string `xorm:"createTime" json:"createTime" xorm:"createTime"` // layout: 2006-01-02 15:04:05
		//Aipid          int `json:"-" xorm:"-"`
	}

	TMessageReceiverCh struct {
		//用于向 Channel 中存放的结构体
		Id           string `json:"id" xorm:"id"`
		MessageId    string `json:"messageId" xorm:"messageId"`
		ReadStatus   int32  `json:"readStatus" xorm:"readStatus"`     // 是否已经阅读过,0：未读,1：已读
		HandleStatus int32  `json:"handleStatus" xorm:"handleStatus"` // 是否处理该条消息业务,0：未处理,2：不需要处理，3：已处理
		ReadTime     string `json:"readTime" xorm:"readTime"`         // layout: 2006-01-02 15:04:05
		HandleTime   string `json:"handleTime" xorm:"handleTime"`
	}

	OperationCenterParam struct {
		// 操作中心传递给函数 GetListForOperation 的参数
		Page, PageSize           int
		HandleStatus, ReadStatus int
		MessageType              int
		StartTime, EndTime       string
		UserId, LeagueId         string
	}

	OperationCenterResp struct {
		AppearMessageReceiver
		SenderUserId       string `json:"-"`
		SenderName         string `json:"senderName"`
		SendLeagueId       string `json:"senderLeagueId"`
		SendLeagueName     string `json:"senderLeagueName"`
		SendLeagueFullName string `json:"senderLeagueFullName"`
		SenderType         int    `json:"senderType"`
		MessageType        int    `json:"messageType"`
		MessageTitle       string `json:"messageTitle"`
		MessageContent     string `json:"messageContent"`
		BusinessId         string `json:"businessId"`
		ProcessId          string `json:"-"`
		StepId             string `json:"-"`
	}
)











// 向es中添加message
func (m *TMessage) AppendDoc() error {
	indexResp, err := search.EsClient.Index().Index("fulltextsearchcore0").
		Type("TMessage").BodyJson(m).Id(m.Id).Do(context.Background())
	if err != nil {
		return err
	}
	logger.Debug("Append Doc for Message,indexResp:%+v", indexResp)
	if indexResp.Shards.Successful >= 1 {
		return nil
	}
	return errors.New("Shard count is less 1")
}
func (m *TMessage) UpdateDoc(doc *TMessage) error {
	updateResp, err := search.EsClient.Update().Index("fulltextsearchcore0").
		Type("TMessage").Id(m.Id).Doc(doc).DetectNoop(false).Do(context.Background())
	if err != nil {
		return err
	}
	logger.Debug("update for message and resp:%+v", updateResp)
	if updateResp.Shards.Successful >= 1 {
		return nil
	}
	return errors.New("Update Failed.	")
}

func (mr *TMessageReceive) AppendDoc() error {
	indexResp, err := search.EsClient.Index().Index("fulltextsearchcore1").
		Type("TMessageReceiver").BodyJson(mr).Id(mr.Id).Do(context.Background())
	if err != nil {
		return err
	}
	logger.Debug("Append Doc for MessageReceiver,indexResp:%+v", indexResp)
	if indexResp.Shards.Successful >= 1 {
		return nil
	}
	return errors.New("Shard count is less 1")
}

// es更新消息状态
func (mr *TMessageReceive) UpdateDocForHandleStatusMuti(tc *TMessageReceiverCh) error {
	//2019.7.8 只允许修改操作状态为未审批的，消息提醒无需修改
	var tq elastic.Query
	var boolQuery *elastic.BoolQuery

	tq = elastic.NewTermsQuery("messageId", mr.MessageId)
	boolQuery = elastic.NewBoolQuery().Must(tq, elastic.NewTermQuery("handleStatus", 0))

	sc := elastic.NewScript("ctx._source.handleStatus=params.handleStatus;ctx._source.handleTime=params.handleTime;ctx._source.readStatus=params.readStatus;ctx._source.readTime=params.readTime")
	sc.Param("handleStatus", tc.HandleStatus)
	sc.Param("handleTime", tc.HandleTime)
	sc.Param("readStatus", tc.ReadStatus)
	sc.Param("readTime", tc.ReadTime)
	updateResp, err := search.EsClient.UpdateByQuery().Index("fulltextsearchcore1").
		Type("TMessageReceiver").Query(boolQuery).Script(sc).Do(context.Background())
	if err != nil {
		return err
	}
	logger.Debug("update for message receiver and resp:%+v", updateResp)
	if updateResp.Updated >= 1 || updateResp.Total == 0 {
		return nil
	}
	return errors.New("Update Failed.	")
}

//func (mr *TMessageReceive) UpdateDbForReadStatus(rs int, updateTime string) error {
//	//将消息放入Channel 中异步处理
//	tm := TMessageReceiverCh{
//		Id:         mr.Id,
//		ReadStatus: int32(rs),
//		ReadTime:   updateTime,
//	}
//	err := tm.AddChannelUpdate()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (mr *TMessageReceive) UpdateDocForReadStatus(hs int, updateTime string) error {
//
//	var boolQuery *elastic.BoolQuery
//	tq1 := elastic.NewTermQuery("id", mr.Id)
//	tq2 := elastic.NewTermQuery("readStatus", 0)
//	boolQuery = elastic.NewBoolQuery().Must(tq1, tq2)
//	sc := elastic.NewScript("ctx._source.readStatus=params.readStatus;ctx._source.readTime=params.readTime")
//	sc.Param("readStatus", hs)
//	sc.Param("readTime", updateTime)
//	updateResp, err := search.EsClient.UpdateByQuery().Index("fulltextsearchcore1").
//		Type("TMessageReceive").Query(boolQuery).Script(sc).Do(context.Background())
//	if err != nil {
//		return err
//	}
//	logger.Debug("update for message receiver and resp:%+v", updateResp)
//	if updateResp.Updated >= 1 || updateResp.Total == 0 {
//		return nil
//	}
//	return errors.New("Update Failed.	")
//}

//func (tt *TMessage) checkTxErr(tx *goqu.TxDatabase, result sql.Result, err error) error {
//	if err == nil {
//		return nil
//	}
//	if rErr := tx.Rollback(); rErr != nil {
//		return rErr
//	}
//	if result != nil {
//		var rowsCount int64
//		if rowsCount, err = result.RowsAffected(); err != nil {
//			//if rErr := tx.Rollback(); rErr != nil {
//			//	return rErr
//			//}
//			return err
//		}
//		logger.Debug("insert transfer and affect row:%d", rowsCount)
//	}
//	return err
//}

//func (tt *TMessageReceive) checkTxErr(tx *goqu.TxDatabase, result sql.Result, err error) error {
//	if err == nil {
//		return nil
//	}
//	if rErr := tx.Rollback(); rErr != nil {
//		return rErr
//	}
//	if result != nil {
//		var rowsCount int64
//		if rowsCount, err = result.RowsAffected(); err != nil {
//			//if rErr := tx.Rollback(); rErr != nil {
//			//	return rErr
//			//}
//			return err
//		}
//		logger.Debug("insert transfer and affect row:%d", rowsCount)
//	}
//	return err
//}

//func (mes *TMessage) InsertTo(table string, d interface{}) error {
//	err := AddChannelInsert(d)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//职务变更，删除团员，批量接转 发送消息
func (ms *TMessage) InsertMesAndMr(mess []TMessage, mrs []TMessageReceive) error {
	//err := ms.InsertTo("t_message", mess)
	//if err != nil {
	//	return err
	//}
	//err = ms.InsertTo("t_message_receive", mrs)
	//if err != nil {
	//	return err
	//}
	err := AddChannelInsert(mess)
	if err != nil {
		return err
	}
	err = AddChannelInsert(mrs)
	if err != nil {
		return err
	}
	for _, message := range mess {
		err := message.AppendDoc()
		if err != nil {
			logger.Error("append message To es error:%v", err)
			return err
		}
	}
	for _, mr := range mrs {
		err := mr.AppendDoc()
		if err != nil {
			logger.Error("append message recevicer errror:%v", err)
			return err
		}
	}
	return nil
}

// 操作中心消息列表顶部的统计
func GetAggByHandleStatus(param *OperationCenterParam) (map[int]int64, error) {
	boolQuery := elastic.NewBoolQuery()
	if len(param.LeagueId) > 0 { // 组织的消息（管理员）
		terms := []elastic.Query{
			elastic.NewTermQuery("receiverType", 1),
			elastic.NewTermQuery("leagueId", param.LeagueId),
		}

		bq := elastic.NewBoolQuery().Must(terms...)

		boolQuery = boolQuery.Should(bq).MinimumNumberShouldMatch(1)
	} else {
		boolQuery = elastic.NewBoolQuery().Must(elastic.NewTermQuery("receiverType", 0),
			elastic.NewTermQuery("userId", param.UserId))
	}
	// 聚合
	aggs := elastic.NewTermsAggregation().Field("handleStatus").OrderByCount(false).Size(30)
	res, err := search.EsClient.Search(search.GetCoreIdx1()).Type(search.GetCoreType1()).Query(boolQuery).Aggregation("messStatis", aggs).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	agg, found := res.Aggregations.Terms("messStatis")
	if !found {
		return nil, errors.New("Not Found Aggregations data.")
	}

	ret := make(map[int]int64, len(agg.Buckets))
	for _, bucket := range agg.Buckets {

		k, ok := bucket.KeyNumber.Int64()
		if ok != nil {
			return nil, errors.New("Type Asset Error")
		}
		ret[int(k)] = bucket.DocCount
	}

	return ret, nil
}

func GetAggByReadStatus(param *OperationCenterParam) (map[int]int64, error) {
	boolQuery := elastic.NewBoolQuery()
	if len(param.LeagueId) > 0 { // 组织的消息（管理员）
		terms := []elastic.Query{
			elastic.NewTermQuery("receiverType", 1),
			elastic.NewTermQuery("leagueId", param.LeagueId),
		}

		bq := elastic.NewBoolQuery().Must(terms...)

		boolQuery = boolQuery.Should(bq).MinimumNumberShouldMatch(1)
	} else {
		boolQuery = elastic.NewBoolQuery().Must(elastic.NewTermQuery("receiverType", 0),
			elastic.NewTermQuery("userId", param.UserId))
	}
	// 聚合
	aggs := elastic.NewTermsAggregation().Field("readStatus").OrderByCount(false).Size(30)
	res, err := search.EsClient.Search(search.GetCoreIdx1()).Type(search.GetCoreType1()).Query(boolQuery).Aggregation("messReadStatis", aggs).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	agg, found := res.Aggregations.Terms("messReadStatis")
	if !found {
		return nil, errors.New("Not Found Aggregations data.")
	}

	ret := make(map[int]int64, len(agg.Buckets))
	for _, bucket := range agg.Buckets {
		k, ok := bucket.KeyNumber.Int64()
		if ok != nil {
			return nil, errors.New("Type Asset Error")
		}
		ret[int(k)] = bucket.DocCount
	}

	return ret, nil
}

// 顶部菜单操作中心列表数据
func GetListForOperation(param *OperationCenterParam, totalRow *int) ([]OperationCenterResp, error) {
	if param.MessageType > -1 {
		return getListForMessage(param, totalRow)
	}
	boolQuery := elastic.NewBoolQuery()
	if len(param.LeagueId) > 0 { // 组织的消息（管理员）
		terms := []elastic.Query{
			elastic.NewTermQuery("receiverType", 1),
			elastic.NewTermQuery("leagueId", param.LeagueId),
		}

		bq := elastic.NewBoolQuery().Must(terms...)

		boolQuery = boolQuery.Should(bq).MinimumNumberShouldMatch(1)
	} else {
		boolQuery = elastic.NewBoolQuery().Must(elastic.NewTermQuery("receiverType", 0),
			elastic.NewTermQuery("userId", param.UserId))
	}
	if param.ReadStatus > -1 {
		boolQuery = boolQuery.Filter(elastic.NewTermQuery("readStatus", param.ReadStatus))
	}

	if param.HandleStatus > -1 {
		boolQuery = boolQuery.Filter(elastic.NewTermQuery("handleStatus", param.HandleStatus))
	}

	if len(param.StartTime) > 0 {
		boolQuery = boolQuery.Filter(elastic.NewRangeQuery("createTime").Gte(param.StartTime))
	}

	if len(param.EndTime) > 0 {
		boolQuery = boolQuery.Filter(elastic.NewRangeQuery("createTime").Lte(param.EndTime))
	}

	fromIdx := (param.Page - 1) * param.PageSize
	//临时解决方案  等待优化
	res, err := search.EsClient.Search(search.GetCoreIdx1()).Type(search.GetCoreType1()).From(fromIdx).Size(10).Query(boolQuery).Sort("handleStatus", true).
		Sort("readStatus", true).Sort("createTime", false).Do(context.Background())
	if err != nil {
		return nil, err
	}
	*totalRow = int(res.Hits.TotalHits)
	if res.Hits.TotalHits > 0 {
		r := make([]AppearMessageReceiver, 0, res.Hits.TotalHits)
		for _, hit := range res.Hits.Hits {
			t := AppearMessageReceiver{}
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				return nil, err
			}
			r = append(r, t)
		}
		// 父子关系的查询需要在建立mapping时建立父子映射，且在同一个index里，
		// 目前不支持，只能检索后再过滤
		messIds := make([]interface{}, 0, res.Hits.TotalHits)
		for _, item := range r {
			messIds = append(messIds, item.MessageId)
		}
		logger.Debug(messIds)
		var re *elastic.SearchResult
		//查询所有是因为teams查询出来的结果是按照t_message本身数据的顺序匹配messIds
		if param.MessageType > -1 {
			re, err = search.EsClient.Search(search.GetCoreIdx0()).Type(search.GetCoreType0()).
				Query(
					elastic.NewBoolQuery().Must(
						elastic.NewTermsQuery("id", messIds...),
					).Must(elastic.NewTermQuery("messageType", param.MessageType)),
				).From(0).Size(len(messIds)).Do(context.Background())
		} else {
			re, err = search.EsClient.Search(search.GetCoreIdx0()).Type(search.GetCoreType0()).
				Query(
					//elastic.NewBoolQuery().Must(
					elastic.NewTermsQuery("id", messIds...),
					//),
				).From(0).Size(len(messIds)).Do(context.Background())
		}
		if err != nil {
			return nil, err
		}
		if re.Hits.TotalHits > 0 {
			rr := make([]OperationCenterResp, 0, re.Hits.TotalHits)
		OUTLOOP:
			for _, mr := range r { // 保证排序结果的一致
				for _, hit := range re.Hits.Hits {
					t := TMessage{}
					err := json.Unmarshal(*hit.Source, &t)
					if err != nil {
						return nil, err
					}

					if mr.MessageId == t.Id {
						oc := OperationCenterResp{
							mr,
							t.UserId,
							t.Name,
							t.LeagueId,
							t.LeagueName,
							t.LeagueFullName,
							int(t.SenderType),
							int(t.MessageType),
							t.MessageTitle,
							t.MessageContent,
							t.BusinessId,
							t.ProcessId,
							t.StepId,
						}
						rr = append(rr, oc)
						continue OUTLOOP
					}
				}
			}
			//*totalRow = len(rr)
			//if param.Page*param.PageSize > *totalRow {
			//	return rr[fromIdx:], nil
			//}
			//return rr[fromIdx : param.Page*param.PageSize], nil
			return rr, nil
		} else {
			return nil, nil
		}

	}

	return nil, nil
}

func getListForMessage(param *OperationCenterParam, totalRow *int) ([]OperationCenterResp, error) {

	boolQuery := elastic.NewBoolQuery()
	if len(param.LeagueId) > 0 { // 组织的消息（管理员）
		terms := []elastic.Query{
			elastic.NewTermQuery("receiverType", 1),
			elastic.NewTermQuery("leagueId", param.LeagueId),
		}

		bq := elastic.NewBoolQuery().Must(terms...)

		boolQuery = boolQuery.Should(bq).MinimumNumberShouldMatch(1)
	} else {
		boolQuery = elastic.NewBoolQuery().Must(elastic.NewTermQuery("receiverType", 0),
			elastic.NewTermQuery("userId", param.UserId))
	}
	if param.ReadStatus > -1 {
		boolQuery = boolQuery.Filter(elastic.NewTermQuery("readStatus", param.ReadStatus))
	}

	if param.HandleStatus > -1 {
		boolQuery = boolQuery.Filter(elastic.NewTermQuery("handleStatus", param.HandleStatus))
	}

	if len(param.StartTime) > 0 {
		boolQuery = boolQuery.Filter(elastic.NewRangeQuery("createTime").Gte(param.StartTime))
	}

	if len(param.EndTime) > 0 {
		boolQuery = boolQuery.Filter(elastic.NewRangeQuery("createTime").Lte(param.EndTime))
	}

	fromIdx := (param.Page - 1) * param.PageSize
	//临时解决方案  等待优化
	res, err := search.EsClient.Search(search.GetCoreIdx1()).Type(search.GetCoreType1()).From(0).Size(10000).Query(boolQuery).Sort("handleStatus", true).
		Sort("readStatus", true).Sort("createTime", false).Do(context.Background())
	if err != nil {
		return nil, err
	}
	//*totalRow = int(res.Hits.TotalHits)
	if res.Hits.TotalHits > 0 {
		r := make([]AppearMessageReceiver, 0, res.Hits.TotalHits)
		for _, hit := range res.Hits.Hits {
			t := AppearMessageReceiver{}
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				return nil, err
			}
			r = append(r, t)
		}
		// 父子关系的查询需要在建立mapping时建立父子映射，且在同一个index里，
		// 目前不支持，只能检索后再过滤
		messIds := make([]interface{}, 0, res.Hits.TotalHits)
		for _, item := range r {
			messIds = append(messIds, item.MessageId)
		}
		logger.Debug(messIds)
		var re *elastic.SearchResult
		//查询所有是因为teams查询出来的结果是按照t_message本身数据的顺序匹配messIds
		if param.MessageType > -1 {
			re, err = search.EsClient.Search(search.GetCoreIdx0()).Type(search.GetCoreType0()).
				Query(
					elastic.NewBoolQuery().Must(
						elastic.NewTermsQuery("id", messIds...),
					).Must(elastic.NewTermQuery("messageType", param.MessageType)),
				).From(0).Size(len(messIds)).Do(context.Background())
		} else {
			re, err = search.EsClient.Search(search.GetCoreIdx0()).Type(search.GetCoreType0()).
				Query(
					//elastic.NewBoolQuery().Must(
					elastic.NewTermsQuery("id", messIds...),
					//),
				).From(0).Size(len(messIds)).Do(context.Background())
		}
		if err != nil {
			return nil, err
		}
		if re.Hits.TotalHits > 0 {
			rr := make([]OperationCenterResp, 0, re.Hits.TotalHits)
			messages := make(map[string]TMessage, 0)
			for _, hit := range re.Hits.Hits {
				t := TMessage{}
				err := json.Unmarshal(*hit.Source, &t)
				if err != nil {
					return nil, err
				}
				messages[t.Id] = t
			}
			//OUTLOOP:
			for _, mr := range r { // 保证排序结果的一致
				if t, ok := messages[mr.MessageId]; ok {
					oc := OperationCenterResp{
						mr,
						t.UserId,
						t.Name,
						t.LeagueId,
						t.LeagueName,
						t.LeagueFullName,
						int(t.SenderType),
						int(t.MessageType),
						t.MessageTitle,
						t.MessageContent,
						t.BusinessId,
						t.ProcessId,
						t.StepId,
					}
					rr = append(rr, oc)
				}
			}
			*totalRow = len(rr)

			if param.Page*param.PageSize > *totalRow {
				list := make([]OperationCenterResp, *totalRow-fromIdx)
				copy(list, rr[fromIdx:])
				return list, nil
			} else {
				list := make([]OperationCenterResp, param.PageSize)
				copy(list, rr[fromIdx:param.Page*param.PageSize])
				return list, nil
			}

			//return rr,nil
		} else {
			return nil, nil
		}
	}
	return nil, nil
}

func GetMessageListByBusinessIdAndLeagueId(businessId, leagueId string) ([]TMessage, error) {
	if len(businessId) == 0 || len(leagueId) == 0 {
		return nil, errors.New("缺少必要参数")
	}
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(elastic.NewTermQuery("businessId", businessId), elastic.NewTermQuery("leagueId", leagueId))
	searchResult, err := search.EsClient.Search().Index(search.GetCoreIdx0()).Type(search.GetCoreType0()).Query(boolQuery).Do(context.Background())
	if err != nil {
		return nil, err
	}
	messageList := make([]TMessage, 0)
	var message TMessage
	for _, item := range searchResult.Each(reflect.TypeOf(message)) {
		t := item.(TMessage)
		messageList = append(messageList, t)
	}
	return messageList, nil
}

//func HandleMessage(messageId string, handlestatus int) error {
//	boolQuery := elastic.NewBoolQuery()
//	boolQuery = boolQuery.Must(elastic.NewTermQuery("messageId", messageId))
//	searchResult, err := search.EsClient.Search().Index(search.GetCoreIdx1()).Type(search.GetCoreType1()).Query(boolQuery).Do(context.Background())
//	if err != nil {
//		return err
//	}
//	messageReceiveList := make([]TMessageReceive, 0)
//	var messageReceiver TMessageReceive
//	for _, item := range searchResult.Each(reflect.TypeOf(messageReceiver)) {
//		t := item.(TMessageReceive)
//		messageReceiveList = append(messageReceiveList, t)
//	}
//	if len(messageReceiveList) > 0 {
//		for _, messageReceive := range messageReceiveList {
//			handlestatus_temp := messageReceive.HandleStatus
//			if handlestatus_temp == 0 {
//				//将要修改的消息 放入chanal中处理
//				tm := TMessageReceiverCh{
//					HandleStatus: int32(handlestatus),
//					HandleTime:   NowTimeStr(),
//					MessageId:    messageId,
//					ReadStatus:   1,
//					ReadTime:     NowTimeStr(),
//				}
//				err = tm.AddChannelUpdate() //异步向channel 中存放数据更新消息状态
//				if err != nil {
//					return err
//				}
//				err := messageReceive.UpdateDocForHandleStatusMuti(&tm)
//				if err != nil {
//					return err
//				}
//			}
//		}
//		return nil
//	}
//	return nil
//}

func SendMessage(messages []TMessage, messageReceiveList []TMessageReceive) error {
	if len(messageReceiveList) > 0 && len(messages) > 0 {
		//放入 Channel 中异步处理 message
		err := AddChannelInsert(messages)
		if err != nil {
			return err
		}
		for idx, _ := range messageReceiveList {
			id_, err := uuid.NewV4()
			if err != nil {
				return err
			}
			messageReceiveList[idx].Id = id_.String()
			messageReceiveList[idx].CreateTime = NowTimeStr()
			messageReceiveList[idx].HandleTime = "0"
			messageReceiveList[idx].ReadTime = "0"

		}
		//放入 Channel 中异步处理 message
		err = AddChannelInsert(messageReceiveList)
		if err != nil {
			return err
		}

		for _, message := range messages {
			err = message.AppendDoc()
			if err != nil {
				return err
				//return message.checkTxErr(tx, nil, err)
			}
		}
		for _, messageReceiver := range messageReceiveList {
			err = messageReceiver.AppendDoc()
			if err != nil {
				return err
				//return messageReceiver.checkTxErr(tx, nil, err)
			}
		}
		return nil
	}
	return nil
}
