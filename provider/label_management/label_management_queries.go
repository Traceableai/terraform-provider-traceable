package label_management

const (
	CREATE_LABEL_QUERY      = `mutation{createLabel(label:{key:"%s",description:"%s",color:"%s"}){id key description color}}`
	READ_LABELS_QUERY       = `query{labels{results{id key description color}}}`
	UPDATE_LEBEL_QUERY      = `mutation{updateLabel(label:{id:"%s",key:"%s",description:"%s",color:"%s"}){id key description color}}`
	CREATE_LABEL_RULE_QUERY = `mutation {
		createLabelApplicationRule(
			labelApplicationRuleData: {
				name: "%s"
				description: "%s"
				enabled: %t
				conditionList: %s
				%s
			}
		) {
			id
		}
	}`
	READ_LABEL_RULE_QUERY = `{
		labelApplicationRules {
			results {
				id
				labelApplicationRuleData {
					name
					description
					enabled
					conditionList {
						leafCondition {
							keyCondition {
								operator
								value
							}
							valueCondition {
								valueConditionType
								stringCondition {
									value
									operator
									stringConditionValueType
									values
								}
								unaryCondition {
									operator
								}
							}
						}
					}
					action {
						entityTypes
						staticLabels {
							ids
						}
						operation
						type
						dynamicLabelKey
					}
				}
			}
		}
	}`
	LABEL_RULE_UPDATE_QUERY = `mutation {
		updateLabelApplicationRule(
			labelApplicationRule: {
				id: "%s"
				labelApplicationRuleData: {
					name: "%s"
					description: "%s"
					enabled: %t
					conditionList: %s
					%s
				}
			}
		) {
			id
		}
	}`
)
