package malicious_sources

const (
	FETCH_REGION_ID = `query Countries {
    countries {
        results {
            id
            name
        }
    }
}`
	DELETE_MALICIOUS_SOURCES = `mutation {
    deleteMaliciousSourcesRule(
        delete: { id: "%s" }
    ) {
        success
        __typename
    }
}`
	CREATE_IP_TYPE_ALERT = `mutation {
    createMaliciousSourcesRule(
        create: {
            ruleInfo: {
                name: "%s"
                description: "%s"
                action: {
                    eventSeverity: %s
                    ruleActionType: %s
                    effects: [%s]
                }
                conditions: [
                    {
                        conditionType: IP_LOCATION_TYPE
                        ipLocationTypeCondition: { ipLocationTypes: [%s] }
                    }
                ]
            }
            %s
        }
    ) {
        id
        __typename
    }
}`
	ENVIRONMENT_SCOPE_QUERY = `scope: { environmentScope: { environmentIds: %s } }`
	UPDATE_IP_TYPE_ALERT    = `mutation {
    updateMaliciousSourcesRule(
        update: {
            rule: {
                id: "%s"
                info: {
                    name: "%s"
                    description: "%s"
                    action: { eventSeverity: %s, ruleActionType: %s, effects: [%s] }
                    conditions: [
                        {
                            conditionType: IP_LOCATION_TYPE
                            ipLocationTypeCondition: { ipLocationTypes: %s }
                        }
                    ]
                }
                status: { disabled: false }
                %s
            }
        }
    ) {
        id
        __typename
    }
}`
	UPDATE_IP_TYPE_BLOCK = `mutation {
    updateMaliciousSourcesRule(
        update: {
            rule: {
                id: "%s"
                info:{
                    name: "%s"
                    description: "%s"
                    action: {
                        eventSeverity: %s
                        ruleActionType: %s
                        %s
                        effects: []
                    }
                    conditions: [
                        {
                            conditionType: IP_LOCATION_TYPE
                            ipLocationTypeCondition: { ipLocationTypes: %s }
                        }
                    ]
                }
                status: { disabled: false }
                %s
            }
        }
    ) {
        id
        __typename
    }
}`

	CREATE_IP_TYPE_BLOCK = `mutation {
    createMaliciousSourcesRule(
        create: {
            ruleInfo: {
                name: "%s"
                description: "%s"
                action: {
                    eventSeverity: %s
                    ruleActionType: %s
                    %s
                    effects: []
                }
                conditions: [
                    {
                        conditionType: IP_LOCATION_TYPE
                        ipLocationTypeCondition: { ipLocationTypes: %s }
                    }
                ]
            }
            %s
        }
    ) {
        id
        __typename
    }
}`

	MALICOUS_SOURCES_READ = `{
  maliciousSourcesRules {
    results {
      id
      info {
        name
        description
        action {
          eventSeverity
          expirationDetails {
            expirationDuration
            expirationTimestamp
            __typename
          }
          ruleActionType
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
        conditions {
          conditionType
          ipLocationTypeCondition {
            ipLocationTypes
            __typename
          }
          emailDomainCondition {
            dataLeakedEmail
            emailRegexes
            emailDomains
            disposableEmailDomain
            emailFraudScore {
              emailFraudScoreType
              minEmailFraudScore
              minEmailFraudScoreLevel
              __typename
            }
            __typename
          }
          __typename
        }
        __typename
      }
      scope {
        environmentScope {
          environmentIds
          __typename
        }
        __typename
      }
      status {
        disabled
        internal
        __typename
      }
      __typename
    }
    __typename
  }
}`
	DELETE_QUERY_IP_RANGE = `mutation {deleteIpRangeRule(delete: {id: "%s"}) {success}}`
	DELETE_REGION         = `mutation {
  deleteRegionRule(input: {id: "%s"}) {
    success
    __typename
  }
}`
	CREATE_IP_RANGE_BLOCK = `mutation {
    createIpRangeRule(
        create: {
            ruleDetails: {
                name: "%s"
                rawIpRangeData: %s
                ruleAction: %s
                description: "%s"
                eventSeverity: %s
				%s
                effects: []
            }
            %s
        }
    ) {
        id
        __typename
    }
}`
	EMAIL_FRAUD_SCORE_QUERY = `emailFraudScore: {
                                emailFraudScoreType: MIN_SEVERITY
                                minEmailFraudScoreLevel: %s
                            }`
	CREATE_EMAIL_DOMAIN_BLOCK = `mutation {
    createMaliciousSourcesRule(
        create: {
            ruleInfo: {
                name: "%s"
                description: "%s"
                action: {
                    eventSeverity: %s
                    ruleActionType: %s
                    %s
                }
                conditions: [
                    {
                        conditionType: EMAIL_DOMAIN
                        emailDomainCondition: {
                            dataLeakedEmail: %t
                            disposableEmailDomain: %t
                            emailDomains: %s
                            emailRegexes: %s
                            %s
                        }
                    }
                ]
            }
            %s
        }
    ) {
        id
        __typename
    }
}`
	CREATE_EMAIL_DOMAIN_ALERT = `mutation {
    createMaliciousSourcesRule(
        create: {
            ruleInfo: {
                name: "%s"
                description: "%s"
                action: {
                    eventSeverity: %s
                    ruleActionType: %s
                }
                conditions: [
                    {
                        conditionType: EMAIL_DOMAIN
                        emailDomainCondition: {
                            dataLeakedEmail: %t
                            disposableEmailDomain: %t
                            emailDomains: %s
                            emailRegexes: %s
                            %s
                        }
                    }
                ]
            }
            %s
        }
    ) {
        id
        __typename
    }
}`
	UPDATE_EMAIL_DOMAIN_BLOCK = `mutation {
    updateMaliciousSourcesRule(
        update: {
                rule: {
                    id: "%s"
                    info: {
                        name: "%s"
                        description: "%s"
                        action: {
                            eventSeverity: %s
                            ruleActionType: %s
                            %s
                        }
                        conditions: [
                            {
                                conditionType: EMAIL_DOMAIN
                                emailDomainCondition: {
                                    dataLeakedEmail: %t
                                    disposableEmailDomain: %t
                                    emailDomains: %s
                                    emailRegexes: %s
                                    %s
                                }
                            }
                        ]
                    }
                    status: { disabled: false }
                    %s
                }
            }
    ) {
        id
        __typename
    }
}`
	UPDATE_EMAIL_DOMAIN_ALERT = `mutation {
    updateMaliciousSourcesRule(
        update: {
            rule: {
                id: "%s"
                info: {
                    name: "%s"
                    description: "%s"
                    action: { eventSeverity: %s, ruleActionType: %s }
                    conditions: [
                    {
                        conditionType: EMAIL_DOMAIN
                        emailDomainCondition: {
                            dataLeakedEmail: %t
                            disposableEmailDomain: %t
                            emailDomains: %s
                            emailRegexes: %s
                            %s
                        }
                    }
                ]
                }
                status: { disabled: true }
                %s
            }
        }
    ) {
        id
        __typename
    }
}`
	CREATE_REGION_RULE_BLOCK = `mutation {
    createRegionRule(
        input: {
            name: "%s"
            regionIds: %s
            type: %s
            description: "%s"
            eventSeverity: %s
			%s
            effects: []
			%s
        }
    ) {
        id
        __typename
    }
}`
	CREATE_REGION_RULE_ALERT = `mutation {
    createRegionRule(
        input: {
            name: "%s"
            regionIds: %s
            type: %s
            description: "%s"
            eventSeverity: %s
            effects: [%s]
			%s
        }
    ) {
        id
        __typename
    }
}`
	UPDATE_REGION_RULE_ALERT = `mutation {
   updateRegionRule(
        input: {
			id: "%s"
            name: "%s"
            regionIds: %s
            type: %s
            description: "%s"
            eventSeverity: %s
            effects: [%s]
			%s
        }
    ) {
        id
        __typename
    }
}`
	UPDATE_REGION_RULE_BLOCK = `mutation {
    updateRegionRule(
        input: {
			id: "%s"
            name: "%s"
            regionIds: %s
            type: %s
            description: "%s"
            eventSeverity: %s
			%s
            effects: []
			%s
        }
    ) {
        id
        __typename
    }
}`
	CREATE_IP_RANGE_ALLOW = `mutation {
    createIpRangeRule(
        create: {
            ruleDetails: {
                name: "%s"
                rawIpRangeData: %s
                ruleAction: %s
                description: "%s"
                effects: []
				%s
            }
			%s
        }
    ) {
        id
        __typename
    }
}`
	CREATE_IP_RANGE_ALERT = `mutation {
    createIpRangeRule(
        create: {
            ruleDetails: {
                name: "%s"
                rawIpRangeData: %s
                ruleAction: %s
                description: "%s"
                eventSeverity: %s
                effects: [%s]
            }
			%s
        }
    ) {
        id
        __typename
    }
}`
	UPDATE_IP_RANGE_ALERT = `mutation {
    updateIpRangeRule(
        update: {
            id: "%s"
            ruleDetails: {
                name: "%s"
                rawIpRangeData: %s
                ruleAction: %s
                description: "%s"
                eventSeverity: %s
                effects: [%s]
            }
			%s
        }
    ) {
        id
        __typename
    }
}`
	UPDATE_IP_RANGE_ALLOW = `mutation {
    updateIpRangeRule(
        update: {
            id: "%s"
            ruleDetails: {
                name: "%s"
                rawIpRangeData: %s
                ruleAction: %s
                description: "%s"
				%s
                effects: []
            }
			%s
        }
    ) {
        id
        __typename
    }
}`

	UPDATE_IP_RANGE_BLOCK = `mutation {
    updateIpRangeRule(
        update: {
            id: "%s"
            ruleDetails: {
                name: "%s"
                rawIpRangeData: %s
                ruleAction: %s
                description: "%s"
                eventSeverity: %s
				%s
                effects: []
            }
            %s
        }
    ) {
        id
        __typename
    }
}`
	REGION_READ_QUERY = `{
    regionRules {
        results {
            id
            name
            type
            description
            disabled
            internal
            eventSeverity
            expiration {
                duration
                timestamp
                __typename
            }
            regions {
                id
                name
                country {
                    isoCode
                    name
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
        __typename
    }
}
`
	IP_RANGE_READ_QUERY = `{
  ipRangeRules {
    results {
      id
      internal
      disabled
      ruleDetails {
        name
        description
        rawIpRangeData
        ruleAction
        eventSeverity
        expirationDetails {
          expirationDuration
          expirationTimestamp
          __typename
        }
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
)
