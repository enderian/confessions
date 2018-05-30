package form

import (
	"encoding/json"
	"github.com/enderian/confessions/database"
	"github.com/enderian/confessions/model"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"gopkg.in/h2non/filetype.v1"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

func SecretSubmit(ctx *fasthttp.RequestCtx) {

	multiPartForm, err := ctx.Request.MultipartForm()
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(returnError("Σφάλμα κατά την λήψη δεδομένων!"))
		return
	}

	carrierId := string(ctx.Path())[1 : strings.Index(string(ctx.Path())[1:], "/")+1]
	carrier, err := database.FindCarrier(carrierId)
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.SetBody(returnError("Άκυρος πάροχος!"))
		return
	}

	formData := struct {
		Content string            `json:"content"`
		Options map[string]string `json:"options"`
		Captcha string            `json:"captcha"`
	}{}

	json.Unmarshal([]byte(multiPartForm.Value["secret"][0]), &formData)

	if carrier.Form.IsEnableCaptcha {
		if !VerifyReCaptcha(formData.Captcha) {
			ctx.SetStatusCode(400)
			ctx.SetBody(returnError("Δεν ολοκληρώσατε την πρόκληση Captcha!"))
			return
		}
	}

	hasContent := formData.Content != ""
	hasImage := len(multiPartForm.File["file"]) > 0 && multiPartForm.File["file"] != nil

	if !hasContent && !hasImage {
		ctx.SetBody(returnError("Το μυστικό σας πρέπει να περιέχει κείμενο ή εικόνα!"))
		ctx.SetStatusCode(400)
		return
	}

	imageId := ""
	if hasImage {

		uid, err := uuid.NewRandom()
		filename := strings.Replace(uid.String(), "-", "", -1) + "-" + multiPartForm.File["file"][0].Filename
		file, err := multiPartForm.File["file"][0].Open()

		if err != nil {
			ctx.SetStatusCode(400)
			ctx.SetBody(returnError("Σφάλμα κατά την λήψη δεδομένων!"))
			return
		}

		bytes, err := ioutil.ReadAll(file)
		if !filetype.IsImage(bytes) {
			ctx.SetStatusCode(400)
			ctx.SetBody(returnError("Το αρχείο δεν ήταν έγκυρη εικόνα!"))
			return
		}
		ioutil.WriteFile(ImageDirectory+filename, bytes, 0755)
		imageId = filename
	}

	uid, err := uuid.NewRandom()
	idRunes := []rune(strings.Replace(uid.String(), "-", "", -1))
	id := string(idRunes[0:12])
	source := ConstructSourceData(ctx)

	content := ""
	if hasContent {
		content = strings.Trim(formData.Content, " \n")
	}

	var options = make(map[string]string)
	for k, v := range carrier.Form.OptionSets {
		input, has1 := formData.Options["option-"+k]
		has1 = has1 && len(input) > 0
		customInput, has2 := formData.Options["option-"+k+"-custom"]
		has2 = has2 && len(customInput) > 0

		if has1 || has2 {
			if input == "custom" && v.AllowCustom {
				if has2 {
					options[k] = customInput
				}
			}
			if contains(v.Options, input) {
				value, ok := v.OptionDisplay[input]
				if ok {
					options[k] = value
				} else {
					options[k] = input
				}
			}
		}
	}

	secret := model.Secret{
		Carrier:    carrier.Id,
		Id:         id,
		Status:     model.SENT,
		Content:    content,
		SourceData: source,
		ImageId:    imageId,
		Options:    options,
	}
	database.SaveSecret(secret)

	ctx.Write(func() []byte {
		js, _ := json.Marshal(struct {
			Id string `json:"id"`
		}{
			Id: id,
		})
		return js
	}())
}

func ConstructSourceData(ctx *fasthttp.RequestCtx) model.SecretSourceData {
	ip := string(ctx.Request.Header.Peek("CF-Connecting-IP"))
	country := string(ctx.Request.Header.Peek("CF-IPCountry"))
	addr, err := net.LookupAddr(ip)
	hostname := "?"
	if err == nil {
		hostname = addr[0]
	}

	return model.SecretSourceData{
		Timestamp: time.Now(),
		IpAddress: ip,
		Country:   country,
		Hostname:  strings.Trim(hostname, "."),
		UserAgent: string(ctx.Request.Header.Peek("User-Agent")),
		RayID:     string(ctx.Request.Header.Peek("CF-RAY")),
	}
}

func returnError(error string) []byte {
	js, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: error,
	})
	return js
	return js
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
