package twitter_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/themis/twitter"
)

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

type ApiTestSuite struct {
	suite.Suite

	api *twitter.API
}

func (suite *ApiTestSuite) SetupSuite() {
	bearer := os.Getenv("TWITTER_BEARER")
	suite.api = twitter.NewAPI(bearer)
}

func (suite *ApiTestSuite) TestGetTweet() {
	_, err := suite.api.GetTweet("1306130480622964736")
	suite.Require().NoError(err)
}

func (suite *ApiTestSuite) TestGetBio() {
	_, err := suite.api.GetUser("ricmontagnin")
	suite.Require().NoError(err)
}
