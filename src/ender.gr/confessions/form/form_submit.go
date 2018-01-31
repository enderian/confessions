package form

import (
	"ender.gr/confessions/model"
	"github.com/valyala/fasthttp"
	"github.com/google/uuid"
	"strings"
	"io/ioutil"
	"net"
	"time"
)

func SecretSubmit(ctx *fasthttp.RequestCtx) {

	mult_form, err := ctx.Request.MultipartForm()
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.Response.SetBody([]byte("Σφάλμα κατά την λήψη δεδομένων!"))
		return
	}

	carrier, err := model.FindCarrier(mult_form.Value["carrier"][0])
	if err != nil {
		ctx.SetStatusCode(400)
		ctx.Response.SetBody([]byte("Άκυρος πάροχος!"))
		return
	}

	if carrier.Form.IsEnableCaptcha {
		if len(mult_form.Value["g-recaptcha-response"]) == 0 || !VerifyReCaptcha(mult_form.Value["g-recaptcha-response"][0]) {
			RenderForm(ctx, carrier, "Δεν ολοκληρώσατε την πρόκληση Captcha!", "")
			return
		}
	}

	has_content := len(mult_form.Value["content"]) > 0 && mult_form.Value["content"][0] != ""
	has_image := len(mult_form.File["image"]) > 0 && mult_form.File["image"] != nil

	if !has_content && !has_image {
		RenderForm(ctx, carrier, "Το μυστικό σας πρέπει να περιέχει ή κείμενο ή εικόνα!", "")
		return
	}

	image_id := ""
	if has_image {
		//web server is setup to filter large image requests

		uid, err := uuid.NewRandom()
		filename := strings.Replace(uid.String(), "-", "", -1) + "-" + mult_form.File["image"][0].Filename
		file, err := mult_form.File["image"][0].Open()

		if err != nil {
			RenderForm(ctx, carrier, "Σφάλμα κατά την λήψη δεδομένων!", "")
			return
		}

		bytes, err := ioutil.ReadAll(file)
		ioutil.WriteFile(ImageDirectory + filename, bytes, 0755)
		image_id = filename
	}

	uid, err := uuid.NewRandom()
	id_runes := []rune(strings.Replace(uid.String(), "-", "", -1))
	id := string(id_runes[0:12])
	source := ConstructSourceData(ctx)

	content := ""
	if has_content {
		content = strings.Trim(mult_form.Value["content"][0], " \n")
	}

	var options map[string]string = make(map[string]string)
	for k, v := range carrier.Form.OptionSets {
		if len(mult_form.Value["option-" + k]) > 0 && mult_form.Value["option-" + k][0] != "" {

			input := mult_form.Value["option-" + k][0]
			if input == "custom" && v.AllowCustom {
				if len(mult_form.Value["option-" + k + "-custom"]) > 0 &&
					mult_form.Value["option-" + k + "-custom"][0] != "" {
					options[k] = mult_form.Value["option-"+k+"-custom"][0]
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
		ImageId:    image_id,
		Options: 	options,
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