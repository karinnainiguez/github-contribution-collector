AWS_REGION = us-west-2
OUTPUT_CFN_STACK = "github-contribution-collector"
OUTPUT_BUCKET = "github-contribution-collector-cfn"

# This can optionally be specified if the deployer wants to use a different AWS
# profile from their config. If not specified the default profile is used.
# 
# make AWS_PROFILE=username deploy
#
AWS_PROFILE = 

# These must be defined by the deployer when running make deploy which can be
# done by running the command like this:
# 
# make SES_VERIFIED_EMAIL=some@example.com GITHUB_API_KEY=key deploy
#
SES_VERIFIED_EMAIL =
GITHUB_API_KEY =

# Setup the AWS command with the proper region and profile
# Run this with `make AWS_PROFILE=whatever` to set the profile at runtime
AWS = aws --region=$(AWS_REGION)
ifdef AWS_PROFILE
	AWS = aws --region=$(AWS_REGION) --profile=$(AWS_PROFILE)
endif

github-contribution-collector: $(wildcard *.go)
	go build

github-contribution-collector.zip: github-contribution-collector
	zip github-contribution-collector.zip github-contribution-collector

# Test if the bucket exists by querying S3 and create it if not
.account-configured:
	@set -eux; \
	if `$(AWS) s3api list-buckets --query="Buckets[].Name | contains(@, '$(OUTPUT_BUCKET)')")` == "true"; then \
		touch .account-configured; \
	else \
		$(AWS) s3api create-bucket \
			--create-bucket-configuration LocationConstraint=$(AWS_REGION) \
			--acl private \
			--bucket $(OUTPUT_BUCKET) \
		;\
		$(AWS) s3api put-bucket-versioning \
			--bucket $(OUTPUT_BUCKET) \
			--versioning-configuration Status=Enabled \
		;\
		$(AWS) s3api put-bucket-encryption \
			--bucket $(OUTPUT_BUCKET) \
			--server-side-encryption-configuration \
				'{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"aws:kms"}}]}' \
		;\
		touch .account-configured; \
	fi

# Transform the template to replace zip references and upload artifacts to S3
packaged-template.yaml: .account-configured github-contribution-collector.zip cfn-deploy.yaml
	$(AWS) cloudformation package \
		--template-file cfn-deploy.yaml \
		--s3-bucket $(OUTPUT_BUCKET) \
		--output-template-file packaged-template.yaml

# Deploy the previously uploaded artifacts using CloudFormation
deploy: packaged-template.yaml
	@if test -z "$(SES_VERIFIED_EMAIL)" ; then \
		echo "SES_VERIFIED_EMAIL is undefined"; \
		exit 1; \
	fi

	@if test -z "$(GITHUB_API_KEY)"; then \
		echo "GITHUB_API_KEY is undefined"; \
		exit 1; \
	fi

	$(AWS) cloudformation deploy \
		--stack-name $(OUTPUT_CFN_STACK) \
		--template-file ./packaged-template.yaml \
		--s3-bucket $(OUTPUT_BUCKET) \
		--force-upload \
		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
		--parameter-overrides GitHubAPIKey=$(GITHUB_API_KEY) SESVerifiedEmail=$(SES_VERIFIED_EMAIL)

clean:
	rm -f \
		github-contribution-collector \
		github-contribution-collector.zip \
		packaged-template.yaml
