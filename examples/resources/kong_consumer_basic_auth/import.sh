# Consumer ACLs can be import by consumer id/username:id/username
terraform import kong_consumer.example <consumer_id>:<basic_auth_id>
terraform import kong_consumer.example <consumer_username>:<basic_auth_id>
terraform import kong_consumer.example <consumer_id>:<basic_auth_username>
terraform import kong_consumer.example <consumer_username>:<basic_auth_username>