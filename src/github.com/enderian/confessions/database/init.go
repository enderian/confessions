package database

import (
	"gopkg.in/mgo.v2"
)

func InitConfessionsDatabase() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	carrierCollection = session.DB("ender-confessions").C("Carrier")
	secretCollection = session.DB("ender-confessions").C("Secret")
	secretArchiveCollection = session.DB("ender-confessions").C("SecretArchive")
}