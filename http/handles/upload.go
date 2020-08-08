package handles

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"os"
	"path"
	"path/filepath"
	"schoolserver/common/ecode"
)

func UploadImg(c echo.Context) error {
	// Read form fields
	name := c.FormValue("name")
	photoDir := c.FormValue("photoDir")
	fmt.Println(name)
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fileExt := path.Ext(file.Filename)
	/*fileExtLower := strings.ToLower(fileExt)*/
	// Destination
	/*dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()*/

	id_, err := uuid.NewV4()
	if err != nil {
		return echo.NewHTTPError(67918, "获取数据发生错误，请联系系统管理员")
	}
	fname := id_.String() + fileExt
	realPath := filepath.Join("static/img"+photoDir, fname)
	outputFile, err := os.OpenFile(realPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return echo.NewHTTPError(7782, "创建文件失败")
	}
	defer outputFile.Close()

	// Copy
	if _, err = io.Copy(outputFile, src); err != nil {
		return err
	}
	ret := struct {
		FileId           string `json:"fileId"`
		FileUrl          string `json:"fileUrl"`
		FileName         string `json:"fileName"`
		ThumbnailFileId  string `json:"thumbnailFileId"`
		ThumbnailFileUrl string `json:"thumbnailFileUrl"`
	}{
		"wqpei",
		"http://127.0.0.1:3334/img"+photoDir+"/"+fname,
		"我的图片",
		"7845",
		"http://127.0.0.1:3334/img"+photoDir+"/"+fname,
	}
	return Success(c, ecode.OK, ret)
}

