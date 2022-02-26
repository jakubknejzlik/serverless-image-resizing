include .env
export

sso:
	aws sso login
init:
	terraform init
apply:
	terraform apply
destroy:
	terraform destroy

build-lambda-function:
	cd lambda && GO111MODULE=on GOARCH=amd64 GOOS=linux go build -o main main.go && cp -R ../icons ./icons && zip ../lambda.zip main ./icons/* && rm main && rm -r ./icons && cd ..

test-thumbnail:
	cd lambda && go run main.go

test:
	curl -I -L --max-redirs 5 $$TEST_URL