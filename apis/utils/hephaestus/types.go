package hephaestus

import (
	"crypto/sha256"
	"encoding/json"
)

type SaveDataReq struct {
	Username         string           `json:"username"`
	VerificationData VerificationData `json:"verification_data"`

	BotSignature string `json:"signature"`
}

type VerificationData struct {
	Address   string `json:"address"`
	PubKey    string `json:"pub_key"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

func (v VerificationData) ToBytes() ([]byte, error) {
	return json.Marshal(&v)
}

// ToSignBytes returns the bytes representation of v that should be used when signing v
func (v VerificationData) ToSignBytes() ([]byte, error) {
	bz, err := v.ToBytes()
	if err != nil {
		return nil, err
	}

	// Hash the message using SHA-256
	msgHash := sha256.New()
	_, err = msgHash.Write(bz)
	if err != nil {
		return nil, err
	}

	return msgHash.Sum(nil), nil
}
