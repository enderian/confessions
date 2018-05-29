package database

import (
	"gopkg.in/mgo.v2"
	"errors"
	"log"
	"github.com/enderian/confessions/model"
	"gopkg.in/mgo.v2/bson"
)

var carrierCollection *mgo.Collection
var carrierArchiveCollection *mgo.Collection

func FindCarrier(id string) (model.Carrier, error) {
	var carrier model.Carrier
	err := carrierCollection.Find(bson.M{"id": id}).One(&carrier)
	if err != nil {
		return model.Carrier{}, errors.New("carrier could not be found")
	} else {
		return carrier, nil
	}
}

func FindCarrierByFacebookPage(page string) (model.Carrier, error) {
	var carrier model.Carrier
	err := carrierCollection.Find(bson.M{"facebookPage": page}).One(&carrier)
	if err != nil {
		return model.Carrier{}, errors.New("carrier could not be found")
	} else {
		return carrier, nil
	}
}

func FindCarriers() []model.Carrier {
	var results []model.Carrier
	carrierCollection.Find(bson.M{}).All(&results)
	return results
}

func SaveCarrier(carrier model.Carrier) {
	_, err := carrierCollection.Upsert(bson.M{"id": carrier.Id}, bson.M{"$set": carrier})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func DeleteCarrier(carrier model.Carrier) {
	go func() {
		secretCursor := secretCollection.Find(bson.M{"carrier":  carrier.Id}).Iter()
		var secret model.Secret
		for secretCursor.Next(&secret) {
			secret.Carrier = carrier.Id + "_deleted"
			secretArchiveCollection.Insert(secret)
		}
		secretCollection.RemoveAll(bson.M{"carrier":  carrier.Id})
	}()

	carrierArchiveCollection.Insert(carrier)
	carrierCollection.Remove(bson.M{"id": carrier.Id})
}

