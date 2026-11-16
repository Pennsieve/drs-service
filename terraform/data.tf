data "aws_caller_identity" "current" {}

data "aws_region" "current_region" {}

# Import Account Data
data "terraform_remote_state" "account" {
  backend = "s3"

  config = {
    bucket = "${var.aws_account}-terraform-state"
    key    = "aws/terraform.tfstate"
    region = "us-east-1"
  }
}

# Import Region Data
data "terraform_remote_state" "region" {
  backend = "s3"

  config = {
    bucket = "${var.aws_account}-terraform-state"
    key    = "aws/${data.aws_region.current_region.name}/terraform.tfstate"
    region = "us-east-1"
  }
}

# Import VPC Data
data "terraform_remote_state" "vpc" {
  backend = "s3"

  config = {
    bucket  = "${var.aws_account}-terraform-state"
    key     = "aws/${data.aws_region.current_region.name}/${var.vpc_name}/terraform.tfstate"
    region  = "us-east-1"
    profile = var.aws_account
  }
}

# Import Platform Infrastructure Data
data "terraform_remote_state" "platform_infrastructure" {
  backend = "s3"

  config = {
    bucket  = "${var.aws_account}-terraform-state"
    key     = "aws/${data.aws_region.current_region.name}/${var.vpc_name}/${var.environment_name}/platform-infrastructure/terraform.tfstate"
    region  = "us-east-1"
    profile = var.aws_account
  }
}

# Import AWS Default SecretsManager KMS Key
data "aws_kms_key" "ssm_kms_key" {
  key_id = "alias/aws/secretsmanager"
}

# Import API Gateway authorizer
# The shared gateway authorizer lambda now lives in auth-service's Terraform state
# (moved from pennsieve-go-api as part of the Phase 0 authorizer consolidation / F2).
# auth-service exports the same output names (authorizer_lambda_invoke_uri,
# authorizer_invocation_role), so the reads in gateway.tf are unchanged.
data "terraform_remote_state" "api_gateway" {
  backend = "s3"

  config = {
    bucket  = "${var.aws_account}-terraform-state"
    key     = "aws/${data.aws_region.current_region.name}/${var.vpc_name}/${var.environment_name}/auth-service/terraform.tfstate"
    region  = "us-east-1"
    profile = var.aws_account
  }
}

# Import Postgres
data "terraform_remote_state" "pennsieve_postgres" {
  backend = "s3"

  config = {
    bucket  = "${var.aws_account}-terraform-state"
    key     = "aws/${data.aws_region.current_region.name}/${var.vpc_name}/${var.environment_name}/pennsieve-postgres/terraform.tfstate"
    region  = "us-east-1"
    profile = var.aws_account
  }
}
