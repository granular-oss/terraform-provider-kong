package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccKongPluginDataSource(t *testing.T) {
	corsConfig := `{"credentials":false,"exposed_headers":null,"headers":null,"max_age":null,"methods":["GET","HEAD","PUT","PATCH","POST","DELETE","OPTIONS","TRACE","CONNECT"],"origins":null,"preflight_continue":false,"private_network":false}`
	jwtConfig := `{"anonymous":null,"claims_to_verify":null,"cookie_names":[],"header_names":["authorization"],"key_claim_name":"iss","maximum_expiration":0,"run_on_preflight":true,"secret_is_base64":false,"uri_param_names":["jwt"]}`
	lambdaConfig := `{"aws_assume_role_arn":null,"aws_imds_protocol_version":"v1","aws_key":null,"aws_region":null,"aws_role_session_name":"kong","aws_secret":null,"awsgateway_compatible":false,"base64_encode_body":true,"disable_https":false,"forward_request_body":false,"forward_request_headers":false,"forward_request_method":false,"forward_request_uri":false,"function_name":null,"host":null,"invocation_type":"RequestResponse","is_proxy_integration":false,"keepalive":60000,"log_type":"Tail","port":443,"proxy_url":null,"qualifier":null,"skip_large_bodies":true,"timeout":60000,"unhandled_status":null}`
	servicePluginAsserts := []statecheck.StateCheck{
		assertStringValue("data.kong_plugin.test", "id", "19a479c1-3999-471f-a13b-467421be7e1c"),
		assertStringValue("data.kong_plugin.test", "name", "cors"),
		assertStringValue("data.kong_plugin.test", "service_id", "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"),
		assertBoolValue("data.kong_plugin.test", "enabled", true),
		assertStringValue("data.kong_plugin.test", "config_json", corsConfig),
	}
	routePluginAsserts := []statecheck.StateCheck{
		assertStringValue("data.kong_plugin.test", "id", "6c32041b-9ab7-4ca2-a7a0-7138eb0be9be"),
		assertStringValue("data.kong_plugin.test", "name", "jwt"),
		assertStringValue("data.kong_plugin.test", "route_id", "772e8b14-b020-4164-b52d-3643f6310821"),
		assertBoolValue("data.kong_plugin.test", "enabled", true),
		assertStringValue("data.kong_plugin.test", "config_json", jwtConfig),
	}
	consumerPluginAsserts := []statecheck.StateCheck{
		assertStringValue("data.kong_plugin.test", "id", "2a5e754b-7cb8-4801-8a01-4bd4c7e54a01"),
		assertStringValue("data.kong_plugin.test", "name", "aws-lambda"),
		assertStringValue("data.kong_plugin.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
		assertBoolValue("data.kong_plugin.test", "enabled", true),
		assertStringValue("data.kong_plugin.test", "config_json", lambdaConfig),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	service_id = "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"
	name = "cors"
}`,
				ConfigStateChecks: servicePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	service_name = "test"
	name = "cors"
}`,
				ConfigStateChecks: servicePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	service_id = "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"
	id = "19a479c1-3999-471f-a13b-467421be7e1c"
}`,
				ConfigStateChecks: servicePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	service_name = "test"
	id = "19a479c1-3999-471f-a13b-467421be7e1c"
}`,
				ConfigStateChecks: servicePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	route_id = "772e8b14-b020-4164-b52d-3643f6310821"
	name = "jwt"
}`,
				ConfigStateChecks: routePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	route_name = "test"
	name = "jwt"
}`,
				ConfigStateChecks: routePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	route_id = "772e8b14-b020-4164-b52d-3643f6310821"
	id = "6c32041b-9ab7-4ca2-a7a0-7138eb0be9be"
}`,
				ConfigStateChecks: routePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	route_name = "test"
	id = "6c32041b-9ab7-4ca2-a7a0-7138eb0be9be"
}`,
				ConfigStateChecks: routePluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
	name = "aws-lambda"
}`,
				ConfigStateChecks: consumerPluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	consumer_name = "test"
	name = "aws-lambda"
}`,
				ConfigStateChecks: consumerPluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	id = "2a5e754b-7cb8-4801-8a01-4bd4c7e54a01"
	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
}`,
				ConfigStateChecks: consumerPluginAsserts,
			},
			{
				Config: providerConfig + `
