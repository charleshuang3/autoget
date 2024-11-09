please help me outline what config file field the backend needed. gen example yaml

-------------------------

based on backend/example/config.example.yaml gen a lib to read parse validate config in backend.

auth type should be google_auth or none or api_key
if api_key, auth.api_key should be set
if google_auth, auth.google_auth should be set and all child field

-------------------------

in validate, all field should be checked has value

-------------------------

add tests for config.go

-------------------------

for backend/config/config_test.go please also cover LoadConfig()
