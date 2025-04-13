package midd

import (
	"encoding/json"
	"filepro/config"
	"filepro/db"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Midd_admin(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	if cookie == config.Config.AdminToken {
		if c.Path() == "/setting" {
			return c.Next()
		} else {
			return c.Redirect("/setting")
		}
	} else {
		if c.Path() == "/admin" {
			return c.Next()
		} else {
			return c.Redirect("/admin")
		}
	}
}

func Midd_check(c *fiber.Ctx) error {
	code := c.Query("code")
	psd := c.Query("psd")
	p, e := db.PickupGet(code)
	if e != nil {
		return c.Send([]byte("参数错误"))
	} else if psd != p.PickupPassword {
		return c.Send([]byte("参数错误"))
	} else {
		return c.Next()
	}
}

func Midd_role(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	if cookie == config.Config.AdminToken {
		return c.Next()
	} else {
		return c.Send([]byte("权限不够"))
	}
}

func Midd_upload(c *fiber.Ctx) error {
	//cookie := c.Cookies("admin")
	//if cookie == config.Config.AdminToken {
	//	return c.Next()
	//}
	form, _ := c.MultipartForm()
	files := form.File["files"]
	var filesize int64 = 0
	for _, file := range files {
		filesize += file.Size / 1024 / 1024
	}
	if len(files) > config.Config.SingleUploadLimit {
		return c.JSON(fiber.Map{"title": "上传错误"})
	}
	if config.Config.IsUpload == false {
		return c.JSON(fiber.Map{"title": "上传错误"})
	}
	if filesize > config.Config.FileSizeLimit {
		return c.JSON(fiber.Map{"title": "上传错误"})

	}
	return c.Next()
}

type f struct {
	Pickupcode string `json:"code"`
}

func Midd_download(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")

	if cookie == config.Config.AdminToken {
		return c.Next()
	}
	var fs f
	body := c.Body()
	json.Unmarshal(body, &fs)
	d, _ := db.PickupGet(fs.Pickupcode)

	if d.IsDownloaded && config.Config.IsDownload && d.RemainingDownloads != 0 && time.Now().Unix() < d.ExpirationTime {
		d.RemainingDownloads -= 1
		db.PickupSets(d.PickupCode, d)
		return c.Next()
	} else {
		return c.SendStatus(fiber.StatusForbidden)
	}

}
