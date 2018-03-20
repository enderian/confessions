package form

import (
	"ender.gr/confessions/model"
	"github.com/valyala/fasthttp"
	"github.com/google/uuid"
	"strings"
	"io/ioutil"
	"net"
	"time"
	"gopkg.in/h2non/filetype.v1"
)

func SecretSubmit(ctx *fasthttp.RequestCtx) {

	multForm, err := ctx.Request.MultipartForm()
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.Response.SetBody([]byte("Σφάλμα κατά την λήψη δεδομένων!"))
		return
	}

	carrier, err := model.FindCarrier(multForm.Value["carrier"][0])
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.Response.SetBody([]byte("Άκυρος πάροχος!"))
		return
	}

	if carrier.Form.IsEnableCaptcha {
		if len(multForm.Value["g-recaptcha-response"]) == 0 || !VerifyReCaptcha(multForm.Value["g-recaptcha-response"][0]) {
			RenderForm(ctx, carrier, "Δεν ολοκληρώσατε την πρόκληση Captcha!", "")
			return
		}
	}

	hasContent := len(multForm.Value["content"]) > 0 && multForm.Value["content"][0] != ""
	hasImage := len(multForm.File["image"]) > 0 && multForm.File["image"] != nil

	if !hasContent && !hasImage {
		RenderForm(ctx, carrier, "Το μυστικό σας πρέπει να περιέχει ή κείμενο ή εικόνα!", "")
		return
	}

	imageId := ""
	if hasImage {
		//web server is setup to filter large image requests

		uid, err := uuid.NewRandom()
		filename := strings.Replace(uid.String(), "-", "", -1) + "-" + multForm.File["image"][0].Filename
		file, err := multForm.File["image"][0].Open()

		if err != nil {
			RenderForm(ctx, carrier, "Σφάλμα κατά την λήψη δεδομένων!", "")
			return
		}

		bytes, err := ioutil.ReadAll(file)
		if !filetype.IsImage(bytes) {
			RenderForm(ctx, carrier, "Το αρχείο δεν ήταν έγκυρη εικόνα!", "")
			return
		}
		ioutil.WriteFile(ImageDirectory + filename, bytes, 0755)
		imageId = filename
	}

	uid, err := uuid.NewRandom()
	idRunes := []rune(strings.Replace(uid.String(), "-", "", -1))
	id := string(idRunes[0:12])
	source := ConstructSourceData(ctx)

	content := ""
	if hasContent {
		content = strings.Trim(multForm.Value["content"][0], " \n")
	}

	var options = make(map[string]string)
	for k, v := range carrier.Form.OptionSets {
		if len(multForm.Value["option-" + k]) > 0 && multForm.Value["option-" + k][0] != "" {

			input := multForm.Value["option-" + k][0]
			if input == "custom" && v.AllowCustom {
				if len(multForm.Value["option-" + k + "-custom"]) > 0 &&
					multForm.Value["option-" + k + "-custom"][0] != "" {
					options[k] = multForm.Value["option-"+k+"-custom"][0]
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
	secret.Save()
	RenderForm(ctx, carrier, "", secret.Id)
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
		Timestamp: makeTimestamp(),
		IpAddress: ip,
		Country: country,
		Hostname: strings.Trim(hostname, "."),
		RayID: string(ctx.Request.Header.Peek("CF-RAY")),
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}