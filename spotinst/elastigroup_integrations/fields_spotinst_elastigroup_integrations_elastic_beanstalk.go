package elastigroup_integrations

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func SetupElasticBeanstalk(fieldsMap map[commons.FieldName]*commons.GenericField) {
	fieldsMap[IntegrationElasticBeanstalk] = commons.NewGenericField(
		commons.ElastigroupIntegrations,
		IntegrationElasticBeanstalk,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(EnvironmentID): {
						Type:     schema.TypeString,
						Required: true,
					},
					string(DeploymentPreferences): {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(AutomaticRoll): {
									Type:     schema.TypeBool,
									Required: true,
								},
								string(BatchSizePercentage): {
									Type:     schema.TypeInt,
									Optional: true,
								},
								string(GracePeriod): {
									Type:     schema.TypeInt,
									Optional: true,
								},
								string(Strategy): {
									Type:     schema.TypeList,
									Required: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											string(Action): {
												Type:     schema.TypeString,
												Required: true,
											},
											string(ShouldDrainInstances): {
												Type:     schema.TypeBool,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			if v, ok := resourceData.GetOk(string(IntegrationElasticBeanstalk)); ok {
				if integration, err := expandAWSGroupElasticBeanstalkIntegration(v); err != nil {
					return err
				} else {
					elastigroup.Integration.SetElasticBeanstalk(integration)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.ElastigroupWrapper)
			elastigroup := egWrapper.GetElastigroup()
			var value *aws.ElasticBeanstalkIntegration = nil

			if v, ok := resourceData.GetOk(string(IntegrationElasticBeanstalk)); ok {
				if integration, err := expandAWSGroupElasticBeanstalkIntegration(v); err != nil {
					return err
				} else {
					value = integration
				}
			}
			elastigroup.Integration.SetElasticBeanstalk(value)
			return nil
		},
		nil,
	)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func expandAWSGroupElasticBeanstalkIntegration(data interface{}) (*aws.ElasticBeanstalkIntegration, error) {
	integration := &aws.ElasticBeanstalkIntegration{}
	list := data.([]interface{})

	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(EnvironmentID)].(string); ok && v != "" {
			integration.SetEnvironmentID(spotinst.String(v))
		}

		if v, ok := m[string(DeploymentPreferences)]; ok {
			deploymentPrefs, err := expandAWSGroupElasticBeanstalkIntegrationDeploymentPreferences(v)

			if err != nil {
				return nil, err
			}

			if deploymentPrefs != nil {
				integration.SetDeploymentPreferences(deploymentPrefs)
			}
		}
	}
	return integration, nil
}

func expandAWSGroupElasticBeanstalkIntegrationDeploymentPreferences(data interface{}) (*aws.DeploymentPreferences, error) {
	if list := data.([]interface{}); len(list) > 0 {
		deploymentPrefs := &aws.DeploymentPreferences{}

		if list != nil && list[0] != nil {
			m := list[0].(map[string]interface{})

			if v, ok := m[string(AutomaticRoll)].(bool); ok {
				deploymentPrefs.SetAutomaticRoll(spotinst.Bool(v))
			} else {
				return nil, errors.New("invalid deployment preferences attributes: set_automatic_roll missing")

			}
			if v, ok := m[string(BatchSizePercentage)].(int); ok && v > 0 {
				deploymentPrefs.SetBatchSizePercentage(spotinst.Int(v))
			}
			if v, ok := m[string(GracePeriod)].(int); ok && v >= 0 {
				deploymentPrefs.SetGracePeriod(spotinst.Int(v))
			}
			if v, ok := m[string(Strategy)]; ok {
				strategy, err := expandAWSGroupElasticBeanstalkIntegrationStrategy(v)
				if err != nil {
					return nil, err
				}
				if strategy != nil {
					deploymentPrefs.SetBeanstalkStrategy(strategy)
				} else {
					return nil, errors.New("invalid deployment preferences attributes: strategy missing")
				}
			}
		}
		return deploymentPrefs, nil
	}
	return nil, nil
}

func expandAWSGroupElasticBeanstalkIntegrationStrategy(data interface{}) (*aws.BeanstalkStrategy, error) {
	if list := data.([]interface{}); len(list) > 0 {
		strategy := &aws.BeanstalkStrategy{}

		if list != nil && list[0] != nil {
			m := list[0].(map[string]interface{})
			if v, ok := m[string(Action)].(string); ok && v != "" {
				strategy.SetAction(spotinst.String(v))
			} else {
				return nil, errors.New("invalid strategy attributes: action missing")
			}
			if v, ok := m[string(ShouldDrainInstances)].(bool); ok {
				strategy.SetShouldDrainInstances(spotinst.Bool(v))
			}
		}
		return strategy, nil
	}
	return nil, nil
}
