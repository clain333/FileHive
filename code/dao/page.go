package dao

import (
	"filepro/config"
	"filepro/db"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func DaoUpload(c *fiber.Ctx) error {
	return c.Render("upload", nil)
}

func DaoGet(c *fiber.Ctx) error {
	return c.Render("get", nil)

}
func DaoLogin(c *fiber.Ctx) error {
	return c.Render("login", nil)
}
func DaoSetting(c *fiber.Ctx) error {
	fff, _ := db.PickupGetAll()
	for i, pickup := range fff {
		fff[i].FileSize = int64(float64(pickup.FileSize) / 1024 / 1024)
		if pickup.IsDownloaded {
			fff[i].IsD = "允许"
		} else {
			fff[i].IsD = "禁止"
		}
		if len(pickup.FileInfo) > 1 {
			fff[i].Text = "file.zip"
		} else {
			fff[i].Text = pickup.FileInfo[0].Filename
		}
		fff[i].Index = i + 2
	}
	return c.Render("setting", fiber.Map{
		"adminpassword": config.Config.AdminPassword,
		"exptime":       config.Config.ExpirationTime,
		"uploadlimbt":   config.Config.SingleUploadLimit,
		"filesize":      config.Config.FileSizeLimit,
		"DownloadLimit": config.Config.DownloadLimit,
		"isdownload":    config.Config.IsDownload,
		"isupload":      config.Config.IsUpload,
		"filecodenum":   config.Config.PickupNum,
		"pickinfo":      fff,
	})
}
func DaoFile(c *fiber.Ctx) error {
	code := c.Query("code")
	p, _ := db.PickupGet(code)
	size := p.FileSize / 1024 / 1024
	url := c.BaseURL() + "/" + p.QrCode
	filename := p.FileInfo[0].Filename

	if len(p.FileInfo) > 1 {
		for _, n := range p.FileInfo {
			filename = fmt.Sprintf("%s,%s", filename, n.Filename)
		}
	}
	num := len(p.FileInfo)
	return c.Render("file", fiber.Map{
		"code":     p.PickupCode,
		"psd":      p.PickupPassword,
		"isdown":   p.IsDownloaded,
		"filesize": size,
		"exptime":  p.ExpirationTime,
		"text":     p.Text,
		"qruel":    url,
		"name":     filename,
		"filenum":  num,
	})
}
