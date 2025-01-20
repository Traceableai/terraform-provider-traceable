package custom_signature

const (
	BLOCK_UPDATE_QUERY = `mutation {
                                    updateCustomSignatureRule(
                                        update: {
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: { eventType: %s, effects: [] , eventSeverity: %s}
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                            %s
                                            id: "%s"
                                            disabled: %t
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`

	BLOCK_CREATE_QUERY = `mutation {
                                    createCustomSignatureRule(
                                        create: {
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: { eventType: %s, effects: [] ,eventSeverity: %s}
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                            %s
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`

	TEST_UPDATE_QUERY = `mutation {
                                    updateCustomSignatureRule(
                                        update: {
                                            id: "%s"
                                            disabled: %t
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: {
                                                eventType: %s,
                                                effects: [
                                                    %s
                                                ]
                                            }
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`

	TEST_CREATE_QUERY = `mutation {
                                    createCustomSignatureRule(
                                        create: {
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: {
                                                eventType: %s,
                                                effects: [
                                                    %s
                                                ]
                                            }
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`

	ALLOW_UPDATE_QUERY = `mutation {
                                    updateCustomSignatureRule(
                                        update: {
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: { eventType: %s, effects: [] }
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                            %s
                                            id: "%s"
                                            disabled: %t
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`

	ALLOW_CREATE_QUERY = `mutation {
                                    createCustomSignatureRule(
                                        create: {
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: { eventType: %s, effects: [] }
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                            %s
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`

	DELETE_QUERY = `mutation {
                             deleteCustomSignatureRule(delete: {id: "%s"}) {
                               success
                               __typename
                             }
                           }`

	ALERT_UPDATE_QUERY = `mutation {
                                    updateCustomSignatureRule(
                                        update: {
                                            id: "%s"
                                            disabled: %t
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: {
                                                eventType: %s,
                                                effects: [
                                                    %s
                                                ] ,
                                                eventSeverity: %s
                                            }
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`

	READ_QUERY = `{
		customSignatureRules {
		  results {
			id
			name
			description
			disabled
			internal
			blockingExpirationDuration
			blockingExpirationTime
			ruleSource
			ruleEffect {
			  eventType
			  eventSeverity
			  effects {
				ruleEffectType
				agentEffect {
				  agentModifications {
					agentModificationType
					headerInjection {
					  key
					  value
					  headerCategory
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
			ruleDefinition {
			  labels {
				key
				value
				__typename
			  }
			  clauseGroup {
				clauseOperator
				clauses {
				  clauseType
				  matchExpression {
					matchKey
					matchOperator
					matchValue
					matchCategory
					__typename
				  }
				  keyValueExpression {
					keyValueTag
					matchKey
					matchValue
					keyMatchOperator
					valueMatchOperator
					matchCategory
					__typename
				  }
				  attributeKeyValueExpression {
					keyCondition {
					  operator
					  value
					  __typename
					}
					valueCondition {
					  operator
					  value
					  __typename
					}
					__typename
				  }
				  customSecRule {
					inputSecRuleString
					__typename
				  }
				  __typename
				}
				__typename
			  }
			  __typename
			}
			ruleScope {
			  environmentScope {
				environmentIds
				__typename
			  }
			  __typename
			}
			__typename
		  }
		  __typename
		}
	  }`

	ALERT_CREATE_QUERY = `mutation {
                                    createCustomSignatureRule(
                                        create: {
                                            name: "%s"
                                            description: "%s"
                                            ruleEffect: {
                                                eventType: %s,
                                                effects: [
                                                    %s
                                                ] ,
                                                eventSeverity: %s
                                            }
                                            internal: false
                                            ruleDefinition: {
                                                labels: []
                                                clauseGroup: {
                                                    clauseOperator: AND
                                                    clauses: [
                                                        %s
                                                        %s
                                                        %s
                                                    ]
                                                }
                                            }
                                            %s
                                        }
                                    ) {
                                        id
                                        __typename
                                    }
                                }`
	ENVIRONMENT_SCOPE_QUERY = `ruleScope: {
                                   environmentScope: { environmentIds: [%s] }
                               }`
	REQ_RES_CONDITION_QUERY = `{
								clauseType: MATCH_EXPRESSION
								matchExpression: {
									matchKey: %s
									matchCategory: %s
									matchOperator: %s
									matchValue: "%s"
								}
							}`
	ATTRIBUTE_VALUE_CONDITION_QUERY = `valueCondition: { operator: %s, value: "%s" }`
	ATTRIBUTES_BASED_QUERY          = ` {
								clauseType: ATTRIBUTE_KEY_VALUE_EXPRESSION
								attributeKeyValueExpression: {
									keyCondition: { operator: %s, value: "%s" }
									%s
								}
							}`
	CUSTOM_SEC_RULE_QUERY = `{
								clauseType: CUSTOM_SEC_RULE
								customSecRule: {
									inputSecRuleString: "%s"
								}
							}`
	AGENT_EFFECT_QUERY_TEMPLATE = `{
							agentModificationType: HEADER_INJECTION
							headerInjection: {
								key: "%s"
								value: "%s"
								headerCategory: REQUEST
							}
						}`
	CUSTOM_HEADER_INJECTION_QUERY = `{
                                               ruleEffectType: AGENT_EFFECT
                                               agentEffect: {
                                                   agentModifications: [
                                                       %s
                                                   ]
                                               }
                                           }`
)
