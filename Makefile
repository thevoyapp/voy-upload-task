UPLOAD_DIR := upload
PROFILE := default
REGION := us-east-2
ACCOUNT := 076279718063
CONTAINER := voy-upload-task
REPO_NAME := voy/upload-task
IMAGE_URL := $(ACCOUNT).dkr.ecr.$(REGION).amazonaws.com/$(REPO_NAME)
PARAMETER_NAME := /ecs/$(CONTAINER)/version
VERSION := v0.0.1

build:
	docker build -t $(CONTAINER) .

run: build
	docker run -p 80:80 $(CONTAINER)

update: build
	`aws ecr get-login --region us-east-2 --no-include-email --profile default`
	docker tag $(CONTAINER) $(IMAGE_URL)
	docker push $(IMAGE_URL)

deploy-repo:
	aws cloudformation deploy \
		--template-file repo.yml \
		--stack-name "repo-$(CONTAINER)" \
		--s3-bucket "cf-$(REGION)-$(ACCOUNT)-bucket" \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides "RepositoryName=$(REPO_NAME)" \
			"ContainerName=$(CONTAINER)" \
			"ParameterName=$(PARAMETER_NAME)" \
		--s3-prefix cluster/ \
		--profile $(PROFILE) \
		--region $(REGION)

deploy-service:
		aws ssm put-parameter --name $(PARAMETER_NAME)/temp --type String --value $(VERSION) \
			--overwrite --profile $(PROFILE) --region $(REGION)
		aws cloudformation deploy \
			--template-file service.yml \
			--stack-name "$(CONTAINER)-service" \
			--s3-bucket "cf-$(REGION)-$(ACCOUNT)-bucket" \
			--parameter-overrides "StackName=voy-app-cluster" \
				"Version=$(VERSION)" \
				"ParameterPath=$(PARAMETER_NAME)" \
				"ImageUrl=$(IMAGE_URL)" \
				"ServiceName=$(CONTAINER)" \
				"Priority=1" \
				"DesiredCount=1" \
				"Path=/content/*" \
				"Role=arn:aws:iam::$(ACCOUNT):role/cluster/$(CONTAINER)-Role" \
				"ContainerPort=80" \
			--capabilities CAPABILITY_IAM \
			--s3-prefix cluster/ \
			--profile $(PROFILE) \
			--region $(REGION)

remove-old:
	aws ssm put-parameter --name $(PARAMETER_NAME) --type String --value $(VERSION) \
		--overwrite --profile $(PROFILE) --region $(REGION)
