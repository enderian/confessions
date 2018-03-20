package model

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"log"
)

var CarrierCollection *mgo.Collection
var CarrierArchiveCollection *mgo.Collection

type Carrier struct {
	Id string
	Name string

	EffectiveFormat string `bson:"effectiveFormat"`
	EffectiveHashtag string `bson:"effectiveHash"`
	EffectiveId int `bson:"effectiveId"`

	FacebookPage string `bson:"facebookPage"`
	FacebookInfo CarrierFacebookInfo `bson:"facebookInformation"`
	Form CarrierForm

	BoostedHosts map[string]string `bson:"boostedHosts"`
	Statistics map[string]int64 `bson:"statistics"`
}

type CarrierForm struct {
	Enabled bool `json:"enabled"`
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	SecretPrompt string `bson:"secretPrompt" json:"secret_prompt"`
	ImagePrompt string `bson:"imagePrompt" json:"image_prompt"`
	SubmitPrompt string `bson:"submitPrompt" json:"submit_prompt"`
	SentMessage string `bson:"sentMessage" json:"sent_message"`
	BackgroundUrl string `bson:"backgroundUrl" json:"background_url"`
	LogoUrl string `bson:"logoUrl" json:"logo_url"`
	TitleColor string `bson:"titleColor" json:"title_color"`
	CustomCss string `bson:"customCss" json:"custom_css"`

	OptionSets map[string]CarrierOptions `bson:"optionSets" json:"option_sets"`

	IsAcceptsImage bool `bson:"isAcceptsImage" json:"is_accepts_image"`
	IsEnableCaptcha bool `bson:"isEnableCaptcha" json:"is_enable_captcha"`

	IsBoosted bool `bson:"-" json:"is_boosted"`
	BoostedReason string `bson:"-" json:"boosted_reason"`
}

type CarrierFacebookInfo struct {
	PageName string `bson:"pageName"`
	UserName string `bson:"userName"`
	UserToken CarrierFacebookToken `bson:"userToken"`
	PageToken CarrierFacebookToken `bson:"pageToken"`
	HasToken bool `bson:"hasToken"`
}

type CarrierFacebookToken struct {
	AccessToken string `bson:"access_token" json:"access_token"`
	MachineId string `bson:"machine_id" json:"machine_id"`
	ExpiresIn int64 `bson:"expires_in" json:"expires_in"`
}

type CarrierOptions struct {
	Name string `json:"name"`
	Options []string `json:"options"`
	OptionDisplay map[string]string `bson:"optionsDisplay" json:"option_display"`
	AllowCustom bool `bson:"allowCustom" json:"allow_custom"`

	SubmittedValue string `bson:"-" json:"submitted_value"`
}

func FindCarrier(id string) (Carrier, error) {
	var carrier Carrier
	err := CarrierCollection.Find(bson.M{"id": id}).One(&carrier)
	if err != nil {
		return Carrier{}, errors.New("carrier could not be found")
	} else {
		return carrier, nil
	}
}

func FindCarriers() []Carrier {
	var results []Carrier
	CarrierCollection.Find(bson.M{}).All(&results)
	return results
}

func (carrier Carrier) IsBoostedHost(ip string, hostname string) (bool, string) {
	for host, reason := range carrier.BoostedHosts {
		host = strings.Replace(host, "_", ".", -1)
		if hostname == host || host == ip {
			return true, reason
		}
	}
	for host, reason := range carrier.BoostedHosts {
		host = strings.Replace(host, "_", ".", -1)
		if strings.HasSuffix(hostname, host) {
			return true, reason
		}
	}
	return false, ""
}

func (carrier Carrier) Save() {
	_, err := CarrierCollection.Upsert(bson.M{"id": carrier.Id}, bson.M{"$set": carrier})
	if err != nil {
		log.Fatal(err.Error())
	}
}
