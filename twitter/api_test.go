package twitter_test

import (
	"github.com/desmos-labs/themis/twitter"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

type ApiTestSuite struct {
	suite.Suite

	api *twitter.Api
}

func (suite *ApiTestSuite) SetupSuite() {
	bearer := os.Getenv("TWITTER_BEARER")
	suite.api = twitter.NewApi(bearer)
}

func (suite *ApiTestSuite) TestGetTweet() {
	_, err := suite.api.GetTweet("1306130480622964736")
	suite.Require().NoError(err)
}

func (suite *ApiTestSuite) TestGetBio() {
	_, err := suite.api.GetUser("ricmontagnin")
	suite.Require().NoError(err)
}
