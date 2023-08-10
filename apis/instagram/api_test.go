package instagram_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/themis/apis/instagram"
)

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

type ApiTestSuite struct {
	suite.Suite

	api *instagram.API
}

func (suite *ApiTestSuite) SetupSuite() {
	suite.api = instagram.NewAPI()
}

func (suite *ApiTestSuite) TestGetUser() {
	user, err := suite.api.GetUserMedia("IGQVJWejVfNDBQVkRBeVp4VDg0RkxFczJsVVR1NWZAudUNLRnQtV2RjckNBam1RcmFuWTlVT1dzZA0hnTFdfTHNZAQ3FrOVdXQTFJN2pzTmd3MmVHdHV0UWpmYUdsNEh3Y3VRZA012VlBUZAUdjTUtCeVZABaAZDZD")
	suite.Require().NoError(err)
	suite.Require().NotEmpty(user.Username)
	suite.Require().NotEmpty(user.Caption)
}