data "kong_plugin" "test" {
	id = "2a5e754b-7cb8-4801-8a01-4bd4c7e54a01"
	consumer_name = "test"
}`,
				ConfigStateChecks: consumerPluginAsserts,
			},
		},
	})
}

func TestAccKongPluginResourceService(t *testing.T) {
	emptyCorsConfig := `{"credentials":false,"exposed_headers":null,"headers":null,"max_age":null,"methods":["GET","HEAD","PUT","PATCH","POST","DELETE","OPTIONS","TRACE","CONNECT"],"origins":null,"preflight_continue":false,"private_network":false}`
	updatedCorsConfig := `{"credentials":false,"exposed_headers":null,"headers":null,"max_age":null,"methods":["GET"],"origins":[".*\\.example2.com"],"preflight_continue":false,"private_network":false}`
	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_plugin" "test" {
 	service_id = kong_service.plugin.id
	name = "cors"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "cors"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_service.plugin", "id", "kong_plugin.test", "service_id"),
					assertStringValue("kong_plugin.test", "config_json", "{}"),
					assertStringValue("kong_plugin.test", "computed_config", emptyCorsConfig),
				},
			},
			{
				ResourceName:            "kong_plugin.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_json", "strict_match"},
			},
			{
				ResourceName:            "kong_plugin.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_json", "strict_match"},
				ImportStateId:           "service:service_plugin:cors",
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_plugin" "test" {
 	service_id = kong_service.plugin.id
	name = "cors"
	config_json = jsonencode(
		{
		origins = [".*\\.example2.com"]
		methods = ["GET"]
		}
	)
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "cors"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_service.plugin", "id", "kong_plugin.test", "service_id"),
					assertStringValue("kong_plugin.test", "config_json", `{"methods":["GET"],"origins":[".*\\.example2.com"]}`),
					assertStringValue("kong_plugin.test", "computed_config", updatedCorsConfig),
				},
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_plugin" "test" {
 	service_id = kong_service.plugin.id
	name = "cors"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "cors"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_service.plugin", "id", "kong_plugin.test", "service_id"),
					assertStringValue("kong_plugin.test", "config_json", "{}"),
					assertStringValue("kong_plugin.test", "computed_config", emptyCorsConfig),
				},
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_plugin" "test" {
 	service_id = kong_service.plugin.id
	name = "jwt"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin2" {
	name = "service_plugin2"
	host = "testing.com"
}
resource "kong_plugin" "test" {
 	service_id = kong_service.plugin2.id
	name = "jwt"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
				},
			},
		},
	})

}

func TestAccKongPluginResourceRoute(t *testing.T) {
	emptyCorsConfig := `{"credentials":false,"exposed_headers":null,"headers":null,"max_age":null,"methods":["GET","HEAD","PUT","PATCH","POST","DELETE","OPTIONS","TRACE","CONNECT"],"origins":null,"preflight_continue":false,"private_network":false}`
	updatedCorsConfig := `{"credentials":false,"exposed_headers":null,"headers":null,"max_age":null,"methods":["GET"],"origins":[".*\\.example2.com"],"preflight_continue":false,"private_network":false}`
	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_route" "plugin" {
	name  = "route_plugin"
	paths = ["/foobar"]
	service_id = kong_service.plugin.id
}
resource "kong_plugin" "test" {
 	route_id = kong_route.plugin.id
	name = "cors"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "cors"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_route.plugin", "id", "kong_plugin.test", "route_id"),
					assertStringValue("kong_plugin.test", "config_json", "{}"),
					assertStringValue("kong_plugin.test", "computed_config", emptyCorsConfig),
				},
			},
			{
				ResourceName:            "kong_plugin.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_json", "strict_match"},
			},
			{
				ResourceName:            "kong_plugin.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_json", "strict_match"},
				ImportStateId:           "route:route_plugin:cors",
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_route" "plugin" {
	name  = "route_plugin"
	paths = ["/foobar"]
	service_id = kong_service.plugin.id
}
resource "kong_plugin" "test" {
 	route_id = kong_route.plugin.id
	name = "cors"
	config_json = jsonencode(
		{
		origins = [".*\\.example2.com"]
		methods = ["GET"]
		}
	)
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "cors"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_route.plugin", "id", "kong_plugin.test", "route_id"),
					assertStringValue("kong_plugin.test", "config_json", `{"methods":["GET"],"origins":[".*\\.example2.com"]}`),
					assertStringValue("kong_plugin.test", "computed_config", updatedCorsConfig),
				},
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_route" "plugin" {
	name  = "route_plugin"
	paths = ["/foobar"]
	service_id = kong_service.plugin.id
}
resource "kong_plugin" "test" {
 	route_id = kong_route.plugin.id
	name = "cors"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "cors"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_route.plugin", "id", "kong_plugin.test", "route_id"),
					assertStringValue("kong_plugin.test", "config_json", "{}"),
					assertStringValue("kong_plugin.test", "computed_config", emptyCorsConfig),
				},
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_route" "plugin" {
	name  = "route_plugin"
	paths = ["/foobar"]
	service_id = kong_service.plugin.id
}
resource "kong_plugin" "test" {
 	route_id = kong_route.plugin.id
	name = "jwt"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: providerConfig + `
