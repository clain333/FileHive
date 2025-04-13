package dao

import (
	"encoding/json"
	"filepro/config"
	"filepro/db"
	"filepro/model"
	"filepro/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"time"
)

func Dao_uploadfile(c *fiber.Ctx) error {
	description := c.FormValue("description")
	password := c.FormValue("password")
	form, _ := c.MultipartForm()
	files := form.File["files"]
	switch {
	case len(files) == 0:
		return c.JSON(utils.Returninfo("10001", "上传失败", "请选择文件", "", ""))
	case password == "":
		return c.JSON(utils.Returninfo("10002", "上传失败", "请输入取件密码", "", ""))
	}

	switch len(files) {
	case 1:

		ext := filepath.Ext(files[0].Filename)
		id := utils.GenerateRandomString(6)
		c.SaveFile(files[0], fmt.Sprintf("%s/%s", model.Dir, id+ext))
		ha := utils.CreateHa(id + ext)
		err := db.FileGet(ha)
		if err != nil {
			fsss := db.File{
				HashValue:   ha,
				NewFilename: id + ext,
			}
			fsss.FileSave()
		} else {
			utils.RemoveFile(id + ext)
		}

		strings := c.BaseURL()
		code := utils.GenerateRandomString(config.Config.PickupNum)
		qrurl := fmt.Sprintf("%s/fileinfo?code=%s&psd=%s", strings, code, password)
		utils.CreateQr(qrurl, fmt.Sprintf("qr/pic/%s%s", code, ".png"))
		config.ConfigLock.Lock()
		p := db.PickupInfo{
			PickupCode:         code,
			PickupPassword:     password,
			QrCode:             "pic/" + code + ".png",
			IsDownloaded:       config.Config.IsDownload,
			FileSize:           files[0].Size,
			Text:               description,
			ExpirationTime:     time.Now().Unix() + config.Config.ExpirationTime*24*60*60,
			RemainingDownloads: config.Config.DownloadLimit,
			FileInfo: []db.FileInfo{
				db.FileInfo{
					HashValue: ha,
					Filename:  files[0].Filename,
				},
			},
		}
		config.ConfigLock.Unlock()

		p.PickupSave()

		return c.JSON(utils.Returninfo("20001", "上传成功", description, password, p.PickupCode))

	default:
		dbis := make([]db.FileInfo, 0)
		var snum int64 = 0
		filesss := []db.HashName{}
		for _, file := range files {
			snum += file.Size
			ext := filepath.Ext(file.Filename)
			id := utils.GenerateRandomString(6)
			c.SaveFile(file, fmt.Sprintf("%s/%s", model.Dir, id+ext))
			ha := utils.CreateHa(id + ext)
			err := db.FileGet(ha)
			if err != nil {
				fsss := db.File{
					HashValue:   ha,
					NewFilename: id + ext,
				}
				fsss.FileSave()
			} else {
				utils.RemoveFile(id + ext)
			}
			dbis = append(dbis, db.FileInfo{
				HashValue: ha,
				Filename:  file.Filename,
			})
			ffss := db.FileGets(ha)
			filesss = append(filesss, db.HashName{
				Oldname: file.Filename,
				Newname: ffss.NewFilename,
			})
		}
		strings := c.BaseURL()
		code := utils.GenerateRandomString(config.Config.PickupNum)
		qrurl := fmt.Sprintf("%s/fileinfo?code=%s&psd=%s", strings, code, password)
		utils.CreateQr(qrurl, fmt.Sprintf("qr/pic/%s%s", code, ".png"))
		config.ConfigLock.Lock()
		p := db.PickupInfo{
			PickupCode:         code,
			PickupPassword:     password,
			QrCode:             "pic/" + code + ".png",
			IsDownloaded:       config.Config.IsDownload,
			FileSize:           snum,
			Text:               description,
			ExpirationTime:     time.Now().Unix() + config.Config.ExpirationTime*24*60*60,
			RemainingDownloads: config.Config.DownloadLimit,
			FileInfo:           dbis,
		}
		config.ConfigLock.Unlock()
		p.PickupSave()
		utils.CreateZip(p.PickupCode, filesss)
		return c.JSON(utils.Returninfo("20001", "上传成功", description, password, p.PickupCode))
	}

}

type f struct {
	Password   string `json:"password"`
	Pickupcode string `json:"pickupcode"`
}

func Dao_getcode(c *fiber.Ctx) error {
	var fs f
	body := c.Body()
	json.Unmarshal(body, &fs)
	p, e := db.PickupGet(fs.Pickupcode)
	if e != nil {
		return c.JSON(utils.Returninfo("50001", "取件失败", "取件码不存在", "", ""))
	} else if fs.Password != p.PickupPassword {
		return c.JSON(utils.Returninfo("50001", "取件失败", "密码错误", "", ""))
	} else {
		return c.JSON(utils.Returninfo("10001", "取件成功", "", p.PickupPassword, p.PickupCode))
	}
}

type ff struct {
	Password   string `json:"psd"`
	Pickupcode string `json:"code"`
	FileNum    int64  `json:"filenum"`
}

func Dao_download(c *fiber.Ctx) error {
	var fs ff
	body := c.Body()
	json.Unmarshal(body, &fs)

	p, _ := db.PickupGet(fs.Pickupcode)
	if p.PickupPassword != fs.Password {
		c.Send([]byte("在网页下载"))
	}

	if len(p.FileInfo) == 1 {
		fss := db.FileGets(p.FileInfo[0].HashValue)
		return c.Download(fmt.Sprintf("./data/%s", fss.NewFilename))
	} else {
		return c.Download(fmt.Sprintf("./zips/%s", p.PickupCode+".zip"))
	}
}

type s struct {
	Psd string `json:"psd"`
}

func Dao_checkadmin(c *fiber.Ctx) error {
	rsq := c.Body()
	var ss s
	json.Unmarshal(rsq, &ss)
	if ss.Psd != config.Config.AdminPassword {
		return c.JSON(fiber.Map{"title": "密码错误"})
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "admin"
	cookie.Value = config.Config.AdminToken
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.Path = "/"
	cookie.MaxAge = 86400
	cookie.HTTPOnly = true
	cookie.SameSite = "Lax"
	c.Cookie(cookie)
	return c.JSON(fiber.Map{"title": "登录成功"})
}
