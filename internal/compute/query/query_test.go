package query

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ParseQueryStrSuite struct {
	suite.Suite
}

func (suite *ParseQueryStrSuite) TestUnknownCommand1() {
	// Arrange
	const queryStr = "test unknown command"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.Error(err)
	suite.Nil(q)
}

func (suite *ParseQueryStrSuite) TestUnknownCommand2() {
	// Arrange
	const queryStr = "nothing"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.Error(err)
	suite.Nil(q)
}

func (suite *ParseQueryStrSuite) TestGETSuccess() {
	// Arrange
	const queryStr = "GET key"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.NoError(err)
	suite.Require().NotNil(q)
	suite.Equal(CommandGET, q.Command())
	suite.Equal(1, len(q.args))
	suite.Equal("key", q.Arg(0))
}

func (suite *ParseQueryStrSuite) TestGETError() {
	// Arrange
	const queryStr = "GET key value"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.Error(err)
	suite.Nil(q)
}

func (suite *ParseQueryStrSuite) TestSETSuccess() {
	// Arrange
	const queryStr = "SET key value"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.NoError(err)
	suite.Require().NotNil(q)
	suite.Equal(CommandSET, q.Command())
	suite.Equal(2, len(q.args))
	suite.Equal("key", q.Arg(0))
	suite.Equal("value", q.Arg(1))
}

func (suite *ParseQueryStrSuite) TestSETError() {
	// Arrange
	const queryStr = "SET key"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.Error(err)
	suite.Nil(q)
}

func (suite *ParseQueryStrSuite) TestDELSuccess() {
	// Arrange
	const queryStr = "DEL key"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.NoError(err)
	suite.Require().NotNil(q)
	suite.Equal(CommandDEL, q.Command())
	suite.Equal(1, len(q.args))
	suite.Equal("key", q.Arg(0))
}

func (suite *ParseQueryStrSuite) TestDELError() {
	// Arrange
	const queryStr = "DEL key value"

	// Act
	q, err := ParseQueryStr(queryStr)

	// Assert
	suite.Error(err)
	suite.Nil(q)
}

func TestParseQueryStr(t *testing.T) {
	suite.Run(t, new(ParseQueryStrSuite))
}
