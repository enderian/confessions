package main

import (
	"ender.gr/confessions/form"
	"ender.gr/confessions/model"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"encoding/json"
	"time"
	"log"
	"os"
	"io"
	"ender.gr/confessions/index"
)

type Configuration struct {
	Port string `json:"port"`
	ConfessionsImages string `json:"confessions_images"`

	ReCaptchaSiteKey string `json:"recaptcha_key"`
	ReCaptchaSiteSecret string `json:"recaptcha_secret"`
}

func main() {

	logFile, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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

	connect()
	router := fasthttprouter.New()

	form.ReCaptchaSiteKey = config.ReCaptchaSiteKey
	form.ReCaptchaSiteSecret = config.ReCaptchaSiteSecret
	form.ImageDirectory = config.ConfessionsImages
	form.SetupForm()

	go registerCarriers(router)
	index.RegisterIndex(router)

	router.POST("/secret", form.StatusRead)
	router.POST("/submit", form.SecretSubmit)
    router.GET("/admin", adminHandler)

	router.NotFound = fasthttp.FSHandler("./frontend", 0)
	start(router, config.Port)
}

func registerCarriers(router *fasthttprouter.Router) {

	var registered []string
	for {
		for _, k := range model.FindCarriers() {
			for _, b := range registered {
				if b == k.Id {
					goto Skip
				}
			}
			registered = append(registered, k.Id)
			router.GET("/" + k.Id, form.CarrierForm)
			log.Println("Registered " + k.Id + " as an available carrier.")
			Skip:
		}

		time.Sleep(time.Minute * 30)
	}
}

func adminHandler(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("https://admin.ender.gr/confessions/", 301)
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
	log.Printf("ender confessions running on %s\n", port)

	err := fasthttp.ListenAndServe(port, router.Handler)
	if err != nil {
		log.Fatalf("confessions could not start!\nError: %s\n", err.Error())
	}
}
