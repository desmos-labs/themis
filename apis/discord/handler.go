package discord

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/desmos-labs/themis/utils"
	"path"
	"strings"
)

type Handler struct {
	cfg    *Config
	pubKey *rsa.PublicKey
}

func NewHandler(cfg *Config) *Handler {
	pubKey, err := utils.ReadPublicKeyFromFile(cfg.HephaestusPubKeyPath)
	if err != nil {
		panic(err)
	}

	return &Handler{
		cfg:    cfg,
		pubKey: pubKey,
	}
}

func (h *Handler) getFilePathByUsername(username string) string {
	return path.Join(h.cfg.StoreFolderPath, strings.ToLower(username))
}

// HandleSaveData handles the request that is performed when creating some .data
func (h *Handler) HandleSaveData(req SaveDataReq) error {
	// 1. Verify the signature
	msg, err := req.VerificationData.ToSignBytes()
	if err != nil {
		return err
	}

	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		return err
	}

	err = rsa.VerifyPSS(h.pubKey, crypto.SHA256, msg, signature, nil)
	if err != nil {
		return fmt.Errorf("invalid signature provided")
	}

	// 2. Write the file
	filePath := h.getFilePathByUsername(req.VerificationData.Username)
	return utils.WriteFile(filePath, req.VerificationData)
}

// GetVerificationDataForUser returns the verification .data for the user.
// If no .data is found, nil is returned instead.
func (h *Handler) GetVerificationDataForUser(user string) (*VerificationData, error) {
	filePath := h.getFilePathByUsername(user)

	// Read the file
	content, err := utils.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, nil
	}

	// Parse the contents
	var data VerificationData
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
