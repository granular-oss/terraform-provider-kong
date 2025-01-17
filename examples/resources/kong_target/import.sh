# Upstreams can be import be {upstream name/id}|{target id/target}
terraform import kong_upstream.example "f20ebf48-b96f-4dfa-951d-b2de304aa071|f20ebf48-b96f-4dfa-951d-b2de304aa071"
terraform import kong_upstream.example "example-upstream|example.com:8000"