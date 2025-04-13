package config

import (
	"filepro/db"
	"filepro/utils"
	"sync"
)

var (
	Config     db.Config
	ConfigLock = sync.RWMutex{}
)

func InitConfig() {
	c, e := db.ConfigGet()
	if e != nil {
		config := db.Config{
			AdminPassword:     "123456",
			ExpirationTime:    7,
			SingleUploadLimit: 200,
			FileSizeLimit:     4096,
			DownloadLimit:     20,
			IsUpload:          true,
			IsDownload:        true,
			PickupNum:         3,
			AdminToken:        utils.GenerateRandomString(12),
		}
		config.ConfigSave()
		Config = config
	} else {
		Config = *c
	}
}
