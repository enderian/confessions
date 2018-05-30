package form

import (
	"encoding/json"
	"github.com/enderian/confessions/database"
	"github.com/enderian/confessions/model"
	"github.com/valyala/fasthttp"
	"strings"
)

func StatusRead(ctx *fasthttp.RequestCtx) {

	carrierId := string(ctx.Path())[1 : strings.Index(string(ctx.Path())[1:], "/")+1]
	secret, err := database.FindSecret(ctx.UserValue("id").(string))
	if err != nil || carrierId != secret.Carrier {
		err, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: "Δεν βρέθηκε τέτοιο μυστικό.",
		})
		ctx.SetBody(err)
		ctx.SetStatusCode(404)
		return
	}

	carrier, err := database.FindCarrier(carrierId)
	statusProcess(secret, carrier, ctx)
}

func StatusPatch(ctx *fasthttp.RequestCtx) {
	carrierId := string(ctx.Path())[1 : strings.Index(string(ctx.Path())[1:], "/")+1]
	secret, err := database.FindSecret(ctx.UserValue("id").(string))
	if err != nil || carrierId != secret.Carrier {
		err, _ := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: "Δεν βρέθηκε τέτοιο μυστικό.",
		})
		ctx.SetBody(err)
		ctx.SetStatusCode(404)
		return
	}

	patch := struct {
		Action string `json:"action"`
	}{}
	json.Unmarshal(ctx.PostBody(), &patch)

	if patch.Action == "delete" && secret.Status == model.SENT {
		secret.Status = model.DELETED
	}

	carrier, err := database.FindCarrier(carrierId)
	statusProcess(secret, carrier, ctx)
}

func statusProcess(secret model.Secret, carrier model.Carrier, ctx *fasthttp.RequestCtx) {

	secret.ChecksData = append(secret.ChecksData, ConstructSourceData(ctx))
	database.SaveSecret(secret)

	content := secret.OriginalContent
	publishUrl := ""
	status := 0
	deletable := false

	if secret.Status == model.PUBLISHED {
		status = 1
		publishUrl = "https://www.facebook.com/" + carrier.FacebookPage + "/posts/" + secret.PublishData.FacebookPostId
	}
	if secret.Status == model.DELETED {
		status = 2
	}
	if secret.Status == model.SENT {
		deletable = true
	}
	if content == "" {
		content = secret.Content
	}

	ctx.Response.Header.Add("Content-Type", "application/json")
	bytes, _ := json.Marshal(struct {
		Id            string `json:"id"`
		Content       string `json:"content"`
		ContainsImage bool   `json:"containsImage"`
		PublishUrl    string `json:"publishUrl,omitempty"`
		Status        int    `json:"status"`
		Deletable     bool   `json:"deletable"`
	}{
		Id:            secret.Id,
		Content:       content,
		ContainsImage: secret.ImageId != "",
		PublishUrl:    publishUrl,
		Status:        status,
		Deletable:     deletable,
	})
	ctx.SetBody(bytes)
}
