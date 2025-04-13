package dao

import (
	"encoding/json"
	"filepro/config"
	"filepro/db"
	"filepro/utils"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
)

type ffff struct {
	Code string `json:"code"`
}

func Dao_change(c *fiber.Ctx) error {
	fffff := c.Body()
	var fs ffff
	json.Unmarshal(fffff, &fs)
	db.PickupSet(fs.Code)
	return c.SendStatus(fiber.StatusOK)
}

func Dao_delete(c *fiber.Ctx) error {
	fffff := c.Body()
	var fs ffff
	json.Unmarshal(fffff, &fs)
	db.PickupDelete(fs.Code)
	return c.SendStatus(fiber.StatusOK)
}

type fffffs struct {
	Status string `json:"status"`
}

func Dao_changeall(c *fiber.Ctx) error {
	fffff := c.Body()
	var fs fffffs
	json.Unmarshal(fffff, &fs)
	if fs.Status == "true" {
		db.PickSetAll(true)
	} else {
		db.PickSetAll(false)
	}
	return c.SendStatus(fiber.StatusOK)
}

func Dao_deleteall(c *fiber.Ctx) error {
	db.PickupDelteAll()
	return c.SendStatus(fiber.StatusOK)
}

type ConfigData struct {
	AdminPassword       string `json:"adminPassword"`      //密码
	ExpireTime          string `json:"expireTime"`         //过期时间
	SingleUploadLimit   string `json:"singleUploadLimit"`  //上传数量
	FileSizeLimit       string `json:"fileSizeLimit"`      //上传文件大小
	DownloadCountLimit  string `json:"downloadCountLimit"` //下载的次数限制
	FileCodeNum         string `json:"fileCodeNum"`        //取件码位数
	UploadRestriction   bool   `json:"uploadRestriction"`
	DownloadRestriction bool   `json:"downloadRestriction"`
}

func Dao_savesetting(c *fiber.Ctx) error {
	fffff := c.Body()
	var fs ConfigData
	json.Unmarshal(fffff, &fs)

	i2, _ := strconv.ParseInt(fs.ExpireTime, 10, 64)
	i3, _ := strconv.Atoi(fs.SingleUploadLimit)
	i4, _ := strconv.Atoi(fs.DownloadCountLimit)
	i5, _ := strconv.Atoi(fs.FileCodeNum)
	i6, _ := strconv.ParseInt(fs.FileSizeLimit, 10, 64)
	config.ConfigLock.Lock()
	if fs.AdminPassword != config.Config.AdminPassword {
		config.Config = db.Config{
			AdminPassword:     fs.AdminPassword,
			ExpirationTime:    i2,
			SingleUploadLimit: i3,
			FileSizeLimit:     i6,
			IsUpload:          fs.UploadRestriction,
			IsDownload:        fs.DownloadRestriction,
			DownloadLimit:     i4,
			PickupNum:         i5,
			AdminToken:        utils.GenerateRandomString(12),
		}
	} else {
		config.Config = db.Config{
			AdminPassword:     fs.AdminPassword,
			ExpirationTime:    i2,
			SingleUploadLimit: i3,
			FileSizeLimit:     i6,
			IsUpload:          fs.UploadRestriction,
			IsDownload:        fs.DownloadRestriction,
			DownloadLimit:     i4,
			PickupNum:         i5,
			AdminToken:        config.Config.AdminToken,
		}
	}
	db.ConfigSet(config.Config)
	config.ConfigLock.Unlock()

	return c.SendStatus(fiber.StatusOK)
}

func Dao_recover(c *fiber.Ctx) error {
	os.RemoveAll("./data")
	os.RemoveAll("./qr/pic")
	os.RemoveAll("./zips")
	utils.Createdir("data")
	utils.Createdir("zips")
	utils.Createdir("qr/pic")
	dbs := db.DbInstance
	dbs.Close()
	os.RemoveAll("./db/badgerdb")
	db.Start()
	config.InitConfig()
	return c.SendStatus(fiber.StatusOK)
}
