package db

var PhotoTable = "t_photo"

type TPhoto struct {
	Id         int32  `json:"id" xorm:"id"`
	FileId     string `json:"fileId" xorm:"fileId"`
	UserId     string `json:"userId" xorm:"userId"`
	FileType   int32  `json:"fileType" xorm:"fileType"`
	FileDir    string `json:"fileDir" xorm:"fileDir"`
	FileSite   string `json:"fileSite" xorm:"fileSite"`
	Cover      string `json:"cover" xorm:"cover"`
	FileUrl    string `json:"fileUrl" xorm:"fileUrl"`
	CreateTime string `json:"createTime" xorm:"createTime"`
}

func InsertPhoto(t *TPhoto) error {
	tx := engineSchool.NewSession()
	total, err := tx.Table(PhotoTable).Insert(t)

	if err != nil {
		return tx.Rollback()
	}
	if total <= 0 {
		return tx.Rollback()
	}
	return tx.Commit()
}

func GetPhotoDirList(userId string, fileType int) ([]TPhoto, error) {
	d := make([]TPhoto, 0, 4)
	err:= engineSchool.Table(PhotoTable).Where("fileType=? and  userId=?", fileType, userId).Desc("createTime").Limit(100, 0).Find(&d)

	if err!=nil{
		return nil,err
	}
	return d, nil

}
