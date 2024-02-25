package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

func main() {
	client := http.Client{}
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	get_boundary := writer.Boundary()
	_ = writer.WriteField("is_multiple", "1")
	_ = writer.WriteField("model", "anime_model_lovelive")

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "1.jpg"))
	h.Set("Content-Type", "image/png")
	file_parameter, _ := writer.CreatePart(h)
	file, err := os.ReadFile("demo.png")
	if err != nil {
		fmt.Println("画像の読み込みに失敗しました！")
		return
	}
	fmt.Println(writer.FormDataContentType())

	defer writer.Close()
	file_parameter.Write(file)
	buffer.Write([]byte("\r\n" + "\r\n" + "--" + get_boundary + "--\r\n"))

	apiUrl := "https://aiapiv2.animedb.cn/ai/api/detect"

	req, err := http.NewRequest("POST", apiUrl, buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return
	}

	content, err := client.Do(req)

	defer content.Body.Close()
	all, err := io.ReadAll(content.Body)
	if err != nil {
		return
	}
	fmt.Println(string(all))
	fmt.Println("識別終了！")
}
