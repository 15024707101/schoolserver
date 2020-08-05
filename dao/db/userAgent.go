package db

var LoginHistoryTable = "t_login_history"

type TUserAgent struct {
	Id             int32  `json:"id" xorm:"id"`
	UserId         string `json:"userId" xorm:"userId"`
	Name           string `json:"name" xorm:"name"`
	LoginTime      string `json:"loginTime" xorm:"loginTime"`
	LoginEquipment string `json:"loginEquipment" xorm:"loginEquipment"`
	LoginAddress   string `json:"loginAddress" xorm:"loginAddress"`
	UserAgent      string `json:"userAgent" xorm:"userAgent"`
	PwdLevel       int32  `json:"pwdLevel" xorm:"pwdLevel"`
}

func InsertLoginHistory(tg *TUserAgent) error {
	tx := engineSchool.NewSession()
	total, err := tx.Table(LoginHistoryTable).Insert(tg)

	if err != nil {
		return tx.Rollback()
	}
	if total <= 0 {
		return tx.Rollback()
	}
	return tx.Commit()
}

//人公用的 插入方法
func (tt *TUserAgent) Insert(f func() error) error {
	tx := engine.NewSession()
	//if err != nil {
	//	return err
	//}
	defer tx.Close()
	err := engineWrap(tx, func() error {
		//tMap := make(map[string]interface{})
		//err := utils.Struct2MapByTagDb(tMap, tt)
		//if err != nil {
		//	return err
		//}
		_, err:= tx.Table("t_biz_transfer").Insert(tt)
		//_, err = insertTransfer.Exec()
		if err != nil {
			return err
		}

		if f != nil {
			e := f()
			if e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func GetloginHistory(userId  string) ([]TUserAgent, error) {
	d := make([]TUserAgent, 0, 4)
	_ = engineSchool.Table(LoginHistoryTable).Where("userId=?", userId).Desc("loginTime").Limit(30, 0).Find(&d)
	return d, nil
}
