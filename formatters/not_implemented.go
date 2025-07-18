package formatters

import (
	"fmt"

	"github.com/Barterio/oasdiff/checker"
	"github.com/Barterio/oasdiff/diff"
	"github.com/Barterio/oasdiff/load"
	"github.com/getkin/kin-openapi/openapi3"
)

type notImplementedFormatter struct{}

func (f notImplementedFormatter) RenderDiff(*diff.Diff, RenderOpts) ([]byte, error) {
	return notImplemented()
}

func (f notImplementedFormatter) RenderSummary(*diff.Diff, RenderOpts) ([]byte, error) {
	return notImplemented()
}

func (f notImplementedFormatter) RenderChangelog(checker.Changes, RenderOpts, *load.SpecInfoPair) ([]byte, error) {
	return notImplemented()
}

func (f notImplementedFormatter) RenderChecks(Checks, RenderOpts) ([]byte, error) {
	return notImplemented()
}

func (f notImplementedFormatter) RenderFlatten(*openapi3.T, RenderOpts) ([]byte, error) {
	return notImplemented()
}

func notImplemented() ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
