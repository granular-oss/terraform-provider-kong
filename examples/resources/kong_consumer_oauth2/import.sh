# Consumer ACLs can be import by consumer id/username:id
terraform import kong_consumer_oauth2.example <consumer_id>:<oauth_id>
terraform import kong_consumer_oauth2.example <consumer_username>:<oauth_id>