package checker_test

import (
	"testing"

	"github.com/Barterio/oasdiff/checker"
	"github.com/stretchr/testify/require"
)

func TestGetOptionalRuleIds(t *testing.T) {
	require.Len(t, checker.GetOptionalRuleIds(), 7)
}
