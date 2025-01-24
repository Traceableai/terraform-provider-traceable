package data_classification

const (
	VALUE_PATTERN_QUERY = `valuePattern: { operator: %s, value: "%s" }`
	ENV_SCOPED_QUERY    = `scope: {
						type: ENVIRONMENT
						environmentScope: {
							environmentIds: [%s]
						}
					}`
	CREATE_DATA_SET_QUERY = `mutation {
		createDataSet(
			dataSetCreate: {name: "%s", dataTypeIds: null, description: "%s", iconType: "%s"}
		) {
			id
			name
		}
		}`
	UPDATE_DATA_SET_QUERY = `mutation {
		updateDataSet(
			dataSetUpdate: {id: "%s",name: "%s", dataTypeIds: null, description: "%s", iconType: "%s"}
		) {
			id
			name
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

	SPAN_FILTER_OVERRIDES_QUERY = `spanFilter: {
					keyValueFilter: {
						keyPattern: { operator: %s, value: "%s" }
						%s
					}
				}`

	ENV_OVERRIDES_QUERY = `environmentScope: { environmentIds: [%s] }`

	CREATE_OVERRIDES_QUERY = `mutation {
		createDataClassificationOverride(
			input: {
				name: "%s"
				description: "%s"
				%s 
				%s
				behaviorOverride: { dataSuppressionOverride: %s }
			}
		) {
			id
			name
			__typename
		}
	}`
	UPDATE_OVERRIDES_QUERY = `mutation {
		updateDataClassificationOverride(
			input: {
				id: "%s"
				name: "%s"
				description: "%s"
				%s 
				%s
				behaviorOverride: { dataSuppressionOverride: %s }
			}
		) {
			id
			name
			__typename
		}
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
	DATA_SET_READ_QUERY = `{
	dataSets {
		results {
		id
		name
		description
		enabled
		dataSuppression
		dataTypes {
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
			__typename
		}
		sensitivity
		color
		iconType
		__typename
		}
		__typename
	}
	}`
	DELETE_QUERY = `mutation {
			deleteDataType(id: "%s") {
				success
				__typename
			}
		}`
	DELETE_DATA_SET_QUERY = `mutation {
			deleteDataSet(id: "%s")
		}`
	DELETE_OVERRIDES_QUERY = `mutation {
		deleteDataClassificationOverride(
			input: {id: "%s"}
		) {
			success
			__typename
		}
		}`
	OVERRIDES_READ_QUERY = `{
	dataClassificationOverrideRules {
		results {
		id
		name
		description
		environmentScope {
			environmentIds
			__typename
		}
		behaviorOverride {
			dataSuppressionOverride
			__typename
		}
		spanFilter {
			keyValueFilter {
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
			negatedFilter {
			keyValueFilter {
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
			__typename
			}
			__typename
		}
		__typename
		}
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
