package models

// CarOwnerProp represents key/value property
type CarOwnerProp struct {
	ID    string `bson:"_id,omitempty"`
	Value string `bson:"value"`
}
