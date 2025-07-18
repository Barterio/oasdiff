package checker

import (
	"github.com/Barterio/oasdiff/diff"
)

const (
	ResponsePropertyBecameNullableId = "response-property-became-nullable"
	ResponseBodyBecameNullableId     = "response-body-became-nullable"
)

func ResponsePropertyBecameNullableCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					if mediaTypeDiff.SchemaDiff == nil {
						continue
					}

					if mediaTypeDiff.SchemaDiff.NullableDiff != nil && mediaTypeDiff.SchemaDiff.NullableDiff.To == true {
						result = append(result, NewApiChange(
							ResponseBodyBecameNullableId,
							config,
							nil,
							"",
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							nullableDiff := propertyDiff.NullableDiff
							if nullableDiff == nil {
								return
							}
							if nullableDiff.To != true {
								return
							}

							result = append(result, NewApiChange(
								ResponsePropertyBecameNullableId,
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
