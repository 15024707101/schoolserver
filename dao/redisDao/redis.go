package redisDao

import (
	"errors"
	"github.com/go-redis/redis"
	"schoolserver/config"
	"time"
)

const (
	OtherLoginKey    = "OtherLogin:"
	ScanLoginKey     = "scanLoginKey:"
	LoginInfoScanKey = "logInfoScanKey:"

	SessionSet = "sessionset_"
)

const (
	otherLoginKeyDuration    = time.Second*60*60*24 - 1
	loginInfoScanKeyDuration = time.Second * 30
)

var (
	Client *redis.Client
)

func InitWithConfig(config *conf.RedisOptions) error {
	if config == nil {
		panic("redisConfig is nil")
	}
	Client = redis.NewClient(&redis.Options{
		DB:       config.DB,
		Password: config.Password,
		Addr:     config.Addr,
		PoolSize: config.PoolSize,
	})
	if Client == nil {
		return errors.New("create redis Client failed")
	}
	return nil
}

func SetOtherLoginKey(userId string, data []byte) (string, error) {

	return Client.Set(OtherLoginKey+userId, data, otherLoginKeyDuration).Result()
}
func RemoveOtherLoginKey(userId string, data []byte) (int64, error) {
	exists, _ := Client.Exists(OtherLoginKey + userId).Result()
	if exists > 0 {
		return Client.Del(OtherLoginKey + userId).Result()
	}
	return 0, nil
}
func GetOtherLoginKey(userId string) (string, error) {

	return Client.Get(OtherLoginKey + userId).Result()
}

func SetLoginInfoScanKey(uuid string, sessionId []byte) (string, error) {
	return Client.Set(LoginInfoScanKey+uuid, sessionId, loginInfoScanKeyDuration).Result()
}

func RemoveChangeUser(userId string) (err error) {
	if len(userId) == 0 {
		return errors.New("empty userId")
	}

	sessionIds, err := Client.SMembers(SessionSet + userId).Result()
	if err != nil && err == redis.Nil {
		return
	} else if err != nil {
		return
	}
	if len(sessionIds) > 0 {
		for _, sessionId := range sessionIds {
			err = Client.Expire(sessionId, -1).Err()
			if err != nil {
				return
			}
		}
		err = Client.Expire(SessionSet+userId, -1).Err()
		if err != nil {
			return
		}
	}
	return
}
