package parse_test

import (
	"testing"

	"github.com/rowan-gud/pakk/config/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParseCommandTestSuite struct {
	suite.Suite
}

func (s *ParseCommandTestSuite) TestParseFromString() {
	s.Run("ValidCommandString", func() {
		assert := assert.New(s.T())

		var cmd parse.Command
		err := cmd.UnmarshalTOML(
			`foo bar 'should quote"'"should quote'"this \"`,
		)

		if assert.NoError(err) {
			assert.Equal(
				`foo bar 'should quote"'"should quote'"this \"`,
				cmd.Raw(),
			)

			assert.NotNil(cmd.Cmd)
		}
	})
}

func (s *ParseCommandTestSuite) TestParseFromArray() {
	s.Run("ValidCommandArray", func() {
		assert := assert.New(s.T())

		var cmd parse.Command
		err := cmd.UnmarshalTOML([]any{
			"a", "b", "c",
		})

		if assert.NoError(err) {
			assert.ElementsMatch(cmd.Raw(), []string{
				"a", "b", "c",
			})

			assert.NotNil(cmd.Cmd)
		}
	})

	s.Run("NotValidElements", func() {
		assert := assert.New(s.T())

		var cmd parse.Command
		err := cmd.UnmarshalTOML([]any{
			"a", 12, "c",
		})

		if assert.Error(err) {
			assert.Equal(err.Error(), "cannot unmarshal type int to type string")
		}
	})
}

func TestParseCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ParseCommandTestSuite))
}
