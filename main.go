package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func main() {
	resp, err := http.Get("https://update.danmuji.me/getUpdate")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		updateResp := &UpdateResponse{}
		err := json.NewDecoder(resp.Body).Decode(updateResp)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return
		}

		err = downloadAndExtract(updateResp.Link)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Download and extraction successful!")
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}
	cmd := exec.Command("cmd.exe", "/C", "start", "GUI-BilibiliDanmuRobot.exe")
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("upgrade exit")
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
		return err
	}
	//zipFile, err := os.Create("temp.zip")
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	zipFile.Close()
	//	os.Remove("temp.zip")
	//}()
	//
	//_, err = io.Copy(bytes.NewReader(body), resp.Body)
	//if err != nil {
	//	return err
	//}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	//defer zipReader.Close()

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
