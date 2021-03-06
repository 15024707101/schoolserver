package db

var PhotoTable = "t_photo"

type TPhoto struct {
	Id         int32  `json:"id" xorm:"id"`
	FileId     string `json:"fileId" xorm:"fileId"`
	UserId     string `json:"userId" xorm:"userId"`
	FileType   int32  `json:"fileType" xorm:"fileType"`
	AlbumName  string `json:"albumName" xorm:"albumName"`
	FileSite   string `json:"fileSite" xorm:"fileSite"` //（相册或照片在磁盘中的位置，便于找到并删除）
	Cover      string `json:"cover" xorm:"cover"`       //封面
	FileUrl    string `json:"fileUrl" xorm:"fileUrl"`
	CreateTime string `json:"createTime" xorm:"createTime"`
	PhotoCount int32  `json:"photoCount" xorm:"-"`
}

type Photo struct {
	AlbumName  string `json:"albumName" xorm:"albumName"`
	PhotoCount int32  `json:"photoCount" xorm:"photoCount"`
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
	err := engineSchool.Table(PhotoTable).Where("fileType=? and  userId=?", fileType, userId).Desc("createTime").Limit(100, 0).Find(&d)

	if err != nil {
		return nil, err
	}
	return d, nil

}

func GetPhotoList(userId string, albumName string, fileType int) ([]TPhoto, error) {
	d := make([]TPhoto, 0, 4)
	err := engineSchool.Table(PhotoTable).Where("fileType=? and  albumName=? and  userId=?", fileType, albumName, userId).Desc("createTime").Limit(100, 0).Find(&d)

	if err != nil {
		return nil, err
	}
	return d, nil

}

func GetPhotoCount() ([]Photo, error) {
	d := make([]Photo, 0, 4)
	err := engineSchool.SQL("select  albumName ,COUNT(*)  photoCount   from  t_photo where fileType!=1  GROUP BY  albumName  ").Find(&d)

	if err != nil {
		return nil, err
	}
	return d, nil

}

//创建完相册后将 封面的 所属相册 改为当前相册
func UpdetePhotoCover(cover string, albumName string) error {
	tx := engineSchool.NewSession()
	tp := TPhoto{
		AlbumName: albumName,
	}
	_, err := tx.Table(PhotoTable).Where("fileUrl=?", cover).Update(tp)

	//_, err = update.Exec()
	if err != nil {
		return err
	}
	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func DeletePhoto(fileSite string) error {
	tx := engineSchool.NewSession()
	tp := TPhoto{
		FileSite: fileSite,
	}

	_, err := tx.Table(PhotoTable).Where("fileSite=?", fileSite).Delete(tp)

	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
