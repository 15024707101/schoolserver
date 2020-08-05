package search

import (
	"context"
	"errors"
	"fmt"
	conf "schoolserver/config"
	"schoolserver/logger"

	"github.com/olivere/elastic"

	"sync"
)

var (
	filename     string = "./logs/es.log"
	limitCount   int64  = 50 // 每次往后端msgpack传送的记录数
	elasticHosts []conf.EsConfig
)

var core0index string = "fulltextsearchcore0"
var core1index string = "fulltextsearchcore1"
var core0type string = "TMessage"
var core1type string = "TMessageReceiver"
var EsClient *elastic.Client

func GetCoreIdx1() string {
	return core1index
}

func GetCoreIdx0() string {
	return core0index
}

func GetCoreType0() string {
	return core0type
}

func GetCoreType1() string {
	return core1type
}

type L struct{}

func (l *L) Printf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

var once sync.Once
var errEs error

func onceInitEsConn() {
	var err error
	EsClient, err = elastic.NewClient(
		elastic.SetURL(getElasticHost()...),
		elastic.SetSniff(true),
		elastic.SetHealthcheck(true),
		//elastic.SetHealthcheckInterval(10*time.Second),
		//elastic.SetMaxRetries(5),
		elastic.SetErrorLog(&L{}),
		elastic.SetInfoLog(&L{}))

	if err != nil {
		errEs = err
		return
	}
	return
}
func InitEsConn() error {
	once.Do(onceInitEsConn)
	return errEs
}

func GetHomeUnCount(userId, leagueId string) (map[string]int64, error) {
	InitEsConn()
	err := errEs
	if err != nil {
		return nil, err
	}
	client := EsClient
	retMap := make(map[string]int64)
	count, err := GetHomeUnReadCount(userId, leagueId, client)
	if err != nil {
		return nil, err
	}

	retMap["unRead"] = count
	count, err = GetHomeUnHandleCount(userId, leagueId, client)
	if err != nil {
		return nil, err
	}
	retMap["unHandle"] = count

	return retMap, nil
}
func GetHomeUnReadCount(userId, leagueId string, client *elastic.Client) (int64, error) {
	//client, err := elastic.NewClient(elastic.SetURL(getElasticHost()))
	//if err != nil {
	//	return -1, err
	//}
	if client == nil {
		return -1, errors.New("elastic.Clent is nil")
	}
	if len(leagueId) > 0 {
		boolQuery := elastic.NewBoolQuery()
		termQueries := []elastic.Query{
			elastic.NewTermQuery("receiverType", 1),
			elastic.NewTermQuery("leagueId", leagueId),
		}
		boolQuery = boolQuery.Must(termQueries...)

		boolQuery = boolQuery.Filter(elastic.NewTermQuery("readStatus", 0))
		numForLeague, err := client.Count(core1index).Type(core1type).Query(boolQuery).Do(context.Background())

		if err != nil {
			return -1, err
		}
		return numForLeague, nil

	} else if len(userId) > 0 {
		var queryStr = `{
			"query":{
				"bool":{
					"must":[
						{"term":{"receiverType":"0"}},
						{"term":{"userId":"%s"}}
					],
				   "filter":{
						"term":{"readStatus":0}
					}
				}
			}
		}`
		queryStr = fmt.Sprintf(queryStr, userId)
		numForUser, err := client.Count(core1index).Type(core1type).BodyString(queryStr).Do(context.Background())
		if err != nil {
			return -1, err
		}
		return numForUser, nil
	}

	return 0, nil
}

func GetNotMemberHomeUnReadCount(userId string, client *elastic.Client) (int64, error) {
	//client, err := elastic.NewClient(elastic.SetURL(getElasticHost()))
	//if err != nil {
	//	return -1, err
	//}
	if client == nil {
		return -1, errors.New("elastic.Clent is nil")
	}

	if len(userId) <= 0 {
		return 0, nil
	}

	var queryStr = `{
			"query":{
				"bool":{
					"must":[
						{"term":{"receiverType":"0"}},
						{"term":{"userId":"%s"}}
					],
				   "filter":{
						"term":{"readStatus":0}
					}
				}
			}
		}`
	queryStr = fmt.Sprintf(queryStr, userId)
	numForUser, err := client.Count(core1index).Type(core1type).BodyString(queryStr).Do(context.Background())
	if err != nil {
		return -1, err
	}

	return numForUser, nil
}

func GetHomeUnHandleCount(userId, leagueId string, client *elastic.Client) (int64, error) {
	//client, err := elastic.NewClient(elastic.SetURL(getElasticHost()))
	//if err != nil {
	//	return -1, err
	//}
	if client == nil {
		return -1, errors.New("elastic.Clent is nil")
	}

	var queryStr string
	if len(leagueId) > 0 {
		queryStr = `{
		   "query": {
				"bool":{
					"should":[
						{
							"bool":{
								"must_not":{
									"term":{"userId":"%s"}
								},
								"must":[
									{"term":{"receiverType":"1"}},
		
									{"term":{"leagueId":"%s"}}
								],
								"filter":{
									"term":{"handleStatus":0}
								}
							}
						},
						{
							"bool":{
								"must":[
									{"term":{"receiverType":"0"}},
									{"term":{"userId":"%s"}}
								],
							   "filter":{
									"term":{"handleStatus":0}
								}
							}
						}
					],
					"minimum_should_match": 1
				}
			}
		}`
		queryStr = fmt.Sprintf(queryStr, userId, leagueId, userId)
	} else {
		queryStr = `{
			"query":{
				"bool":{
					"must":[
						{"term":{"receiverType":"0"}},
						{"term":{"userId":"%s"}}
					],
				   "filter":{
						"term":{"handleStatus":0}
					}
				}
			}
		}`
		queryStr = fmt.Sprintf(queryStr, userId)
	}

	return client.Count(core1index).Type(core1type).BodyString(queryStr).Do(context.Background())
}
func InitES(config *conf.AllConfig) {
	elasticHosts = config.Es

}
func getElasticHost() []string {

	urls := make([]string, 0, len(elasticHosts))
	//nth := rand.Intn(len(elasticHosts))
	//host := fmt.Sprintf("%s:%d", elasticHosts[nth].EsHost, elasticHosts[nth].EsPort)
	for _, host := range elasticHosts {
		url := fmt.Sprintf("%s:%d", host.EsHost, host.EsPort)
		urls = append(urls, url)
	}

	return urls
}
