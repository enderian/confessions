package form

import (
	"github.com/enderian/confessions/model"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"github.com/tyler-sommer/stick"
	"github.com/enderian/confessions/database"
)

func StatusRead(ctx *fasthttp.RequestCtx) {

	request := ctx.Request.PostArgs()
	secret, err := database.FindSecret(string(request.Peek("id")))
	if err != nil || string(request.Peek("carrier")) != secret.Carrier {
		StatusRender(ctx, model.Secret{})
		return
	}
	if string(request.Peek("action")) == "delete" {
		secret.Status = model.DELETED
	}

	secret.ChecksData = append(secret.ChecksData, ConstructSourceData(ctx))
	database.SaveSecret(secret)

	StatusRender(ctx, secret)
}

func StatusRender(ctx *fasthttp.RequestCtx, secret model.Secret)  {
	env := stick.New(nil)
	file, _ := ioutil.ReadFile("./templates/form_status.twig")
	ctx.SetContentType("text/html")

	published := ""
	if secret.Status == model.PUBLISHED {
		carrier, err := database.FindCarrier(secret.Carrier)
		if err == nil {
			published = "https://www.facebook.com/" + carrier.FacebookPage + "/posts/" + secret.PublishData.FacebookPostId
		}
	}

	values := map[string]stick.Value{
		"secret": secret,
		"deletable": secret.Status == model.SENT,
		"published": published,
	}
	env.Execute(string(file), ctx, values)
}
