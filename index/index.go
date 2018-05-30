package index

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"html/template"
)

var indexTemplate *template.Template
var helpTemplate *template.Template
var privacyTemplate *template.Template
var joinTemplate *template.Template

func RegisterIndex(router *fasthttprouter.Router) {

	indexTemplate = template.Must(template.New("template.html"), nil)
	helpTemplate = template.Must(template.New("template.html"), nil)
	privacyTemplate = template.Must(template.New("template.html"), nil)
	joinTemplate = template.Must(template.New("template.html"), nil)
	indexTemplate.ParseFiles("templates/home.html", "templates/index_template.html", "templates/template.html")
	helpTemplate.ParseFiles("templates/help.html", "templates/index_template.html", "templates/template.html")
	privacyTemplate.ParseFiles("templates/privacy.html", "templates/index_template.html", "templates/template.html")
	joinTemplate.ParseFiles("templates/join.html", "templates/index_template.html", "templates/template.html")

	router.GET("/", renderIndex)
	router.GET("/help", renderHelp)
	router.GET("/privacy", renderPrivacy)
	router.GET("/join", renderJoin)
}

func renderIndex(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	if err := indexTemplate.Execute(ctx, map[string]interface{}{"HomePage": "active"}); err != nil {
		panic(err)
	}
}

func renderHelp(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	if err := helpTemplate.Execute(ctx, map[string]interface{}{"Help": "active"}); err != nil {
		panic(err)
	}
}

func renderPrivacy(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	if err := privacyTemplate.Execute(ctx, map[string]interface{}{"Privacy": "active"}); err != nil {
		panic(err)
	}
}

func renderJoin(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	if err := joinTemplate.Execute(ctx, map[string]interface{}{"Join": "active"}); err != nil {
		panic(err)
	}
}
