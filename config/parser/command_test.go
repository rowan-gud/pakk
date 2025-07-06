package parser_test

import (
	"testing"

	"github.com/rowan-gud/pakk/config/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParseCommandTestSuite struct {
	suite.Suite
}

func (s *ParseCommandTestSuite) TestParseFromString() {
	s.Run("ValidCommandString", func() {
		assert := assert.New(s.T())

		cmd, err := parser.ParseCommand(`foo bar 'should quote"'"should quote'"this \"`)

		if assert.NoError(err) {
			assert.Equal([]string{
				"foo", "bar", "should quote\"", "should quote'", "this", "\"",
			}, cmd)
		}
	})
}

func (s *ParseCommandTestSuite) TestParseFromArray() {
	s.Run("ValidCommandArray", func() {
		assert := assert.New(s.T())

		cmd, err := parser.ParseCommand([]any{
			"a", "b", "c",
		})

		if assert.NoError(err) {
			assert.ElementsMatch(cmd, []string{
				"a", "b", "c",
			})
		}
	})

	s.Run("NotValidElements", func() {
		assert := assert.New(s.T())

		_, err := parser.ParseCommand([]any{
			"a", 12, "c",
		})

		if assert.Error(err) {
			assert.Equal(err.Error(), "expected elements to be of type string found int")
		}
	})
}

func TestParseCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ParseCommandTestSuite))
}
