package db

type TClass struct {
	Id             int32  `json:"id" xorm:"id"`
	ClassId        string `json:"classId" xorm:"classId"`
	ClassName      string `json:"className" xorm:"className"`
	ClassRoom      string `json:"classRoom" xorm:"classRoom"`
	ClassNum       int32  `json:"classNum" xorm:"classNum"`
	ClassTeacher   string `json:"classTeacher" xorm:"classTeacher"`
	Classregulator string `json:"classregulator" xorm:"classregulator"`
	SpecialtyId    string `json:"specialtyId" xorm:"specialtyId"`
}

func GetClassList() ([]TClass, error) {
	d := make([]TClass, 0, 4)
	_ = engineSchool.Table(ClassTable).Where("1=?", 1).Limit(100, 0).Find(&d)
	return d, nil
}
