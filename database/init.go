package database

import (
	"gopkg.in/mgo.v2"
)

var (
	Username string
	Password string
	Address  string
)

func InitConfessionsDatabase() {
	session, err := mgo.Dial(Address)
	if err != nil {
		panic(err)
	}
	if Username != "" {
		if err = session.Login(&mgo.Credential{
			Username: Username,
			Password: Password,
		}); err != nil {
			panic(err.Error())
		}
	}
	session.SetMode(mgo.Monotonic, true)

	carrierCollection = session.DB("ender-confessions").C("Carrier")
	secretCollection = session.DB("ender-confessions").C("Secret")
	secretArchiveCollection = session.DB("ender-confessions").C("SecretArchive")
}
