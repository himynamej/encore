# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Install dependencies

gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

brew:
	brew update
	brew list encore || brew install encoredev/tap/encore

# ==============================================================================
# Run Project

up:
	encore run -v

GENERATE_ID = $(shell docker ps | grep encoredotdev | cut -c1-12)
SET_ID = $(eval MY_ID=$(GENERATE_ID))

down:
	$(SET_ID)
	docker stop $(MY_ID)
	docker rm $(MY_ID) -v

# ==============================================================================
# Access Project

users:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"

pgcli:
	pgcli $(shell encore db conn-uri url)

curl:
	curl -il "http://127.0.0.1:4000/test?limit=2&offset=2"

# Auth
# export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiI1Y2YzNzI2Ni0zNDczLTQwMDYtOTg0Zi05MzI1MTIyNjc4YjciLCJleHAiOjE3NDE5NzU3NjIsImlhdCI6MTcxMDQzOTc2Miwicm9sZXMiOlsiQURNSU4iXX0.qAhRvfAVtckeqFVkWF5KVMmvWXwh-aY8ffGEEDWtSm79X45f2qqVG4qKz5xL-CbRN1rkpCSOPJxK84ywtVqvl8l55mT89xsQwHYxu8I6EkzMgP4XMUpzL5IFW6FuqPuKDryZ9COMiWPsN1zxFpzQaqJT-CP8XaiB15hGXN9kPQbqYF7ps-eUg6wd0-jLbTPrKuIkDOXL3lgLbXPztRVPxjKeMy3hzs_7KVfoKeqivE7sZT1iI6EpSMwfsQiYVeRCxD-e7tQc3j0kNoXZAfAk2KHKOiq5HOG1eMWAoAJR6sjwKW--igL_aIcXpHx_lOyY6TKRyKkgg1C51URQ1ruVkw

# Unauth
# export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiI1Y2YzNzI3Ni0zNDczLTQwMDYtOTg0Zi05MzI1MTIyNjc4YjciLCJleHAiOjE3NDE5NzU3NjIsImlhdCI6MTcxMDQzOTc2Miwicm9sZXMiOlsiQURNSU4iXX0.qAhRvfAVtckeqFVkWF5KVMmvWXwh-aY8ffGEEDWtSm79X45f2qqVG4qKz5xL-CbRN1rkpCSOPJxK84ywtVqvl8l55mT89xsQwHYxu8I6EkzMgP4XMUpzL5IFW6FuqPuKDryZ9COMiWPsN1zxFpzQaqJT-CP8XaiB15hGXN9kPQbqYF7ps-eUg6wd0-jLbTPrKuIkDOXL3lgLbXPztRVPxjKeMy3hzs_7KVfoKeqivE7sZT1iI6EpSMwfsQiYVeRCxD-e7tQc3j0kNoXZAfAk2KHKOiq5HOG1eMWAoAJR6sjwKW--igL_aIcXpHx_lOyY6TKRyKkgg1C51URQ1ruVkw

auth:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://127.0.0.1:4000/testauth/5cf37266-3473-4006-984f-9325122678b7?limit=2&offset=2"

create:
	curl -il -X POST \
	-d '{"name": "bill", "email": "bill3@ardanlabs.com", "roles": ["ADMIN"], "department": "IT", "password": "123", "passwordConfirm": "123"}' \
	-H "Authorization: Bearer ${TOKEN}" "http://127.0.0.1:4000/v1/users"

token:
	curl -il -X GET \
	--user "admin@example.com:gophers" "http://127.0.0.1:4000/v1/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

update:
	curl -il -X PUT \
	-d '{"name": "jill"}' \
	-H "Authorization: Bearer ${TOKEN}" "http://127.0.0.1:4000/v1/users/6e7bcb19-8389-44a2-9bcf-074d9bcd2bb8"

delete:
	curl -il -X DELETE \
	-H "Authorization: Bearer ${TOKEN}" "http://127.0.0.1:4000/v1/users/6e7bcb19-8389-44a2-9bcf-074d9bcd2bb8"

queryid:
	curl -il -X GET \
	-H "Authorization: Bearer ${TOKEN}" "http://127.0.0.1:4000/v1/users/3a60a6f5-cc42-4f5c-aabd-e86f8ab11057"

query:
	curl -il -X GET \
	-H "Authorization: Bearer ${TOKEN}" "http://127.0.0.1:4000/v1/users?page=1&rows=4"