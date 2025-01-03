package data_classification

const (
	VALUE_PATTERN_QUERY = `valuePattern: { operator: %s, value: "%s" }`
	ENV_SCOPED_QUERY = `scope: {
						type: ENVIRONMENT
						environmentScope: {
							environmentIds: [%s]
						}
					}`
	SCOPED_PATTERN_QUERY = `{
				name: "%s"
				urlMatchPatterns: [%s]
				patternType: %s
				locations: [%s]
				matchType: %s
				%s
				keyPattern: { operator: %s, value: "%s" }
				%s
			}`
	CREATE_QUERY = `mutation {
			createDataType(
				dataTypeRule: {
					name: "%s"
					description: "%s"
					scopedPatterns: [
						%s
					]
					enabled: %t
					dataSuppression: %s
					sensitivity: %s
					datasetIds: [
						%s
					]
				}
			) {
				id
				name
				description
			}
		}`
	UPDATE_QUERY = `mutation {
			updateDataType(
				dataTypeRule: {
					name: "%s"
					description: "%s"
					scopedPatterns: [
						%s
					]
					enabled: %t
					dataSuppression: %s
					sensitivity: %s
					datasetIds: [
						%s
					]
				}
			) {
				id
				name
				description
			}
		}`
	DELETE_QUERY = `mutation {
			deleteDataType(id: "%s") {
				success
				__typename
			}
		}`
	READ_QUERY = `{
			dataTypes {
				results {
					id
					name
					description
					scopedPatterns {
						name
						locations
						patternType
						matchType
						scope {
							type
							environmentScope {
								environmentIds
								__typename
							}
							__typename
						}
						urlMatchPatterns
						keyPattern {
							operator
							value
							__typename
						}
						valuePattern {
							operator
							value
							__typename
						}
						__typename
					}
					suppressionPattern
					enabled
					datasets {
						id
						name
						description
						iconType
						dataTypes {
							id
							__typename
						}
						__typename
					}
					dataSuppression
					sensitivity
					__typename
				}
				__typename
			}
		}`
)