package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"love-remittance-be-apps/lib/model"
	"net/http"
	"os"
	"strings"

	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/jlaffaye/ftp"
	"github.com/nfnt/resize"
)

func SignMD5(data string) string {
	datas := []byte(data)
	b := md5.Sum(datas)
	pass := hex.EncodeToString(b[:])
	// fmt.Println("MD5")
	// fmt.Printf("%s", pass)
	return pass
}

func SignSHA256(data string) string {
	datas := []byte(data)
	b := sha256.Sum256(datas)
	pass := hex.EncodeToString(b[:])
	// fmt.Println("SHA256")
	// fmt.Printf("%s", pass)
	return pass
}

func GenerateOtp() int {
	hi := 999999
	low := 100000
	return low + rand.Intn(hi-low)
}

func UploadImg(ctx *fiber.Ctx, nameFile string, uuID string) (*string, *model.ErrorData) {
	file, err := ctx.FormFile(nameFile)
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error get " + nameFile + err.Error(),
		}
	}

	var contentType = file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
		return nil, &model.ErrorData{
			Description: "Error type to " + nameFile,
		}
	}
	var imgType string
	if contentType == "image/jpeg" {
		imgType = ".jpeg"
	} else if contentType == "image/png" {
		imgType = ".png"
	} else {
		imgType = ".jpg"
	}

	var finishName string
	if nameFile == "img_signature" {
		finishName = "customer_signature_picture"
	} else if nameFile == "img_self" {
		finishName = "customer_self_picture"
	} else {
		finishName = "customer_self_identity_picture"
	}

	//make connection to ftp
	conns, err := ftp.Dial(os.Getenv("FTP_ADDRESS"))
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error connect to server ftp " + nameFile + err.Error(),
		}
	}

	err = conns.Login(os.Getenv("FTP_USER"), os.Getenv("FTP_PASSWORD"))
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error login to server ftp " + nameFile + err.Error(),
		}
	}

	resultName := uuID + "_" + finishName + "" + imgType
	sourceFile := fmt.Sprintf("/tmp/%s", resultName) //server
	// sourceFile := fmt.Sprintf("./temp/%s", resultName) //local
	if err := ctx.SaveFile(file, sourceFile); err != nil {
		// Handle error
		return nil, &model.ErrorData{
			Description: "Error store to local " + nameFile + err.Error(),
		}
	}

	f, err := os.Open(sourceFile)
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error get file from store local " + nameFile + err.Error(),
		}
	}
	destinatiFile := "." + os.Getenv("FTP_PATH") + "/" + resultName
	errssss := conns.Stor(destinatiFile, f)
	if errssss != nil {
		return nil, &model.ErrorData{
			Description: "Error store file to server " + nameFile + errssss.Error(),
		}
	}
	defer f.Close()

	log.Print("Success upload file = " + nameFile + " with = " + resultName)
	errDel := os.Remove(sourceFile)
	if errDel != nil {
		return nil, &model.ErrorData{
			Description: "Error delete file from local " + nameFile + errDel.Error(),
		}
	}
	log.Print("Delete file from temp = " + nameFile + " with = " + resultName)

	return &resultName, nil
}

func ShowImg(nameImg string) (*string, *model.ErrorData) {
	conns, err := ftp.Dial(os.Getenv("FTP_ADDRESS"))
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error connect to server ftp " + err.Error(),
		}
	}

	err = conns.Login(os.Getenv("FTP_USER"), os.Getenv("FTP_PASSWORD"))
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error login to server ftp " + err.Error(),
		}
	}
	destinatiFile := "." + os.Getenv("FTP_PATH") + "/" + nameImg
	files, err := conns.Retr(destinatiFile)
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error get file from ftp " + err.Error(),
		}
	}

	imgByte, err := io.ReadAll(files)
	files.Close()
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error get file from ftp " + err.Error(),
		}
	}

	fmt.Println(" ->", files, "->", string(imgByte))
	fff := string(imgByte)
	return &fff, nil
}

func GetImgString(imgName string) (string, *model.ErrorData) {
	conn, err := ftp.Dial(os.Getenv("FTP_ADDRESS"))
	if err != nil {
		return "", &model.ErrorData{
			Description: "Error connect to server ftp " + err.Error(),
		}
	}

	err = conn.Login(os.Getenv("FTP_USER"), os.Getenv("FTP_PASSWORD"))
	if err != nil {
		return "", &model.ErrorData{
			Description: "Error login to server ftp " + err.Error(),
		}
	}
	desFile := "." + os.Getenv("FTP_PATH") + "/" + imgName
	files, err := conn.Retr(desFile)
	if err != nil {
		log.Println(err.Error())
	}
	defer files.Close()

	imgByte, err := io.ReadAll(files)
	defer files.Close()
	if err != nil {
		log.Println(err.Error())
	}
	contentType := http.DetectContentType(imgByte)

	switch contentType {
	case "image/png":
		img, err := png.Decode(bytes.NewReader(imgByte))
		if err != nil {
			return "", &model.ErrorData{
				Description: "unable to decode png: " + err.Error(),
			}
		}
		newImage := resize.Resize(100, 0, img, resize.Lanczos3)
		var buf bytes.Buffer
		if err := png.Encode(&buf, newImage); err != nil {
			return "", &model.ErrorData{
				Description: "unable to encode png: " + err.Error(),
			}
		}
		imgByte = buf.Bytes()
		log.Println("done resizing img .png")
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imgByte))
		if err != nil {
			return "", &model.ErrorData{
				Description: "unable to decode jpeg: " + err.Error(),
			}
		}
		newImage := resize.Resize(100, 0, img, resize.Lanczos3)

		var buf bytes.Buffer
		if err := png.Encode(&buf, newImage); err != nil {
			return "", &model.ErrorData{
				Description: "unable to encode png: " + err.Error(),
			}
		}
		imgByte = buf.Bytes()
	default:
		return "", &model.ErrorData{
			Description: "unsupported content typo: " + err.Error(),
		}
	}
	imgBase64Str := base64.StdEncoding.EncodeToString(imgByte)

	return imgBase64Str, nil
}

func GetFileConfig(name string, lang string) string {
	result := os.Getenv("FORM_CONFIG") + "/" + name + "_" + strings.ToUpper(lang) + os.Getenv("FORM_TYPE")
	log.Println("result " + result)
	return result
}

func GetParameterNotes(name string, lang string) (map[string]interface{}, *model.ErrorData) {
	nameFile := os.Getenv("FORM_CONFIG") + "/" + name + "_" + strings.ToUpper(lang) + os.Getenv("FORM_TYPE")
	log.Println("result " + nameFile)

	content, err := os.ReadFile(nameFile)
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Form tidak tersedia: " + err.Error(),
		}
	}
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Error during Unmarshal(): " + err.Error(),
		}
	}
	return payload, nil
}
