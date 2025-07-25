package checker

import (
	"github.com/Barterio/oasdiff/diff"
	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

const (
	ResponseOptionalPropertyRemovedId          = "response-optional-property-removed"
	ResponseOptionalWriteOnlyPropertyRemovedId = "response-optional-write-only-property-removed"
	ResponseOptionalPropertyAddedId            = "response-optional-property-added"
	ResponseOptionalWriteOnlyPropertyAddedId   = "response-optional-write-only-property-added"
)

func ResponseOptionalPropertyUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
							id := ResponseOptionalPropertyRemovedId
							if propertyItem.WriteOnly {
								id = ResponseOptionalWriteOnlyPropertyRemovedId
							}
							if slices.Contains(parent.Base.Required, propertyName) {
								// covered by response-required-property-removed
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
							id := ResponseOptionalPropertyAddedId
							if propertyItem.WriteOnly {
								id = ResponseOptionalWriteOnlyPropertyAddedId
							}

							if slices.Contains(parent.Revision.Required, propertyName) {
								// covered by response-required-property-added
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