//func UploadTwoRedFileHandler(c echo.Context) error {
//
//	currLeagueId := c.Get(middleware.CtxLeagueId).(string)
//	loginedInfo := c.Get(middleware.CtxUser).(*backend.LoginResp)
//
//	file, fileHeader, err := c.Request().FormFile("file")
//	//log.Debug("uploaded file name:%s",fileHeader.Filename)
//	if err != nil {
//		return echo.NewHTTPError(6789, fmt.Errorf("读取文件错误:%v", err))
//	}
//	defer file.Close()
//	defer c.Request().MultipartForm.RemoveAll() // 删除上传的临时文件
//	fileExt := path.Ext(fileHeader.Filename)
//	fileExtLower := strings.ToLower(fileExt)
//	if fileExtLower != ".jpg" && fileExtLower != ".jpeg" && fileExtLower != ".gif" && fileExtLower != ".png" {
//		return echo.NewHTTPError(6667, fmt.Sprintf("文件类型%s不正确", path.Ext(fileHeader.Filename)))
//	}
//	sizeMB := math.Ceil(float64(fileHeader.Size)*100/(1024*1024)) / 100
//	if sizeMB > 10 { // 不能大于 10M
//		return echo.NewHTTPError(7789, "上传文件不能大于10M")
//	}
//
//	// 写入磁盘
//	uridate := time.Now().Format("20060102")
//	//realPath := filepath.Join(string(filepath.Separator),service.Cfg.Sys.UploadLocalDir,"meeting",
//	//	time.Now().Format("20060102"))
//	realPath := filepath.Join(string(filepath.Separator), service.Cfg.Sys.UploadLocalDir, "meeting", uridate)
//	err = os.MkdirAll(realPath, os.ModePerm)
//	if err != nil {
//		log.Error("7780,os.MKdirAll error:%v", err)
//		return echo.NewHTTPError(7780, "创建目录出错")
//	}
//	remoteRealPath := filepath.Join(string(filepath.Separator), service.Cfg.Sys.UploadRemoteDir, "meeting_indent", uridate)
//	inputFile, err := fileHeader.Open()
//	if err != nil {
//		log.Error("7781,open file error:%v", err)
//		return echo.NewHTTPError(7781, "打开传入文件失败")
//	}
//	defer inputFile.Close()
//	id_, err := uuid.NewV4()
//	if err != nil {
//		return echo.NewHTTPError(67918, "获取数据发生错误，请联系系统管理员")
//	}
//	fname := id_.String() + fileExt
//	realPath = filepath.Join(realPath, fname)
//	log.Debug("realPath:%s", realPath)
//	outputFile, err := os.OpenFile(realPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
//	if err != nil {
//		log.Error("7782,os.OpenFile error:%v", err)
//		return echo.NewHTTPError(7782, "创建文件失败")
//	}
//	defer outputFile.Close()
//	inputFile.Seek(0, 0)
//	_, err = io.Copy(outputFile, inputFile)
//	if err != nil {
//		log.Error("7783,io.copy error:%v", err)
//		return echo.NewHTTPError(7783, "写入文件内容失败")
//	}
//	//type re struct {
//	//	backend.HeaderResp
//	//	Url string `json:"url"`
//	//}
//	uri := "/meeting/" + uridate + "/" + fname
//	//ret := re{
//	//	backend.HeaderResp{Status:1},
//	//	uri,
//	//}
//	//return RespCheck(ret,c)+
//	realPath2 := filepath.Join(string(filepath.Separator), service.Cfg.Sys.UploadLocalDir, "meeting_indent", uridate)
//	err = os.MkdirAll(realPath2, os.ModePerm)
//	if err != nil {
//		log.Error("7780,os.MKdirAll error:%v", err)
//		return echo.NewHTTPError(7780, "创建目录出错")
//	}
//	uuid_ := utils.RandSeq()
//	var (
//		lang      int32  = 1
//		ip        string = c.RealIP()
//		deviceID  string = loginedInfo.DeviceId
//		sessionID string = loginedInfo.SessionId
//	)
//	srcPath := service.Cfg.Sys.UploadLocalDir + uri
//	fname2 := id_.String() + path.Ext(path.Ext(srcPath))
//	destPath := filepath.Join(realPath2, fname2)
//	err = utils.ResizeImage(srcPath, destPath, 150, 90)
//	if err != nil {
//		log.Error("9091,ResizeImage error: %v", err)
//		return echo.NewHTTPError(9091, "resize图片时发生错误")
//	}
//	destRemotePath := filepath.Join(remoteRealPath, fname2)
//
//	msgPackReq1501 := backend.FileUploadReq{
//		Code:             backend.CodeFileUpload,
//		SeqNo:            uuid_,
//		Lang:             lang,
//		Ip:               ip,
//		DeviceId:         deviceID,
//		SessionId:        sessionID,
//		CurrentLeagueId:  currLeagueId,
//		FileLocalAddress: destRemotePath,
//	}
//	msgPackResp1501 := backend.FileUploadResp{}
//	respErr := backend.GetDataWithCtx(&msgPackReq1501, &msgPackResp1501, c)
//	if respErr != nil {
//		log.Error("9092,backend.GetDataWithCtx failed:%v", respErr)
//		return echo.NewHTTPError(9092, "上传图片到后端服务器时发生错误")
//	}
//	if msgPackResp1501.Status != 1 {
//		log.Error("9093,backend.GetDataWithCtx Return:%d", msgPackResp1501.Status)
//		return echo.NewHTTPError(9093, fmt.Sprintf("后端服务返回错误：%s", msgPackResp1501.Desc))
//	}
//
//	width, _, err := utils.GetImageInfo(srcPath)
//	if err != nil {
//		log.Error("9090, GetImageInfo error:%v", err)
//		return echo.NewHTTPError(9090, "获取上传图片信息时失败")
//	}
//
//	fileIdPath := service.Cfg.Sys.UploadRemoteDir + uri
//	if width >= 1080 {
//		id_, err := uuid.NewV4()
//		if err != nil {
//			return echo.NewHTTPError(67918, "获取数据发生错误，请联系系统管理员")
//		}
//		fname := id_.String() + path.Ext(path.Ext(srcPath))
//		fileIdPath_ := filepath.Join(realPath2, fname)
//		err = utils.ResizeImage(srcPath, fileIdPath_, 1080, 1080)
//		if err != nil {
//			log.Error("9091,ResizeImage error: %v", err)
//			return echo.NewHTTPError(9091, "resize 1080图片时发生错误")
//		}
//		fileIdPath = filepath.Join(remoteRealPath, fname)
//	}
//	uuid_2 := utils.RandSeq()
//	msgPackReq1501.SeqNo = uuid_2
//	msgPackReq1501.FileLocalAddress = fileIdPath
//	msgPackResp1501_2 := backend.FileUploadResp{}
//	respErr = backend.GetDataWithCtx(&msgPackReq1501, &msgPackResp1501_2, c)
//	if respErr != nil {
//		log.Error("9094,backend.GetDataWithCtx failed:%v", respErr)
//		return echo.NewHTTPError(9094, "上传图片到后端服务器时发生错误")
//	}
//	if msgPackResp1501_2.Status != 1 {
//		log.Error("9095,backend.GetDataWithCtx Return:%d", msgPackResp1501.Status)
//		return echo.NewHTTPError(9095, fmt.Sprintf("后端服务返回错误：%s", msgPackResp1501.Desc))
//	}
//
//	ret := struct {
//		backend.HeaderResp
//		FileId           string `json:"fileId"`
//		FileUrl          string `json:"fileUrl"`
//		FileName         string `json:"fileName"`
//		ThumbnailFileId  string `json:"thumbnailFileId"`
//		ThumbnailFileUrl string `json:"thumbnailFileUrl"`
//	}{
//		msgPackResp1501_2.HeaderResp,
//		msgPackResp1501_2.FileId,
//		msgPackResp1501_2.FileUrl,
//		fname2,
//		msgPackResp1501.FileId,
//		msgPackResp1501.FileUrl,
//	}
//	return RespCheck(ret, c)
//
//}

