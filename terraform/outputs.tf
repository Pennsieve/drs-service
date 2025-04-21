output "service_lambda_arn" {
  value = aws_lambda_function.drs_service_api_lambda.arn
}

output "service_lambda_invoke_arn" {
  value = aws_lambda_function.drs_service_api_lambda.invoke_arn
}

output "service_lambda_function_name" {
  value = aws_lambda_function.drs_service_api_lambda.function_name
}
