# Consumer ACLs can be import by consumer id/username:id/key
terraform import kong_consumer.example <consumer_id>:<jwt_id>
terraform import kong_consumer.example <consumer_username>:<jwt_id>
terraform import kong_consumer.example <consumer_id>:<key>
terraform import kong_consumer.example <consumer_username>:<key>