// Create log group for service lambda.
resource "aws_cloudwatch_log_group" "drs_service_api_lambda_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.drs_service_api_lambda.function_name}"
  retention_in_days = 30
  tags = local.common_tags
}

# TIMESERIES SERVICE API GATEWAY
resource "aws_cloudwatch_log_group" "drs_service_api_gateway_log_group" {
  name =  "${var.environment_name}/${var.service_name}/drs-api-gateway"

  retention_in_days = 30
}
