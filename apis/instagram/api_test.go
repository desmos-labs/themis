package instagram_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/themis/instagram"
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
	user, err := suite.api.GetUser("leobragaz")
	suite.Require().NoError(err)
	suite.Require().NotEmpty(user.Username)
}
