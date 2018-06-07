package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/enderian/confessions/database"
	"github.com/enderian/confessions/form"
	"github.com/enderian/confessions/index"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"os"
	"time"
	"gopkg.in/ini.v1"
)

func main() {

	logFile, err := os.OpenFile("confessions.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Fail to read configuration file: %v", err)
	}

	database.Address = cfg.Section("database").Key("address").MustString("localhost")
	database.Username = cfg.Section("database").Key("username").MustString("")
	database.Password = cfg.Section("database").Key("password").MustString("")
	database.Init()
	router := fasthttprouter.New()

	form.ServiceAlert = cfg.Section("").Key("service_alert").MustString("")
	form.ReCaptchaSiteKey = cfg.Section("recaptcha").Key("key").String()
	form.ReCaptchaSiteSecret = cfg.Section("recaptcha").Key("secret").String()
	form.SetupForm()
	form.SetupS3()

	go registerCarriers(router)
	index.RegisterIndex(router)

	router.GET("/admin", adminHandler)

	router.NotFound = fasthttp.FSHandler("./frontend", 0)
	start(router, cfg.Section("").Key("port").MustString(":8080"))
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
