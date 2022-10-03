package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

var (
	SlackToken        string
	VerificationToken string
)

func init() {
	SlackToken = ""
}

func main() {
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)

	wPreviewImage, err := mw.CreateFormFile("preview_image", "41781157.jpeg")
	if err != nil {
		log.Printf(`Error: CreateFormFile, err=%s`, err)
		return
	}
	img, err := os.ReadFile("41781157.jpeg")
	if err != nil {
		log.Printf("can not read image file. err=%s", err)
		return
	}
	_, err = wPreviewImage.Write(img)
	if err != nil {
		log.Printf("can not write image file. err=%s", err)
	}

	externalId := uuid.NewString()
	_createFormField(mw, "external_id", externalId)
	_createFormField(mw, "external_url", "https://example.com")
	_createFormField(mw, "title", "slackunfurltest")
	_createFormField(mw, "indexable_file_contents", "search_terms.txt")

	url := "https://slack.com/api/files.remote.add"
	contentType := mw.FormDataContentType()
	mw.Close()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Printf("failed to create new http request. err=%s", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+SlackToken)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("failed to add remote file. err=%s", err)
		return
	}

	log.Printf(response.Status)

	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	var resBodyJson interface{}
	json.Unmarshal(resBody, &resBodyJson)
	log.Printf("Response: %+v", resBodyJson)
}

func _createFormField(mw *multipart.Writer, key string, value string) error {
	writer, err := mw.CreateFormField(key)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}
	_, err = writer.Write([]byte(value))
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}
	return nil
}

func addRemoteFile(*slack.RemoteFileParameters) error {

	return nil
}
