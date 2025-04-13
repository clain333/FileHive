package utils

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"filepro/db"
	"filepro/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"[rand.Intn(len("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))]
	}
	return string(result)
}

func Returninfo(code, title, text, password, codepass string) fiber.Map {
	return fiber.Map{
		"code":     code,
		"title":    title,
		"text":     text,
		"password": password,
		"codepass": codepass,
	}

}

func CreateHa(filename string) string {
	file, err := os.Open(fmt.Sprintf("%s/%s", model.Dir, filename)) // 替换为你的文件路径
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := sha256.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatal(err)
	}

	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func RemoveFile(string2 string) {
	os.RemoveAll(fmt.Sprintf("%s/%s", model.Dir, string2))
}

func CreateQr(url, dir string) {
	err := qrcode.WriteFile(url, qrcode.Medium, 256, dir)
	if err != nil {
		panic(err)
	}
}

func CreateZip(zipdir string, files []db.HashName) {
	zipPath := fmt.Sprintf("%s/%s.zip", model.Zipdir, filepath.Base(zipdir))
	zipFile, err := os.Create(zipPath)
	if err != nil {
		fmt.Println("无法创建压缩包:", err)
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, f := range files {
		filePath := fmt.Sprintf("%s/%s", model.Dir, f.Newname)
		fileToZip, err := os.Open(filePath)
		if err != nil {
			fmt.Println("无法打开文件:", filePath, err)
			continue
		}
		defer fileToZip.Close()

		info, err := fileToZip.Stat()
		if err != nil {
			fmt.Println("无法获取文件信息:", filePath, err)
			continue
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			fmt.Println("创建 zip 文件头失败:", err)
			continue
		}
		header.Name = f.Oldname
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			fmt.Println("创建 zip writer 失败:", err)
			continue
		}

		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			fmt.Println("写入 zip 文件失败:", filePath, err)
			continue
		}
	}
	err = zipWriter.Close()
	if err != nil {
		fmt.Println("关闭 zip writer 时出错:", err)
		return
	}
}

func Createdir(s string) {
	folderPath := fmt.Sprintf("./%s", s) // 你想创建的文件夹路径

	// 检查文件夹是否存在
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 文件夹不存在，创建它
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败:", err)
			return
		}
	} else {
	}
}
