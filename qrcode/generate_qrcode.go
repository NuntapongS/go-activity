package qrcode

import (
	"encoding/base64"
	"encoding/json"
	"log"

	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQrCode(activity any) (string, error) {
	var base64Image string
	data, err := json.Marshal(&activity)
	if err != nil {
		return "", err
	}

	code, err := qrcode.Encode(string(data), qrcode.Medium, 2048)
	if err != nil {
		log.Fatal(err)
	}

	base64Image = base64.StdEncoding.EncodeToString(code)
	base64Image = "data:image/png;base64," + base64Image
	return base64Image, nil
}