/*func MemberRecordUpload(c echo.Context) error {
	curLeague := c.Get(middleware.CtxLeague).(*backend.JoinedLeague)
	loginedInfo := c.Get(middleware.CtxUser).(*backend.LoginResp)
	leagueId := c.FormValue("leagueId")
	userId := c.FormValue("userId")
	name := c.FormValue("name")
	userCode := c.FormValue("userCode")
	log.Info(name, ":", userCode)

	err := c.Request().ParseMultipartForm(int64(maxSizeOfDocument)<<20 + 512)
	if err != nil {
		return echo.NewHTTPError(6788, fmt.Errorf("读取文件设置缓冲错误:%v", err))
	}
	fileHeader, err := c.FormFile("recordfile")
	if err != nil {
		return echo.NewHTTPError(6789, fmt.Errorf("读取文件错误:%v", err))
	}
	file, err := fileHeader.Open()
	if err != nil {
		log.Error("7781,open file error:%v", err)
		return echo.NewHTTPError(7781, "打开传入文件失败")
	}
	defer file.Close()
	extension := path.Ext(fileHeader.Filename)
	if extension != ".pdf" && extension != ".PDF" {
		return echo.NewHTTPError(6667, fmt.Sprintf("文件类型%s不正确", extension))
	}

	sizeMB := math.Ceil(float64(fileHeader.Size)*100/(1024*1024)) / 100
	if sizeMB > float64(maxSizeOfDocument) { // 不能大于 规定大小
		return echo.NewHTTPError(7789, "上传文件不能大于", maxSizeOfDocument, "M")
	}

	// 写入磁盘
	uriDate := time.Now().Format("20060102")
	dirPath := filepath.Join(string(filepath.Separator), service.Cfg.Sys.UploadLocalDir, "recordFile", uriDate)

	//文件名如：订单受理1.pdf
	newV4, _ := uuid.NewV4()
	fileName := RecordName + name + newV4.String() + ".pdf"
	localPath, err := saveFileLocally(file, fileName, dirPath)
	if err != nil {
		return echo.NewHTTPError(7783, "写入文件内容失败")
	}

	// remove tmp file, after pass to api
	defer func() {
		err := os.Remove(localPath)
		if err != nil {
			log.Error(err)
		} else {
			log.Info("remove tmp file:%s", localPath)
		}
	}()

	log.Debug("save new file to: %v", localPath)

	// check pdf file

	var pt pdft.PDFt
	err = pt.Open(localPath)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(7993, "读取文件内容失败，请检查文件是否符合规范")
	}

	if pt.GetNumberOfPage() != RequiredPdfPage {
		log.Error("upload eRecord page  not match has %d, required %d", pt.GetNumberOfPage(), RequiredPdfPage)
		return echo.NewHTTPError(7994, "必须上传12页完整的电子版《入团志愿书》")
	}

	//pdfCtx, err := pdfApi.ReadContext(file, pdfcpu.NewDefaultConfiguration())
	//if err != nil {
	//	log.Error(err)
	//	return echo.NewHTTPError(7993, "读取文件内容失败，请检查文件是否符合规范")
	//}
	//err = pdfApi.ValidateContext(pdfCtx)
	//if err != nil {
	//	log.Error(err)
	//	return echo.NewHTTPError(7995, "读取文件内容失败，请检查文件是否符合规范")
	//}
	//
	//if pdfCtx.PageCount != RequiredPdfPage {
	//	log.Error("upload eRecord has %d pages, not match %d", pdfCtx.PageCount, RequiredPdfPage)
	//	return echo.NewHTTPError(7994, "必须上传12页完整的电子版《入团志愿书》")
	//}
	//
	//err = pdfApi.OptimizeContext(pdfCtx)
	//if err != nil {
	//	log.Error(err)
	//	return echo.NewHTTPError(7995, "读取文件内容失败，请检查文件是否符合规范")
	//}

	// check images in pdf
	//pageImages := pdfCtx.Optimize.PageImages
	//images := pdfCtx.Optimize.ImageObjects
	//if len(pageImages) != RequiredPdfPage {
	//	log.Error("upload eRecord has %d pages, not match %d", pdfCtx.PageCount, RequiredPdfPage)
	//	return echo.NewHTTPError(7996, "必须上传12页完整的电子版《入团志愿书》")
	//}
	//
	//for i := 0; i < len(pageImages); i++ {
	//	imageOk := false
	//	if pageImages[i] == nil {
	//		return echo.NewHTTPError(7996, fmt.Sprintf("第%d页没有图片，请检查文件是否符合规范", i+1))
	//	} else {
	//
	//		for imgNumber, _ := range pageImages[i] {
	//			img := images[imgNumber].ImageDict
	//			h, _ := strconv.Atoi(img.Dict["Height"].String())
	//			w, _ := strconv.Atoi(img.Dict["Width"].String())
	//			if h >= imgHeightMin && w >= imgWidthMin {
	//				imageOk = true
	//			}
	//		}
	//	}
	//	if !imageOk {
	//		return echo.NewHTTPError(7996, fmt.Sprintf("第%d页的图片像素异常，请检查文件是否符合规范", i+1))
	//	}
	//}

	//档案所在的文件夹
	uri := "/recordFile/" + uriDate + "/" + fileName

	var (
		lang      int32  = 1
		ip        string = c.RealIP()
		deviceID  string = loginedInfo.DeviceId
		sessionID string = loginedInfo.SessionId
	)
	//判断当前组织是否需要发起档案流程 （1需要，2不需要） 团支部 ，团总支，毕业班需要；其余不需要
	processStatus := 2
	if curLeague.LeagueTypeId == "01TZB" || curLeague.LeagueTypeId == "02TZZ" || curLeague.LeagueTypeId == "06BYBTZB" || curLeague.LeagueTypeId == "00DJZTZB" {
		processStatus = 1
	} else {
		processStatus = 2
	}

	//调用c端接口上传电子档案，并创建档案审批流程
	fileIdPath := service.Cfg.Sys.UploadRemoteDir + uri
	msgPackReq := backend.RecordUploadReq{
		Code:            backend.CodeRecordUpload,
		Lang:            lang,
		Ip:              ip,
		DeviceId:        deviceID,
		SessionId:       sessionID,
		CurrentLeagueId: curLeague.LeagueId,
		LeagueId:        leagueId,
		UserId:          userId,
	}
	msgPackReq.SeqNo = utils.RandSeq()
	msgPackReq.FileLocalAddress = fileIdPath
	msgPackReq.Filename = fileName
	msgPackReq.ProcessStatus = int32(processStatus)
	msgPackResp := backend.RecordUploadResp{}
	respErr := backend.GetDataWithCtx(&msgPackReq, &msgPackResp, c)
	if respErr != nil {
		log.Error("9094,backend.GetDataWithCtx failed:%v", respErr)
		return echo.NewHTTPError(9094, "上传团员档案到后端服务器时发生错误")
	}
	if msgPackResp.Status != 1 {
		log.Error("9095,backend.GetDataWithCtx Return:%d", msgPackResp.Status)
		return echo.NewHTTPError(9095, msgPackResp.Desc)
	}

	//在go端创建档案审批流程

	err = CreateMemberRecord(c, msgPackResp.FileId, processStatus)
	if err != nil {
		log.Error("9094,backend.GetDataWithCtx failed:%v", err)
		return echo.NewHTTPError(9094, err)
	}
	type re struct {
		backend.HeaderResp
		FileUrl string `json:"fileUrl"`
		FileId  string `json:"fileId"`
	}
	ret := re{
		backend.HeaderResp{Status: 1},
		msgPackResp.FileUrl,
		msgPackResp.FileId,
	}
	return RespCheck(ret, c)
}*/
