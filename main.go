package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type UpdateResponse struct {
	Version   string `json:"version"`
	Link      string `json:"link"`
	ChangeLog string `json:"changeLog"`
}

var Version string

func main() {
	log.Printf("")
	log.Println("正在查询版本信息")
	resp, err := http.Get("https://danmuji.neuedu.work/getUpdate")
	if err != nil {
		log.Println("连接版本服务器错误")
		log.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		updateResp := &UpdateResponse{}
		err := json.NewDecoder(resp.Body).Decode(updateResp)
		if err != nil {
			log.Println("版本信息解析失败")
			log.Println("Error decoding JSON response:", err)
			return
		}
		log.Println("正在更新弹幕机")
		err = downloadAndExtract(updateResp.Link)
		if err != nil {
			log.Println("更新弹幕机失败")
			log.Println("Error:", err)
			return
		}
		log.Println("弹幕机更新完成即将启动")
		cmd := exec.Command("cmd.exe", "/C", "start", "GUI-BilibiliDanmuRobot.exe")
		if err := cmd.Start(); err != nil {
			log.Println("启动弹幕机失败，请手动启动")
			log.Println(err)
		}
		log.Println("更新完成即将退出更新程序")
	} else {
		log.Println("更新服务器链接失败")
		log.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}

	time.Sleep(10 * time.Second)
	log.Println("upgrade exit")
	os.Exit(0)
}

func downloadAndExtract(link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("弹幕机下载失败")
		return err
	}
	log.Println("弹幕机下载成功，正在解压")
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	log.Println("解压成功，正在更新软件")
	for _, file := range zipReader.File {
		if filepath.Base(file.Name) == "GUI-BilibiliDanmuRobot.exe" {
			zippedFile, err := file.Open()
			if err != nil {
				return err
			}
			defer zippedFile.Close()

			extractedFile, err := os.OpenFile(filepath.Base(file.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer extractedFile.Close()

			_, err = io.Copy(extractedFile, zippedFile)
			if err != nil {
				return err
			}
			break // 只提取第一个匹配的文件
		}
	}

	return nil
}
