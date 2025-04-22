variable "aws_account" {}

variable "aws_region" {}

variable "environment_name" {}

variable "service_name" {}

variable "vpc_name" {}

variable "domain_name" {}

variable "image_tag" {}

variable "lambda_bucket" {
  default = "pennsieve-cc-lambda-functions-use1"
}

variable "api_domain_name" {}

variable "rds_proxy_user" {}

locals {
  domain_name = data.terraform_remote_state.account.outputs.domain_name
  hosted_zone = data.terraform_remote_state.account.outputs.public_hosted_zone_id

  common_tags = {
    aws_account      = var.aws_account
    aws_region       = data.aws_region.current_region.name
    environment_name = var.environment_name
  }

  rds_db_connect_arn    = "${replace(replace(data.terraform_remote_state.pennsieve_postgres.outputs.rds_proxy_endpoint_arn, ":rds:", ":rds-db:"), ":db-proxy:", ":dbuser:")}/${var.rds_proxy_user}"
}