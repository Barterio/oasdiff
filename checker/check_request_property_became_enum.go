package checker

import (
	"github.com/Barterio/oasdiff/diff"
)

const (
	RequestPropertyBecameEnumId = "request-property-became-enum"
)

func RequestPropertyBecameEnumCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {

			if operationItem.RequestBodyDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}
			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {

						if enumDiff := propertyDiff.EnumDiff; enumDiff == nil || !enumDiff.EnumAdded {
							return
						}

						result = append(result, NewApiChange(
							RequestPropertyBecameEnumId,
							config,
							[]any{propertyFullName(propertyPath, propertyName)},
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
	return result
}
