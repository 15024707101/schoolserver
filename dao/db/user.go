package db

var (
	UserTable      = "t_user"
	ClassTable     = "t_class"
	SpecialtyTable = "t_specialty"
)

type TUser struct {
	Id             int32  `json:"id" xorm:"id"`
	UserId         string `json:"userId" xorm:"userId"`
	Name           string `json:"name" xorm:"name"`
	IdentityCardNo string `json:"identityCardNo" xorm:"identityCardNo"`
	Mobile         string `json:"mobile" xorm:"mobile"`
	Pwd            string `json:"pwd" xorm:"pwd"`
	Sex            string `json:"sex" xorm:"sex"`
	ClassId        string `json:"classId" xorm:"classId"`
	Age            int32  `json:"age" xorm:"age"`
	PersonType     int32  `json:"personType" xorm:"personType"`
	Status         int32  `json:"status" xorm:"status"`
	CreateTime     string `json:"createTime" xorm:"createTime"`
	HeadPortrait   string `json:"headPortrait" xorm:"headPortrait"`
}

func (u *TUser) LoginByPwd() (bool, error) {
	ok, err := engineSchool.Table(UserTable).Where("userId=? and pwd =? ", u.UserId, u.Pwd).Get(u)
	return ok, err
}

func GetUserList() ([]TUser, error) {
	d := make([]TUser, 0, 4)
	_ = engineSchool.Table(UserTable).Where("status=?", 1).Desc("createTime").Limit(100, 0).Find(&d)
	return d, nil

}
