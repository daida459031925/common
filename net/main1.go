package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	//chmod +x /home/pi/smpgo && echo "*/5 * * * * /home/pi/smpgo" | sudo crontab -
	url := "http://192.168.2.179:8089/oauth/token?grant_type=client_credentials"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}

	// 添加基本身份验证
	username := "dlkfdlf94r3409093dj0f9eur09ur03r"
	password := "poiud39876sffferete10perwere6y5t"
	auth := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var accessTokenResponse AccessTokenResponse
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Access Token: %s\n", accessTokenResponse.AccessToken)

	// 使用 access token 获取文件流
	fileURL := "http://192.168.2.179:8089/api/v1.0/customer/1/ecndb?access_token=" + accessTokenResponse.AccessToken
	resp, err = http.Get(fileURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 将文件流写入本地文件
	outFile, err := os.Create("output.zip")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("File written successfully")

	// 解压缩ZIP文件
	//err = unzip("output.zip", "/boot/gw")
	err = unzip("output.zip", ".")
	if err != nil {
		panic(err)
	}

	fmt.Println("ZIP file unzipped successfully")

	// 将解压后的文件内容写入到 /boot/gw/ecn.db
	err = writeFileContent("qwert12345yu.db", "/boot/gw/ecn.db")
	//err = writeFileContent("qwert12345yu.db", "ecn.db")
	if err != nil {
		panic(err)
	}

	fmt.Println("File content written successfully")

	// 在 /tmp 文件夹下查找所有的 .db 文件
	err = findAndReplaceDBFiles("/tmp")
	if err != nil {
		panic(err)
	}

	fmt.Println("File content written successfully")
}

func unzip(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}

			rc, err := file.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}

			outFile.Close()
			rc.Close()
		}
	}
	return nil
}

func writeFileContent(sourceFile, destFile string) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func findAndReplaceDBFiles(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录，则跳过
		if d.IsDir() {
			return nil
		}

		// 如果是 .db 文件
		if strings.HasSuffix(strings.ToLower(d.Name()), ".db") {
			// 使用 writeFileContent 方法将 qwert12345yu.db 文件的内容覆盖到当前 .db 文件中
			err = writeFileContent("qwert12345yu.db", path)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
