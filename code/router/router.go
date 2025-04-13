package router

import (
	"filepro/dao"
	"filepro/midd"
	"github.com/gofiber/fiber/v2"
)

func Router_page(router fiber.Router) {
	router.Get("/", dao.DaoUpload)
	router.Get("/getcode", dao.DaoGet)
	router.Get("/admin", midd.Midd_admin, dao.DaoLogin)
	router.Get("/setting", midd.Midd_admin, dao.DaoSetting)
	router.Get("/fileinfo", midd.Midd_check, dao.DaoFile)
}

func Router_api(router fiber.Router) {
	router.Post("/upload", midd.Midd_upload, dao.Dao_uploadfile)
	router.Post("/getcode", dao.Dao_getcode)
	router.Post("/download", midd.Midd_download, dao.Dao_download)
	router.Post("/admin", dao.Dao_checkadmin)
}

func Router_admin(router fiber.Router) {
	router.Post("/change", midd.Midd_role, dao.Dao_change)
	router.Post("/delete", midd.Midd_role, dao.Dao_delete)
	router.Post("/changeall", midd.Midd_role, dao.Dao_changeall)
	router.Post("/deleteall", midd.Midd_role, dao.Dao_deleteall)
	router.Post("/savesetting", midd.Midd_role, dao.Dao_savesetting)
	router.Post("/recover", midd.Midd_role, dao.Dao_recover)
}
