package rate_limiting

const (
	RATE_LIMIT_QUERY_KEY       = "ENDPOINT_RATE_LIMITING"
	RATE_LIMITING_CREATE_QUERY = `mutation {
                                      createRateLimitingRule(
                                          rateLimitingRuleData: {
                                              category: %s
                                              conditions: [
                                                  %s
                                              ]
                                              enabled: %t
                                              name: "%s"
                                              thresholdActionConfigs: [
                                                  {
                                                      actions: [{
                                                              actionType: %s
                                                              %s: %s
                                                          }]
                                                      thresholdConfigs: [
                                                             %s
                                                      ]
                                                  }
                                              ]
                                              %s
                                              description: "%s"
                                              labels: []
                                              ruleStatus: { internal: false }
                                          }
                                      ) {
                                          id
                                          name
                                      }
                                  }`

	RATE_LIMITING_UPDATE_QUERY = `mutation {
                                    updateRateLimitingRule(
                                        rateLimitingRule: {
                                            id: "%s"
                                            category: %s
                                            conditions: [
                                                %s
                                            ]
                                            enabled: %t
                                            name: "%s"
                                            thresholdActionConfigs: [
                                                {
                                                    actions: [{
                                                            actionType: %s
                                                            %s: %s
                                                        }]
                                                    thresholdConfigs: [
                                                           %s
                                                    ]
                                                }
                                            ]
                                            %s
                                            description: "%s"
                                            labels: []
                                            ruleStatus: { internal: false }
                                        }
                                    ) {
                                        id
                                        name
                                    }
                                }`

	DELETE_RATE_LIMIT_QUERY = `mutation {
                                    deleteRateLimitingRule(id: "%s") {
                                        success
                                        __typename
                                    }
                                }`

	FETCH_RATE_LIMIT_RULES = `{
  rateLimitingRules: rateLimitingRules(category: %s) {
    results {
      category
      conditions {
        leafCondition {
          conditionType
          datatypeCondition {
            datasetIds
            datatypeIds
            dataLocation
            datatypeMatching {
              datatypeMatchingType
              regexBasedMatching {
                customMatchingLocation {
                  metadataType
                  keyCondition {
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
          ipAddressCondition {
            cidrIpRanges
            ipAddresses
            rawInputIpData
            exclude
            ipAddressConditionType
            __typename
          }
          ipReputationCondition {
            minIpReputationSeverity
            __typename
          }
          ipLocationTypeCondition {
            ipLocationTypes
            exclude
            __typename
          }
          emailDomainCondition {
            emailRegexes
            exclude
            __typename
          }
          ipOrganisationCondition {
            ipOrganisationRegexes
            exclude
            __typename
          }
          ipAsnCondition {
            ipAsnRegexes
            exclude
            __typename
          }
          userIdCondition {
            userIdRegexes
            userIds
            exclude
            __typename
          }
          userAgentCondition {
            userAgentRegexes
            exclude
            __typename
          }
          ipConnectionTypeCondition {
            ipConnectionTypes
            exclude
            __typename
          }
          ipAbuseVelocityCondition {
            minIpAbuseVelocity
            __typename
          }
          keyValueCondition {
            keyCondition {
              operator
              value
              __typename
            }
            metadataType
            valueCondition {
              operator
              value
              __typename
            }
            __typename
          }
          regionCondition {
            regionIdentifiers {
              countryIsoCode
              __typename
            }
            exclude
            __typename
          }
          requestScannerTypeCondition {
            scannerTypes
            exclude
            __typename
          }
          scopeCondition {
            entityScope {
              entityIds
              entityType
              __typename
            }
            labelScope {
              labelIds
              labelType
              __typename
            }
            scopeType
            urlScope {
              urlRegexes
              __typename
            }
            __typename
          }
          __typename
        }
        __typename
      }
      description
      enabled
      id
      name
      labels {
        key
        value
        __typename
      }
      thresholdActionConfigs {
        actions {
          actionType
          alert {
            eventSeverity
            __typename
          }
          block {
            duration
            eventSeverity
            __typename
          }
          __typename
        }
        thresholdConfigs {
          apiAggregateType
          rollingWindowThresholdConfig {
            countAllowed
            duration
            __typename
          }
          dynamicThresholdConfig {
            percentageExceedingMeanAllowed
            meanCalculationDuration
            duration
            __typename
          }
          valueBasedThresholdConfig {
            valueType
            uniqueValuesAllowed
            duration
            sensitiveParamsEvaluationType
            __typename
          }
          thresholdConfigType
          userAggregateType
          __typename
        }
        __typename
      }
      transactionActionConfigs {
        action {
          actionType
          alert {
            eventSeverity
            __typename
          }
          allow {
            duration
            __typename
          }
          block {
            duration
            eventSeverity
            __typename
          }
          __typename
        }
        expirationTimestamp
        __typename
      }
      ruleConfigScope {
        environmentScope {
          environmentIds
          __typename
        }
        __typename
      }
      ruleStatus {
        internal
        ruleCreationSource
        __typename
      }
      __typename
    }
    __typename
  }
}`
	ROLLING_WINDOW_THRESHOLD_CONFIG_QUERY = `{
                                                  apiAggregateType: %s
                                                  userAggregateType: %s
                                                  thresholdConfigType: ROLLING_WINDOW
                                                  rollingWindowThresholdConfig: { countAllowed: %d, duration: "%s" }
                                              }`
	ENUMERATION_THRESHOLD_CONFIG_QUERY = `{
                            apiAggregateType: %s
                            userAggregateType: %s
                            thresholdConfigType: VALUE_BASED
                            valueBasedThresholdConfig: {
                                valueType: %s
                                uniqueValuesAllowed: %d
                                duration: "%s"
                            }
                        }`
	ENUMERATION_THRESHOLD_CONFIG_SENSITIVE_PARAM_QUERY = `{
                            apiAggregateType: %s
                            userAggregateType: %s
                            thresholdConfigType: VALUE_BASED
                            valueBasedThresholdConfig: {
                                valueType: %s
                                uniqueValuesAllowed: %d
                                duration: "%s"
                                sensitiveParamsEvaluationType: %s
                            }
                        }`
	DYNAMIC_THRESHOLD_CONFIG_QUERY = `{
                                           apiAggregateType: %s
                                           userAggregateType: %s
                                           thresholdConfigType: DYNAMIC
                                           dynamicThresholdConfig: {
                                               percentageExceedingMeanAllowed: %d
                                               meanCalculationDuration: "%s"
                                               duration: "%s"
                                           }
                                       }`
	USER_ID_REGEXES_QUERY = `{
                                  leafCondition: {
                                      conditionType: USER_ID
                                      userIdCondition: {
                                          userIdRegexes: %s
                                          exclude: %t
                                      }
                                  }
                              }`
	USER_ID_LIST_QUERY = `{
                                  leafCondition: {
                                      conditionType: USER_ID
                                      userIdCondition: {
                                          userIds: %s
                                          exclude: %t
                                      }
                                  }
                              }`

	REQUEST_SCANNER_TYPE_QUERY = `{
                                 leafCondition: {
                                     conditionType: REQUEST_SCANNER_TYPE
                                     requestScannerTypeCondition: {
                                         scannerTypes: %s
                                         exclude: %t
                                     }
                                 }
                             }`

	IP_CONNECTION_TYPE_QUERY = `{
                                     leafCondition: {
                                         conditionType: IP_CONNECTION_TYPE
                                         ipConnectionTypeCondition: {
                                             ipConnectionTypes: %s
                                             exclude: %t
                                         }
                                     }
                                 }`

	IS_ASN_QUERY = `{
                        leafCondition: {
                            conditionType: IP_ASN
                            ipAsnCondition: { ipAsnRegexes: %s, exclude: %t }
                        }
                    }`

	IP_ORGANISATION_QUERY = `{
                         leafCondition: {
                             conditionType: IP_ORGANISATION
                             ipOrganisationCondition: {
                                 ipOrganisationRegexes: %s
                                 exclude: %t
                             }
                         }
                     }`

	REGION_QUERY = `{
                       leafCondition: {
                           conditionType: REGION
                           regionCondition: {
                               exclude: %t
                               regionIdentifiers: [%s]
                           }
                       }
                   }`

	USER_AGENT_QUERY = `{
                            leafCondition: {
                                conditionType: USER_AGENT
                                userAgentCondition: { userAgentRegexes: %s, exclude: %t }
                            }
                        }`

	RAW_INPUT_IP_ADDRESS_QUERY = `{
                                     leafCondition: {
                                         conditionType: IP_ADDRESS
                                         ipAddressCondition: {
                                             rawInputIpData: %s
                                             exclude: %t
                                         }
                                     }
                                 }`

	ALL_EXTERNAL_IP_ADDRESS_QUERY = `{
                                        leafCondition: {
                                            conditionType: IP_ADDRESS
                                            ipAddressCondition: {
                                                ipAddressConditionType: %s
                                                exclude: %t
                                            }
                                        }
                                    }`

	EMAIL_DOMAIN_QUERY = `{
                           leafCondition: {
                               conditionType: EMAIL_DOMAIN
                               emailDomainCondition: { emailRegexes: %s, exclude: %t }
                           }
                       }`

	IP_LOCATION_TYPE_QUERY = `{
                                 leafCondition: {
                                     conditionType: IP_LOCATION_TYPE
                                     ipLocationTypeCondition: {
                                         ipLocationTypes: %s
                                         exclude: %t
                                     }
                                 }
                             }`

	IP_ABUSE_VELOCITY_QUERY = `{
                                    leafCondition: {
                                        conditionType: IP_ABUSE_VELOCITY
                                        ipAbuseVelocityCondition: { minIpAbuseVelocity: %s }
                                    }
                                }`

	IP_REPUTATION_QUERY = `{
                              leafCondition: {
                                  conditionType: IP_REPUTATION
                                  ipReputationCondition: { minIpReputationSeverity: %s }
                              }
                          }`

	ATTRIBUTE_BASED_CONDITIONS_QUERY = `{
                                                leafCondition: {
                                                    conditionType: KEY_VALUE
                                                    keyValueCondition: {
                                                        keyCondition: { operator: %s, value: "%s" }
                                                        metadataType: TAG
                                                        %s
                                                    }
                                                }
                                            }`

	REQ_RES_CONDITIONS_QUERY = `{
                                   leafCondition: {
                                       conditionType: KEY_VALUE
                                       keyValueCondition: {
                                           metadataType: %s
                                           valueCondition: { operator: %s, value: "%s" }
                                       }
                                   }
                               }`

	DATA_LOCATION_STRING        = `dataLocation: %s`
	DATA_TYPES_CONDITIONS_QUERY = `{
                    leafCondition: {
                        conditionType: DATATYPE
                        datatypeCondition: {
                            datasetIds: %s
                            datatypeIds: %s
                            %s
                        }
                    }
                }`

	ENDPOINT_SCOPED_QUERY = `{
                                  leafCondition: {
                                      conditionType: SCOPE
                                      scopeCondition: {
                                          scopeType: ENTITY
                                          entityScope: {
                                              entityIds: %s
                                              entityType: API
                                          }
                                      }
                                  }
                              }`

	LABEL_ID_SCOPED_QUERY = `{
                              leafCondition: {
                                  conditionType: SCOPE
                                  scopeCondition: {
                                      scopeType: LABEL
                                      labelScope: {
                                          labelIds: %s
                                          labelType: API
                                      }
                                  }
                              }
                          }`
	ENVIRONMENT_SCOPE_QUERY = `ruleConfigScope: {
                                                    environmentScope: {
                                                        environmentIds: %s
                                                    }
                                                }`
)
