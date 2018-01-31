package model

import (
	"strings"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type SecretStatus string
var SecretCollection *mgo.Collection
var SecretArchiveCollection *mgo.Collection

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
	StatusDescription string `bson:"statusDescription" json:"status_description"`

	PublishData SecretPublishData `bson:"publishData" json:"publish_data"`
	SourceData SecretSourceData `bson:"sourceData" json:"source_data"`
	ChecksData []SecretSourceData `bson:"checksData" json:"checks_data"`
	Statistics map[string]int64 `bson:"statistics" json:"statistics"`

	Content string `json:"content"`
	OriginalContent string `json:"original_content"`
	ImageId string `bson:"imageId" json:"image_id"`
	Options map[string]string `json:"options"`
	Boosted bool `json:"boosted"`

	FinalForm string `bson:"-" json:"final_form"`
	Properties string `bson:"-" json:"properties"`
}

type SecretPublishData struct {
	QueuedTime int64 `bson:"queuedTime" json:"queued_time"`
	PublishTime int64 `bson:"publishTime" json:"publish_time"`
	Publisher string `bson:"publisher" json:"publisher"`
	PublishTag string `bson:"publishTag" json:"publish_tag"`
	FacebookPostId string `bson:"facebookPostId" json:"facebook_post_id"`
}

type SecretSourceData struct {
	Timestamp int64 `json:"timestamp"`
	IpAddress string `bson:"ipAddress" json:"ip_address"`
	IpPort int `bson:"ipPort" json:"ip_port"`
	Hostname string `bson:"hostname" json:"hostname"`
	Country string `bson:"country" json:"country"`
	RayID string `bson:"rayId" json:"ray_id"`
}

func FindSecret(id string) (Secret, error) {
	var secret Secret
	err := SecretCollection.Find(bson.M{"id": id}).One(&secret)
	if err != nil {
		err = SecretArchiveCollection.Find(bson.M{"id": id}).One(&secret)
		if err != nil {
			return Secret{}, errors.New("Secret could not be found!")
		} else {
			return secret, nil
		}
	} else {
		return secret, nil
	}
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

func (secret Secret) Save() {
	SecretCollection.Upsert(bson.M{"id": secret.Id}, bson.M{"$set": secret})
}