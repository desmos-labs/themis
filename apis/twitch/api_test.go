package twitch_test

import (
	"github.com/desmos-labs/themis/apis/twitch"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

type ApiTestSuite struct {
	suite.Suite

	api *twitch.API
}

func (suite *ApiTestSuite) SetupSuite() {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	suite.api = twitch.NewAPI(clientID, clientSecret)
}

func (suite *ApiTestSuite) TestGetBio() {
	_, err := suite.api.GetBio("riccardomontagnin")
	suite.Require().NoError(err)
}
