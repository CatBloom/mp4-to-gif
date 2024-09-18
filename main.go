package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// 引数でURLを指定する場合の処理
	// if len(os.Args) < 2 {
	// 	fmt.Println("Usage: go run main.go <input_url>")
	// 	return
	// }
	// inputURL := os.Args[1]

	inputURLs, err := textToSlice()
	if err != nil {
		fmt.Printf("error:%s", err.Error())
		return
	}
	if len(inputURLs) == 0 {
		fmt.Println("error:empty <input.txt>")
		return
	}

	outputDir := "/app/outputs"
	tempMP4 := "temp.mp4"

	// 出力ディレクトリが存在するか確認し、存在しない場合は作成
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("error:%s", err.Error())
		return
	}
	fmt.Printf("download start count:%d \n", len(inputURLs))
	for i, v := range inputURLs {
		fmt.Printf("%d downloading mp4...\n", i+1)
		err = downloadFile(v, tempMP4)
		if err != nil {
			fmt.Printf("error:%s", err.Error())
			return
		}

		// MP4ファイル名からGIFファイル名を生成
		baseName := strings.TrimSuffix(filepath.Base(v), filepath.Ext(v))
		outputFile := filepath.Join(outputDir, baseName+".gif")

		// すでに同名ファイルが存在している場合スキップ
		if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
			fmt.Printf("%s skipping...\n", baseName+".gif")
			os.Remove(tempMP4)
			continue
		}

		// FFmpegコマンド
		cmd := exec.Command("ffmpeg", "-i", tempMP4,
			"-filter_complex",
			"[0:v] fps=10,scale=640:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse",
			"-q:v", "0",
			"-crf", "18",
			outputFile)

		err = cmd.Run()
		if err != nil {
			fmt.Printf("error:%s", err.Error())
			return
		}
		os.Remove(tempMP4)
	}
	fmt.Println("gif created successfully")
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	return err
}

func textToSlice() ([]string, error) {
	result := []string{}

	file, err := os.Open("input.txt")
	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			result = append(result, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, err
}
