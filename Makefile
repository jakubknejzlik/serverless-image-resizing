include .env
export

init:
	terraform init
apply:
	terraform apply
destroy:
	terraform destroy

build-lambda-function:
	cd lambda && GO111MODULE=on GOARCH=amd64 GOOS=linux go build -o main main.go && zip ../lambda.zip main && rm main && cd ..

test-thumbnail:
	cd lambda && go run main.go

test:
	curl -I -L --max-redirs 5 https://dpjc5oprjprsh.cloudfront.net/infd-develop-files/24fee41d-85a7-4b91-9b5a-499a7f49e867-157