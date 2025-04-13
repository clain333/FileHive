package db

type File struct {
	HashValue   string `json:"hash_value"`      // 文件的哈希值
	NewFilename string `json:"source_filename"` // 新文件名
}

type PickupInfo struct {
	PickupCode         string     `json:"pickup_code"`         // 取件码
	PickupPassword     string     `json:"pickup_password"`     // 取件密码
	QrCode             string     `json:"qr_code"`             //分享的二维码
	IsDownloaded       bool       `json:"is_downloaded"`       //是否可以下载
	FileSize           int64      `json:" file_size"`          //文件的大小
	Text               string     `json:"text"`                //文件描述
	ExpirationTime     int64      `json:"expiration_time"`     //过期时间
	RemainingDownloads int        `json:"remaining_downloads"` //剩余下载数
	Index              int        `json:"index"`
	IsD                string     `json:"is_d"`
	FileInfo           []FileInfo `json:"file_info"`
}

type FileInfo struct {
	HashValue string `json:"hash_value"` // 文件的哈希值
	Filename  string `json:"filename"`   // 文件名
}

type Config struct {
	AdminPassword     string `json:"admin_password"`      // 管理员密码
	ExpirationTime    int64  `json:"expiration_time"`     // 文件过期时间（单位：天）
	SingleUploadLimit int    `json:"single_upload_limit"` // 单次上传最大文件数限制
	FileSizeLimit     int64  `json:"file_size_limit"`     // 上传总的文件大小限制（单位：字节）
	IsUpload          bool   `json:"is_upload"`           // 是否启用上传
	IsDownload        bool   `json:"is_download"`         // 是否启用下载
	DownloadLimit     int    `json:"download_limit"`      // 每个下载间隔内允许下载的最大次数
	PickupNum         int    `json:"pickup_num"`          //取件码几位数
	AdminToken        string `json:"admin_token"`
}

type HashName struct {
	Oldname string
	Newname string
}
