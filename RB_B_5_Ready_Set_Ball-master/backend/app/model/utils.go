package model

import (
	"gopkg.in/mgo.v2/bson"
)

// mapKeys returns the a slice of [bson.ObjectId]s from a map -> map[bson.ObjectId]float64
func mapKeys(m map[string]float64) []bson.ObjectId {
	keys := make([]bson.ObjectId, len(m))
	i := 0
	for k := range m {
		if bson.IsObjectIdHex(k) {
			keys[i] = bson.ObjectIdHex(k)
		} else {
			keys[i] = bson.ObjectId(k)
		}
		i++
	}

	return keys
}

// IDsToHex converts a slice of [bson.ObjectId] to a slice of the hex representations
//
// It creates a slice of the size of the initial slice + [add]
func IDsToHex(ids []bson.ObjectId, add int) []string {
	hexes := make([]string, len(ids)+add)
	for i := range ids {
		hexes[i] = ids[i].Hex()
	}

	return hexes
}
