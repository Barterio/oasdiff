package checker

import (
	"github.com/Barterio/oasdiff/diff"
	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

const (
	ResponseRequiredPropertyRemovedId          = "response-required-property-removed"
	ResponseRequiredWriteOnlyPropertyRemovedId = "response-required-write-only-property-removed"
	ResponseRequiredPropertyAddedId            = "response-required-property-added"
	ResponseRequiredWriteOnlyPropertyAddedId   = "response-required-write-only-property-added"
)

func ResponseRequiredPropertyUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {

			if operationItem.ResponsesDiff == nil {
				continue
			}

			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					CheckDeletedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
							id := ResponseRequiredPropertyRemovedId
							if propertyItem.WriteOnly {
								id = ResponseRequiredWriteOnlyPropertyRemovedId
							}
							if !slices.Contains(parent.Base.Required, propertyName) {
								// Covered by response-optional-property-removed
								return
							}

							result = append(result, NewApiChange(
								id,
								config,
								[]any{propertyFullName(propertyPath, propertyName), responseStatus},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						})
					CheckAddedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
							id := ResponseRequiredPropertyAddedId
							if propertyItem.WriteOnly {
								id = ResponseRequiredWriteOnlyPropertyAddedId
							}
							if !slices.Contains(parent.Revision.Required, propertyName) {
								// Covered by response-optional-property-added
								return
							}

							result = append(result, NewApiChange(
								id,
								config,
								[]any{propertyFullName(propertyPath, propertyName), responseStatus},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						})
				}
			}
		}
	}
	return result
}
