package handles

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"schoolserver/common/ecode"
	"schoolserver/dao/db"
)

func AppendPhoto(c echo.Context) error {
	photoDir := c.FormValue("photoDir") //上传文件的目录
	userId := c.FormValue("userId")     //上传文件的目录

	//上传图片
	gofile, err := UploadImgUtil(c)
	if err != nil {
		return echo.NewHTTPError(67918, err)
	}
	//记录到数据库

	tp := db.TPhoto{}
	tp.CreateTime = db.NowTimeStr()
	tp.UserId = userId
	tp.FileId = gofile.FileId
	tp.FileDir = FileDirString + photoDir
	tp.FileSite = gofile.FilePath
	tp.FileUrl = gofile.FileUrl
	tp.FileType = 2
	tp.Cover = ""

	err = db.InsertPhoto(&tp)
	if err != nil {
		return FailWithMsg(c, 4002, fmt.Sprintf("添加到数据库时候发生异常：%v", err))
	}

	return Success(c, ecode.OK, gofile)
}

func CreateAlbum(c echo.Context) error {

	cover := c.FormValue("cover")
	albumName := c.FormValue("albumName")
	userId := c.FormValue("userId")

	id_, err := uuid.NewV4()
	if err != nil {
		return FailWithMsg(c, 4002, fmt.Sprintf("生成uuid时异常：%v", err))
	}

	tp := db.TPhoto{}
	tp.CreateTime = db.NowTimeStr()
	tp.UserId = userId
	tp.FileDir = albumName
	tp.FileType = 1
	tp.Cover = cover
	tp.FileId = id_.String()

	err = db.InsertPhoto(&tp)
	if err != nil {
		return FailWithMsg(c, 4002, fmt.Sprintf("添加到数据库时候发生异常：%v", err))
	}

	return Success(c, ecode.OK, "相册创建成功")
}

func GetPhotoDirList(c echo.Context) error {

	userId := c.FormValue("userId")

	d := make([]db.TPhoto, 0, 4)
	d, err := db.GetPhotoDirList(userId, 1)
	if err != nil {
		return FailWithMsg(c, 4002, fmt.Sprintf("添加到数据库时候发生异常：%v", err))
	}

	return Success(c, ecode.OK, d)
}

func GetPhotoList(c echo.Context) error {

	userId := c.FormValue("userId")

	d := make([]db.TPhoto, 0, 4)
	d, err := db.GetPhotoDirList(userId, 2)
	if err != nil {
		return FailWithMsg(c, 4002, fmt.Sprintf("查询数据库时候发生异常：%v", err))
	}

	return Success(c, ecode.OK, d)
}
