# Consumer ACLs can be import by id
terraform import kong_plugin.example f20ebf48-b96f-4dfa-951d-b2de304aa071
# By service id/name and plugin name
terraform import kong_plugin.example service:service_id:plugin_name
# By route id/name and plugin name
terraform import kong_plugin.example route:route_id:plugin_name
# By consumer id/name and plugin name
terraform import kong_plugin.example consumer:consumer_id:plugin_name