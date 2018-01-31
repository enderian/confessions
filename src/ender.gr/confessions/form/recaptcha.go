package form

import (
	"github.com/valyala/fasthttp"
	"net/url"
	"encoding/json"
)

var ReCaptchaSiteKey string
var ReCaptchaSiteSecret string

type VerifyReCaptchaResponse struct {
	Success bool `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

func VerifyReCaptcha(grecaptcharesponse string) bool {
	client :=&fasthttp.Client{}
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	data := url.Values{}
	data.Set("response", grecaptcharesponse)
	data.Set("secret", ReCaptchaSiteSecret)

	request.URI().Update("https://www.google.com/recaptcha/api/siteverify")
	request.Header.SetMethodBytes([]byte("POST"))
	request.SetBody([]byte(data.Encode()))

	resp := VerifyReCaptchaResponse{}
	err := client.Do(request, response)
	if err != nil {
		return false
	}
	json.Unmarshal(response.Body(), &resp)
	return resp.Success
}
