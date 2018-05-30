package database

import (
	"errors"
	"github.com/enderian/confessions/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var secretCollection *mgo.Collection
var secretArchiveCollection *mgo.Collection

func FindSecret(id string) (model.Secret, error) {
	var secret model.Secret
	err := secretCollection.Find(bson.M{"id": id}).One(&secret)
	if err != nil {
		err = secretArchiveCollection.Find(bson.M{"id": id}).One(&secret)
		if err != nil {
			return model.Secret{}, errors.New("secret could not be found")
		} else {
			return secret, nil
		}
	} else {
		return secret, nil
	}
}

func FindSecrets(query DocumentWrapper) *mgo.Query {
	return secretCollection.Find(query)
}

func FindArchivedSecrets(query DocumentWrapper) *mgo.Query {
	return secretArchiveCollection.Find(query)
}

func SaveSecret(secret model.Secret) {
	secretCollection.Upsert(bson.M{"id": secret.Id}, bson.M{"$set": secret})
}

func ArchiveSecret(secret model.Secret) {
	secretCollection.Remove(bson.M{"id": secret.Id})
	secretArchiveCollection.Insert(secret)
}
