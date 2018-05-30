package main

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/enderian/confessions/database"
	"github.com/enderian/confessions/form"
	"github.com/enderian/confessions/index"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Configuration struct {
	Port              string `json:"port"`
	ConfessionsImages string `json:"confessions_images"`

	ReCaptchaSiteKey    string `json:"recaptcha_key"`
	ReCaptchaSiteSecret string `json:"recaptcha_secret"`

	ServiceAlert string `json:"service_alert"`
}

func main() {

	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	config := Configuration{}
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Unable to open configuration file config.json: %s\n", err.Error())
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Unable to open configuration file config.json: %s\n", err.Error())
	}

	database.InitConfessionsDatabase()
	router := fasthttprouter.New()

	form.ServiceAlert = config.ServiceAlert
	form.ReCaptchaSiteKey = config.ReCaptchaSiteKey
	form.ReCaptchaSiteSecret = config.ReCaptchaSiteSecret
	form.ImageDirectory = config.ConfessionsImages
	form.SetupForm()

	go registerCarriers(router)
	index.RegisterIndex(router)

	router.GET("/admin", adminHandler)

	router.NotFound = fasthttp.FSHandler("./frontend", 0)
	start(router, config.Port)
}

func registerCarriers(router *fasthttprouter.Router) {

	var registered []string
	for {
		for _, k := range database.FindCarriers() {
			for _, b := range registered {
				if b == k.Id {
					goto Skip
				}
			}
			registered = append(registered, k.Id)

			router.GET("/"+k.Id, form.CarrierForm)
			router.POST("/"+k.Id+"/submit", form.SecretSubmit)

			router.GET("/"+k.Id+"/secret/:id", form.StatusRead)
			router.PATCH("/"+k.Id+"/secret/:id", form.StatusPatch)

			log.Println("Registered " + k.Id + " as an available carrier.")
		Skip:
		}

		time.Sleep(time.Minute * 30)
	}
}

func adminHandler(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("https://admin.github.com/confessions/", 301)
}

func start(router *fasthttprouter.Router, port string) {
	log.Printf("ender confessions running on %s\n", port)

	err := fasthttp.ListenAndServe(port, router.Handler)
	if err != nil {
		log.Fatalf("confessions could not start!\nError: %s\n", err.Error())
	}
}
