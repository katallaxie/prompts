package prompts_test

import (
	"testing"

	"github.com/katallaxie/prompts"

	"github.com/stretchr/testify/require"
)

func TestNewParamOneOf(t *testing.T) {
	t.Parallel()

	params := prompts.NewParamOneOf()
	params.ParamOneOf = &prompts.ParamOneOfString{}
	require.NotNil(t, params)
	require.NotNil(t, params.GetTool())
	require.NotNil(t, params.GetString())
	require.Nil(t, params.GetInt())
}
