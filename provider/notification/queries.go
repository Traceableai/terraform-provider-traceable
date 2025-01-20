package notification

const (
	DELETE_NOTIFICATION_RULE = `mutation {
		deleteNotificationRule(input: {ruleId: "%s"}) {
		  success
		}
	  }`
	NOTIFICATION_CHANNEL_READ = `{
		notificationChannels {
			results {
				channelId
				channelName
				notificationChannelConfig {
					emailChannelConfigs {
						address
					}
					customWebhookChannelConfigs {
						headers {
							key
							value
							isSecret
						}
						customWebhookChannelIndex: id
						url
					}
					slackWebhookChannelConfigs {
						url
					}
					splunkIntegrationChannelConfigs {
						splunkIntegrationId
					}
					syslogIntegrationChannelConfigs {
						syslogIntegrationId
					}
					s3BucketChannelConfigs {
						bucketName
						region
						authenticationCredentialType
						webIdentityAuthenticationCredential {
							roleArn	
						}	
					}	
				}	
			}
		}
	}`
	NOTIFICATION_RULE_READ = `{
  notificationRules {
    results {
      ruleId
      ruleName
      environmentScope {
        environments
        __typename
      }
      channelId
      integrationTarget {
        type
        integrationId
        __typename
      }
      category
      eventConditions {
        detectedSecurityEventCondition {
          detectedThreatActivityConditions {
            detectedThreatActivityConditionType
            customDetectionCondition {
              customDetectionType
              __typename
            }
            preDefinedDetectionCondition {
              anomalyRuleId
              __typename
            }
            __typename
          }
          severities
          impactLevels
          confidenceLevels
          __typename
        }
        notificationConfigChangeCondition {
          notificationConfigChangeTypes
          notificationTypes
          __typename
        }
        sensitiveDataConfigChangeCondition {
          sensitiveDataConfigChangeTypes
          sensitiveDataConfigTypes
          __typename
        }
        dataCollectionChangeCondition {
          agentType
          agentActivityType
          agentStatusChanges
          __typename
        }
        blockedEventCondition {
          blockedThreatActivityConditions {
            blockedThreatActivityConditionType
            customBlockingCondition {
              customBlockingType
              __typename
            }
            preDefinedBlockingCondition {
              anomalyRuleId
              __typename
            }
            __typename
          }
          __typename
        }
        threatActorStateChangeEventCondition {
          actorStates
          __typename
        }
        actorSeverityStateChangeEventCondition {
          actorSeverities
          actorIpReputationLevels
          __typename
        }
        securityConfigChangeEventCondition {
          securityConfigurationTypes
          __typename
        }
        userChangeEventCondition {
          userChangeTypes
          __typename
        }
        postureEventCondition {
          postureEvents
          riskDeltas
          __typename
        }
        apiNamingRuleConfigChangeEventCondition {
          apiNamingRuleConfigChangeTypes
          __typename
        }
        apiSpecConfigChangeEventCondition {
          apiSpecConfigChangeTypes
          __typename
        }
        excludeSpanRuleConfigChangeEventCondition {
          excludeSpanRuleConfigChangeTypes
          __typename
        }
        labelConfigChangeEventCondition {
          labelConfigChangeTypes
          labelConfigTypes
          __typename
        }
        riskScoringConfigChangeEventCondition {
          riskScoringConfigChangeTypes
          __typename
        }
        integrationChangeEventCondition {
          integrationNotificationEventTypes
          integrationNotificationTypes
          __typename
        }
        threatScoringConfigChangeEventCondition {
          threatScoringConfigTypes
          __typename
        }
        fraudDetectionEventCondition {
          severities
          __typename
        }
        __typename
      }
      rateLimitIntervalDuration
      __typename
    }
    __typename
  }
}`
	DELETE_NOTIFICATION_CHANNEL = `mutation {
		deleteNotificationChannel(
		  input: {channelId: "%s"}
		) {
		  success
		  
		}
	  }`
	UPDATE_NOTIFICATION_CHANNEL = `mutation {
		updateNotificationChannel(
			input: {
				channelId: "%s"
				channelName: "%s"
				notificationChannelConfig: {
					%s
					%s
					%s
					%s
					%s
					%s
				}
			}
		) {
			channelId
		}
	}`
	CREATE_NOTIFICATION_CHANNEL = `mutation {
		createNotificationChannel(
			input: {
				channelName: "%s"
				notificationChannelConfig: {
					%s
					%s
					%s
					%s
					%s
					%s
				}
			}
		) {
			channelId
		}
	}`
)
