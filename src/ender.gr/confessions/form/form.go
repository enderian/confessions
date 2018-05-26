package form

import (
	"ender.gr/confessions/model"
	"github.com/valyala/fasthttp"
	"html/template"
)

var ImageDirectory string
var formTemplate *template.Template

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

func SetupForm() {
	formTemplate = template.Must(template.New("template.html"), nil)
	if _, err := formTemplate.ParseFiles("templates/form.html", "templates/template.html"); err != nil {
		panic(err)
	}
}

func RenderForm(ctx *fasthttp.RequestCtx, carrier model.Carrier, error string, success interface{})  {

	customStyle := "body{ background: url('" + carrier.Form.BackgroundUrl + "') center; " +
		"background-size: cover; background-repeat: no-repeat; } " +
		".jumbotron {color: " + carrier.Form.TitleColor + ";} "
	if carrier.Form.AccentColor != "" {
		customStyle += ".form-jumbotron {background: rgba(" + carrier.Form.AccentColor + ", 0.60);}"
	}
	if carrier.Form.CustomCss != "" {
		customStyle += carrier.Form.CustomCss
	}

	ctx.SetContentType("text/html")
	if err := formTemplate.Execute(ctx, map[string]interface{}{
		"Carrier": carrier,
		"Title": carrier.Name,
		"Icon": "https://graph.facebook.com/" + carrier.FacebookPage + "/picture?type=square",
		"RecaptchaKey": ReCaptchaSiteKey,
		"CustomStyle": template.CSS(customStyle),
	}); err != nil {
		panic(err)
	}
}