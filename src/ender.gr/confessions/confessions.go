package main

import (
	"ender.gr/confessions/form"
	"ender.gr/confessions/model"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"gopkg.in/mgo.v2"
	"strings"
	"github.com/tyler-sommer/stick"
	"io/ioutil"
	"encoding/json"
	"time"
)

type Configuration struct {
	Port string `json:"port"`
	ConfessionsImages string `json:"confessions_images"`

	ReCaptchaSiteKey string `json:"recaptcha_key"`
	ReCaptchaSiteSecret string `json:"recaptcha_secret"`
}

func main() {
	config := Configuration{}
	config_file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Unable to open configuration file config.json")
	}
	err = json.Unmarshal(config_file, &config)
	if err != nil {
		panic("Unable to open configuration file config.json")
	}

	connect()
	router := fasthttprouter.New()

	form.ReCaptchaSiteKey = config.ReCaptchaSiteKey
	form.ReCaptchaSiteSecret = config.ReCaptchaSiteSecret
	form.ImageDirectory = config.ConfessionsImages
	router.POST("/secret", form.StatusRead)
	router.POST("/submit", form.SecretSubmit)

	//Register carriers
	go registerCarriers(router)

	router.GET("/", indexHandler)
    router.GET("/admin", adminHandler)

	router.NotFound = notFoundHandler
	fsHandler = fasthttp.FSHandler("./frontend", 0)
	start(router, config.Port)
}


func registerCarriers(router *fasthttprouter.Router) {

	registered := []string{}
	for {
		for _, k := range model.FindCarriers() {
			for _, b := range registered {
				if b == k.Id {
					goto Skip
				}
			}
			registered = append(registered, k.Id)
			router.GET("/" + k.Id, form.FormCarrier)
			router.GET("/" + k.Id + "/css", form.FormCarrierCss)

			fmt.Println("Registered " + k.Id + " as available carrier.")
			Skip:
		}

		time.Sleep(time.Minute * 30)
	}
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Peek("User-Agent")) == "Googlebot" || strings.HasPrefix(string(ctx.Request.Header.Peek("User-Agent")), "AdsBot-Google") {
		env := stick.New(nil)
		file, _ := ioutil.ReadFile("./frontend/bots/google.twig")
		ctx.SetContentType("text/html")
		env.Execute(string(file), ctx, map[string]stick.Value{
			"title":   "ender confessions",
			"content": "Από το πληκτρολόγιό σου στην αγαπημένη σου σελίδα με ένα κλικ!",
		})
	} else if string(ctx.Request.Header.Peek("User-Agent")) == "Facebot" || strings.HasPrefix(string(ctx.Request.Header.Peek("User-Agent")), "facebookexternalhit") {
			env := stick.New(nil)
			file, _ := ioutil.ReadFile("./frontend/bots/facebook.twig")
			ctx.SetContentType("text/html")
			env.Execute(string(file), ctx, map[string]stick.Value{
				"title": "ender confessions",
				"url": string(ctx.Path()),
				"image": "https://confessions.ender.gr/assets/ender-confessions-opengraph.png",
				"content": "Από το πληκτρολόγιό σου στην αγαπημένη σου σελίδα με ένα κλικ!",
			})
	} else {
		file, _ := ioutil.ReadFile("./frontend/index.html")
		ctx.SetContentType("text/html")
		ctx.SetBody(file)
	}
}

func adminHandler(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("https://admin.ender.gr/confessions/", 301)
}

var fsHandler fasthttp.RequestHandler
func notFoundHandler(ctx *fasthttp.RequestCtx) {
	fsHandler(ctx)
	if ctx.Response.StatusCode() / 100 == 4 {
		ctx.Redirect("/#!/" + string(ctx.Response.StatusCode()), 302)
	}
}

func connect() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	model.CarrierCollection = session.DB("ender-confessions").C("Carrier")
	model.SecretCollection = session.DB("ender-confessions").C("Secret")
	model.SecretArchiveCollection = session.DB("ender-confessions").C("SecretArchive")
}

func start(router *fasthttprouter.Router, port string) {
	fmt.Printf("ender confessions running on %s\n", port)

	err := fasthttp.ListenAndServe(port, router.Handler)
	if err != nil {
		fmt.Printf("confessions could not start!\nError: %s\n", err.Error())
	}
}
