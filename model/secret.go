package model

import (
	"strings"
	"strconv"
	"time"
)

type SecretStatus string

const(
	SENT SecretStatus = "SENT"
	QUEUED SecretStatus = "QUEUED"
	PUBLISHING SecretStatus = "PUBLISHING"
	PUBLISHED SecretStatus = "PUBLISHED"
	DELETED SecretStatus = "DELETED"
	DECLINED SecretStatus = "DECLINED"
	HIDDEN SecretStatus = "HIDDEN"
	REMOVED_FB SecretStatus = "REMOVED_FB"
	FAILED SecretStatus = "FAILED"
)

type Secret struct {
	Carrier string `json:"carrier"`
	Id string `json:"id"`
	Status SecretStatus `json:"status"`
	StatusDescription string `bson:"statusDescription" json:"statusDescription"`

	PublishData SecretPublishData `bson:"publishData" json:"publishData"`
	SourceData SecretSourceData `bson:"sourceData" json:"sourceData"`
	ChecksData []SecretSourceData `bson:"checksData" json:"checksData"`
	Statistics map[string]int64 `bson:"statistics" json:"statistics"`

	Content string `json:"content"`
	OriginalContent string `json:"originalContent"`
	ImageId string `bson:"imageId" json:"imageId"`
	Options map[string]string `json:"options"`
	Boosted bool `json:"boosted"`

	FinalForm string `bson:"-" json:"finalForm"`
	Properties string `bson:"-" json:"properties"`
}

type SecretPublishData struct {
	QueuedTime time.Time `bson:"queuedTime" json:"queuedTime"`
	PublishTime time.Time `bson:"publishTime" json:"publishTime"`
	Publisher string `bson:"publisher" json:"publisher"`
	PublishTag string `bson:"publishTag" json:"publishTag"`
	FacebookPostId string `bson:"facebookPostId" json:"facebookPostId"`
}

type SecretSourceData struct {
	Timestamp time.Time `json:"timestamp"`
	IpAddress string `bson:"ipAddress" json:"ipAddress"`
	IpPort int `bson:"ipPort" json:"ipPort"`
	Hostname string `bson:"hostname" json:"hostname"`
	Country string `bson:"country" json:"country"`
	RayID string `bson:"rayId" json:"rayId"`
}



func (secret *Secret) BuildProperties(carrier *Carrier) {
	secret.Properties = carrier.EffectiveFormat

	for k, v := range secret.Options {
		if _, ok := carrier.Form.OptionSets[k]; ok {
			secret.Properties = strings.Replace(secret.Properties, "{" + k + "}", v, -1)
		}
	}
	for id := range carrier.Form.OptionSets {
		secret.Properties = strings.Replace(secret.Properties, "{" + id + "}", "", -1)
	}

	for strings.Contains(secret.Properties, "--") {
		secret.Properties = strings.Replace(secret.Properties, "--", "-", -1)
	}
	for strings.Contains(secret.Properties, "-)") {
		secret.Properties = strings.Replace(secret.Properties, "-)", ")", -1)
	}
	for strings.Contains(secret.Properties, " )") {
		secret.Properties = strings.Replace(secret.Properties, " )", ")", -1)
	}
	for strings.Contains(secret.Properties, "(-") {
		secret.Properties = strings.Replace(secret.Properties, "(-", "(", -1)
	}
	for strings.Contains(secret.Properties, "( ") {
		secret.Properties = strings.Replace(secret.Properties, "( ", "(", -1)
	}
	for strings.Contains(secret.Properties, " (") {
		secret.Properties = strings.Replace(secret.Properties, " (", "(", -1)
	}
	for strings.Contains(secret.Properties, "{cn}{cn}") {
		secret.Properties = strings.Replace(secret.Properties, "{cn}{cn}", "{cn}", -1)
	}
	
	hasMessage := strings.Contains(secret.Properties, "{message}")
	finalForm := "#" + carrier.EffectiveHashtag + strconv.Itoa(carrier.EffectiveId)
	secret.Content = strings.Trim(secret.Content, " \t\n")

	secret.Properties = strings.Replace(secret.Properties, "{nn}", "\n", -1)
	secret.Properties = strings.Replace(secret.Properties, "{cn}", "\n", -1)
	secret.Properties = strings.Replace(secret.Properties, "{s}", " ", -1)
	secret.Properties = strings.Replace(secret.Properties, "{q}", "\"", -1)
	secret.Properties = strings.Replace(secret.Properties, "()", "", -1)
	secret.Properties = strings.Replace(secret.Properties, "(-)", "", -1)

	if !strings.HasPrefix(secret.Properties, ":") {
		finalForm += " "
	}
	if len(strings.Trim(secret.Properties, " \t")) != 0 {
		finalForm += secret.Properties
	}
	if !hasMessage {
		finalForm += secret.Content
	} else {
		finalForm = strings.Replace(finalForm, "{message}", secret.Content, -1)
	}

	secret.FinalForm = finalForm
	secret.Properties = strings.Replace(secret.Properties, "{message}", "?", -1)
}