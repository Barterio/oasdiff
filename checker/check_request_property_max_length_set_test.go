package checker_test

import (
	"testing"

	"github.com/Barterio/oasdiff/checker"
	"github.com/Barterio/oasdiff/diff"
	"github.com/Barterio/oasdiff/load"
	"github.com/stretchr/testify/require"
)

// CL: setting maxLength of request body
func TestRequestBodyMaxLengthSetCheck(t *testing.T) {
	s1, err := open("../data/checker/request_body_max_length_set_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_max_length_set_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxLengthSetCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMaxLengthSetId,
		Args:        []any{uint64(15)},
		Level:       checker.WARN,
		Comment:     checker.RequestBodyMaxLengthSetId + "-comment",
		Operation:   "POST",
		OperationId: "addPet",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_body_max_length_set_revision.yaml"),
	}, errs[0])
}

// CL: setting maxLength of request propreties
func TestRequestPropertyMaxLengthSetCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_max_length_set_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_max_length_set_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxLengthSetCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMaxLengthSetId,
		Args:        []any{"age", uint64(15)},
		Level:       checker.WARN,
		Comment:     checker.RequestPropertyMaxLengthSetId + "-comment",
		Operation:   "POST",
		OperationId: "addPet",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_max_length_set_revision.yaml"),
	}, errs[0])
}
