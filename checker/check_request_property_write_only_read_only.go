package checker

import (
	"github.com/Barterio/oasdiff/diff"
	"golang.org/x/exp/slices"
)

const (
	RequestOptionalPropertyBecameNonWriteOnlyCheckId = "request-optional-property-became-not-write-only"
	RequestOptionalPropertyBecameWriteOnlyCheckId    = "request-optional-property-became-write-only"
	RequestOptionalPropertyBecameReadOnlyCheckId     = "request-optional-property-became-read-only"
	RequestOptionalPropertyBecameNonReadOnlyCheckId  = "request-optional-property-became-not-read-only"
	RequestRequiredPropertyBecameNonWriteOnlyCheckId = "request-required-property-became-not-write-only"
	RequestRequiredPropertyBecameWriteOnlyCheckId    = "request-required-property-became-write-only"
	RequestRequiredPropertyBecameReadOnlyCheckId     = "request-required-property-became-read-only"
	RequestRequiredPropertyBecameNonReadOnlyCheckId  = "request-required-property-became-not-read-only"
)

func RequestPropertyWriteOnlyReadOnlyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
						writeOnlyDiff := propertyDiff.WriteOnlyDiff
						if writeOnlyDiff == nil {
							return
						}
						if parent.Revision.Properties[propertyName] == nil {
							// removed properties processed by the RequestOptionalPropertyUpdatedCheck check
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if slices.Contains(parent.Base.Required, propertyName) {
							id := RequestRequiredPropertyBecameNonWriteOnlyCheckId
							if writeOnlyDiff.To == true {
								id = RequestRequiredPropertyBecameWriteOnlyCheckId
							}

							result = append(result, NewApiChange(
								id,
								config,
								[]any{propName},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
							return
						}

						id := RequestOptionalPropertyBecameNonWriteOnlyCheckId
						if writeOnlyDiff.To == true {
							id = RequestOptionalPropertyBecameWriteOnlyCheckId
						}
						result = append(result, NewApiChange(
							id,
							config,
							[]any{propName},
							"",
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					})

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						readOnlyDiff := propertyDiff.ReadOnlyDiff
						if readOnlyDiff == nil {
							return
						}
						if parent.Revision.Properties[propertyName] == nil {
							// removed properties processed by the RequestOptionalPropertyUpdatedCheck check
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if slices.Contains(parent.Base.Required, propertyName) {
							id := RequestRequiredPropertyBecameNonReadOnlyCheckId
							if readOnlyDiff.To == true {
								id = RequestRequiredPropertyBecameReadOnlyCheckId
							}
							result = append(result, NewApiChange(
								id,
								config,
								[]any{propName},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
							return
						}

						id := RequestOptionalPropertyBecameNonReadOnlyCheckId
						if readOnlyDiff.To == true {
							id = RequestOptionalPropertyBecameReadOnlyCheckId
						}
						result = append(result, NewApiChange(
							id,
							config,
							[]any{propName},
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
