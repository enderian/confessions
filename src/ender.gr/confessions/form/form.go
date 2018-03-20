package form

import (
	"ender.gr/confessions/model"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strings"
	"github.com/tyler-sommer/stick"
	"fmt"
)

var ImageDirectory string

func CarrierForm(ctx *fasthttp.RequestCtx)  {
	path := string(ctx.Path())[1:]
	carrier, err := model.FindCarrier(path)
	if err != nil {
		ctx.SetBody([]byte("Η φορμα δεν υπάρχει!"))
		return
	}

	if carrier.Form.Enabled {
		RenderForm(ctx, carrier, "", "")
	} else {
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte("Η φορμα προσωρινά έχει απενεργοποιηθεί."))
		ctx.SetContentType("text/html; charset=utf-8")
	}
}

func RenderForm(ctx *fasthttp.RequestCtx, carrier model.Carrier, error string, success string)  {
	env := stick.New(nil)
	file, _ := ioutil.ReadFile("./templates/form.twig")
	ctx.SetContentType("text/html")

	values := map[string]stick.Value{
		"carrier": carrier.Id,
		"carrierFacebook": carrier.FacebookPage,
		"form": carrier.Form,
		"recaptcha": ReCaptchaSiteKey,
	}
	if error != "" {
		values["error"] = error
	}
	if success != "" {
		values["success"] = true
		values["secretId"] = success
	}

	env.Execute(string(file), ctx, values)
}

func FormCarrierCss(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())[1:]
	carrier, err := model.FindCarrier(path[0:strings.Index(path, "/")])
	if err != nil {
		ctx.SetBody([]byte("Η φορμα δεν υπάρχει!"))
		return
	}

	ctx.SetContentType("text/css")
	fmt.Fprint(ctx, carrier.Form.CustomCss)
}
