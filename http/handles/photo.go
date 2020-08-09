package handles

import (
	"fmt"
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
