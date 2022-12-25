package hephaestus

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/desmos-labs/themis/apis/types"

	"github.com/desmos-labs/themis/apis/utils"
)

type Handler struct {
	storeFolderPath string
	pubKey          *rsa.PublicKey
}

func NewHandler(storeFolder string, cfg *Config) *Handler {
	pubKey, err := utils.ReadPublicKeyFromFile(cfg.PubKeyPath)
	if err != nil {
		panic(err)
	}

	return &Handler{
		storeFolderPath: storeFolder,
		pubKey:          pubKey,
	}
}

func (h *Handler) getFilePathByUsername(username string) string {
	return path.Join(h.storeFolderPath, strings.ToLower(username))
}

func (h *Handler) ParseSaveDataRequest(req *http.Request) (types.SaveDataReq, error) {
	jsonBz, err := io.ReadAll(req.Body)
	if err != nil {
		return types.SaveDataReq{}, err
	}

	var data types.SaveDataReq
	return data, json.Unmarshal(jsonBz, &data)
}

// HandleSaveData handles the request that is performed when creating some data
func (h *Handler) HandleSaveData(req types.SaveDataReq) error {
	// 1. Verify the signature
	msg, err := req.VerificationData.ToSignBytes()
	if err != nil {
		return err
	}

	signature, err := hex.DecodeString(req.BotSignature)
	if err != nil {
		return err
	}

	err = rsa.VerifyPSS(h.pubKey, crypto.SHA256, msg, signature, nil)
	if err != nil {
		return fmt.Errorf("invalid signature provided")
	}

	// 2. Write the file
	filePath := h.getFilePathByUsername(req.Username)
	return utils.WriteFile(filePath, req.VerificationData)
}

// GetVerificationDataForUser returns the verification data for the user.
// If no data is found, nil is returned instead.
func (h *Handler) GetVerificationDataForUser(user string) (*types.VerificationData, error) {
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
	var data types.VerificationData
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