resource "kong_service" "plugin" {
	name = "service_plugin"
	host = "testing.com"
}
resource "kong_route" "plugin2" {
	name  = "route_plugin2"
	paths = ["/foobar"]
	service_id = kong_service.plugin.id
}
resource "kong_plugin" "test" {
 	route_id = kong_route.plugin2.id
	name = "jwt"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
				},
			},
		},
	})

}
func TestAccKongPluginResourceConsumer(t *testing.T) {
	emptyLambdaConfig := `{"aws_assume_role_arn":null,"aws_imds_protocol_version":"v1","aws_key":null,"aws_region":null,"aws_role_session_name":"kong","aws_secret":null,"awsgateway_compatible":false,"base64_encode_body":true,"disable_https":false,"forward_request_body":false,"forward_request_headers":false,"forward_request_method":false,"forward_request_uri":false,"function_name":null,"host":null,"invocation_type":"RequestResponse","is_proxy_integration":false,"keepalive":60000,"log_type":"Tail","port":443,"proxy_url":null,"qualifier":null,"skip_large_bodies":true,"timeout":60000,"unhandled_status":null}`
	updatedLambdaConfig := `{"aws_assume_role_arn":null,"aws_imds_protocol_version":"v1","aws_key":null,"aws_region":null,"aws_role_session_name":"kong","aws_secret":null,"awsgateway_compatible":false,"base64_encode_body":true,"disable_https":false,"forward_request_body":false,"forward_request_headers":false,"forward_request_method":false,"forward_request_uri":false,"function_name":"test","host":null,"invocation_type":"RequestResponse","is_proxy_integration":false,"keepalive":60000,"log_type":"Tail","port":443,"proxy_url":null,"qualifier":null,"skip_large_bodies":true,"timeout":60000,"unhandled_status":null}`
	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_consumer" "plugin" {
	username = "consumer_plugin"
}
resource "kong_plugin" "test" {
 	consumer_id = kong_consumer.plugin.id
	name = "aws-lambda"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "aws-lambda"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_consumer.plugin", "id", "kong_plugin.test", "consumer_id"),
					assertStringValue("kong_plugin.test", "config_json", "{}"),
					assertStringValue("kong_plugin.test", "computed_config", emptyLambdaConfig),
				},
			},
			{
				ResourceName:            "kong_plugin.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_json", "strict_match"},
			},
			{
				ResourceName:            "kong_plugin.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_json", "strict_match"},
				ImportStateId:           "consumer:consumer_plugin:aws-lambda",
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "plugin" {
	username = "consumer_plugin"
}
resource "kong_plugin" "test" {
 	consumer_id = kong_consumer.plugin.id
	name = "aws-lambda"
	config_json = jsonencode(
		{
			function_name = "test"
		}
	)
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "aws-lambda"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_consumer.plugin", "id", "kong_plugin.test", "consumer_id"),
					assertStringValue("kong_plugin.test", "config_json", `{"function_name":"test"}`),
					assertStringValue("kong_plugin.test", "computed_config", updatedLambdaConfig),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "plugin" {
	username = "consumer_plugin"
}
resource "kong_plugin" "test" {
 	consumer_id = kong_consumer.plugin.id
	name = "aws-lambda"
	config_json = "{}"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertNotNull("kong_plugin.test", "id"),
					assertStringValue("kong_plugin.test", "name", "aws-lambda"),
					assertBoolValue("kong_plugin.test", "enabled", true),
					assertNullValue("kong_plugin.test", "tags"),
					assertStateMatch("kong_consumer.plugin", "id", "kong_plugin.test", "consumer_id"),
					assertStringValue("kong_plugin.test", "config_json", "{}"),
					assertStringValue("kong_plugin.test", "computed_config", emptyLambdaConfig),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "plugin" {
	username = "consumer_plugin"
}
resource "kong_plugin" "test" {
 	consumer_id = kong_consumer.plugin.id
	name = "azure-functions"
	config_json = jsonencode({
		appname = "test"
		functionname = "test"
	})
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertStateMatch("kong_consumer.plugin", "id", "kong_plugin.test", "consumer_id"),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "plugin2" {
	username = "consumer_plugin2"
}
resource "kong_plugin" "test" {
 	consumer_id = kong_consumer.plugin2.id
	name = "azure-functions"
	config_json = jsonencode({
		appname = "test"
		functionname = "test"
	})
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue(
						"kong_plugin.test", tfjsonpath.New("id"),
					),
					assertStateMatch("kong_consumer.plugin2", "id", "kong_plugin.test", "consumer_id"),
				},
			},
		},
	})

}
