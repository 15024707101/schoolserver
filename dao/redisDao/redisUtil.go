package redisDao

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"time"
)

/*
  向redis中存贮数据
  rediskey:redis 中存储的key
  data：要解析并取出的数据

*/
func SaveRedisString(rediskey string, data interface{}, expiresTime time.Duration) (err error) {
	//向redis存贮数据
	jsonspeak, err := json.Marshal(data)
	if err != nil {
		return echo.NewHTTPError(9011, "格式化数据出错")
	}
	//过期时间设置为10天
	_, err = Client.Set(rediskey, string(jsonspeak), expiresTime).Result()
	if err != nil {
		return echo.NewHTTPError(9011, "存储数据出错")
	}
	return nil
}

//设置redis 的key 的超时时间
func SaveRedisExpire(rediskey string, expiresTime time.Duration) {

	//过期时间设置为10天
	Client.Expire(rediskey, expiresTime)
	return
}

/*从redis中获取数据
  rediskey:redis 中存储的key
  data：要解析并存入的数据

*/
func GetRedisString(rediskey string, data interface{}) error {
	speakInfo, err := Client.Get(rediskey).Result()
	if err != nil {
		if err != redis.Nil {
			return err
		}
	} else if speakInfo != "" {

		err = json.Unmarshal([]byte(speakInfo), data)
		if err != nil {
			return nil
		}
	}

	return nil
}

/*
  单个删除（rediskey 为key值）
*/
func DeleteRedisBuyKey(rediskey string) error {
	/*speakInfo, err := service.Client.Keys(rediskey).Result()
	log.Debug(speakInfo)
	if err != nil {
		log.Error("删除redis异常", err)
	}*/
	if len(rediskey) != 0 {
		Client.Del(rediskey)
	}
	return nil
}

//向redis中存贮hash类型数据
func SaveRedisHash(key string, field string, data interface{}, expiresTime time.Duration) (err error) {
	//向redis存贮数据
	jsonspeak, err := json.Marshal(data)
	if err != nil {
		return echo.NewHTTPError(9011, "格式化数据出错")
	}
	//过期时间设置为10天
	_, err = Client.HSet(key, field, string(jsonspeak)).Result()
	Client.Expire(key, expiresTime)
	if err != nil {
		return echo.NewHTTPError(9011, "存储数据出错")
	}
	return nil
}

//从redis中获取hash类型数据
func GetRedisHash(key string, field string, data interface{}) error {
	speakInfo, err := Client.HGet(key, field).Result()
	if err != nil {
		if err != redis.Nil {
			return err
		}
	} else if speakInfo != "" {

		err = json.Unmarshal([]byte(speakInfo), data)
		if err != nil {
			return nil
		}
	}

	return nil
}

//删除hash类型数据
func DeleteRedisHash(rediskey string) error {
	if len(rediskey) != 0 {
		Client.Expire(rediskey, 0)
	}
	return nil
}
