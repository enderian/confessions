package form

import (
	"github.com/enderian/confessions/database"
	"github.com/enderian/confessions/model"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/valyala/fasthttp"
	"html/template"
	"strconv"
)

var ImageDirectory string
var ServiceAlert string

var formTemplate *template.Template

func CarrierForm(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())[1:]
	carrier, err := database.FindCarrier(path)
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

func RenderForm(ctx *fasthttp.RequestCtx, carrier model.Carrier, error string, success interface{}) {

	barColor := ""
	recaptchaKey := ""
	customStyle := "body{ background: url('" + carrier.Form.BackgroundUrl + "') center; " +
		"background-size: cover; background-repeat: no-repeat; } " +
		".jumbotron {color: " + carrier.Form.TitleColor + ";} "

	if carrier.Form.AccentColor != "" {
		barColor = carrier.Form.AccentColor
		c, err := colorful.Hex(carrier.Form.AccentColor)
		if err == nil {
			customStyle += ".form-jumbotron {background: rgba(" +
				strconv.Itoa(int(c.R*255)) + "," +
				strconv.Itoa(int(c.G*255)) + "," +
				strconv.Itoa(int(c.B*255)) + "," +
				" 0.60) !important;}"
		}
	}
	if carrier.Form.CustomCss != "" {
		customStyle += carrier.Form.CustomCss
	}
	if carrier.Form.IsEnableCaptcha {
		recaptchaKey = ReCaptchaSiteKey
	}

	ctx.SetContentType("text/html")
	if err := formTemplate.Execute(ctx, map[string]interface{}{
		"Title":        carrier.Name,
		"Carrier":      carrier,
		"BarColor":     barColor,
		"RecaptchaKey": recaptchaKey,
		"Icon":         "https://graph.facebook.com/" + carrier.FacebookPage + "/picture?type=square",
		"ServiceAlert":	template.HTML(ServiceAlert),
		"CustomStyle":  template.CSS(customStyle),
	}); err != nil {
		panic(err)
	}
}
