package model



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
	Statistics map[string]interface{} `bson:"statistics"`
}

type CarrierForm struct {
	Enabled bool `json:"enabled"`
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	SecretPrompt string `bson:"secretPrompt" json:"secretPrompt"`
	ImagePrompt string `bson:"imagePrompt" json:"imagePrompt"`
	SubmitPrompt string `bson:"submitPrompt" json:"submitPrompt"`
	SentMessage string `bson:"sentMessage" json:"sentMessage"`
	BackgroundUrl string `bson:"backgroundUrl" json:"backgroundUrl"`
	CustomCss string `bson:"customCss" json:"customCss"`

	AccentColor string `bson:"accentColor" json:"-"`
	TitleColor string `bson:"titleColor" json:"titleColor"`

	OptionSets map[string]CarrierOptions `bson:"optionSets" json:"optionSets"`

	IsAcceptsImage bool `bson:"isAcceptsImage" json:"acceptsImage"`
	IsEnableCaptcha bool `bson:"isEnableCaptcha" json:"requiresCaptcha"`
}

type CarrierFacebookInfo struct {
	PageName string `bson:"pageName"`
	UserName string `bson:"userName"`
	UserToken CarrierFacebookToken `bson:"userToken"`
	PageToken CarrierFacebookToken `bson:"pageToken"`
	HasToken bool `bson:"hasToken"`
}

type CarrierFacebookToken struct {
	AccessToken string `bson:"access_token" json:"accessToken"`
	MachineId string `bson:"machine_id" json:"machineId"`
	ExpiresIn int64 `bson:"expires_in" json:"expiresIn"`
}

type CarrierOptions struct {
	Name string `json:"name"`
	Options []string `json:"options"`
	OptionDisplay map[string]string `bson:"optionsDisplay" json:"optionsDisplay"`
	AllowCustom bool `bson:"allowCustom" json:"allowCustom"`

	SubmittedValue string `bson:"-" json:"submittedValue"`
}