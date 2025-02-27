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
    count
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
              __typename
            }
            valueCondition {
              valueConditionType
              stringCondition {
                value
                operator
                stringConditionValueType
                values
                __typename
              }
              unaryCondition {
                operator
                __typename
              }
              __typename
            }
            __typename
          }
          __typename
        }
        action {
          entityTypes
          staticLabels {
            ids
            __typename
          }
          operation
          type
          dynamicLabelKey
          dynamicLabel {
            expression
            tokenExtractionRules {
              alias
              key
              regexCapture
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
