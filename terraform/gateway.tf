resource "aws_apigatewayv2_api" "drs-service-api" {
  name          = "Serverless DRS API"
  protocol_type = "HTTP"
  description = "API for the lambda-based DRS API"
  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["*"]
    allow_headers = ["*"]
    expose_headers = ["*"]
    max_age = 300
  }
  body          = templatefile("${path.module}/drs_service.yml", {
    authorize_lambda_invoke_uri = data.terraform_remote_state.api_gateway.outputs.authorizer_lambda_invoke_uri,
    gateway_authorizer_role = data.terraform_remote_state.api_gateway.outputs.authorizer_invocation_role,
    drs_service_api_lambda_arn = aws_lambda_function.drs_service_api_lambda.arn
  })
}

resource "aws_apigatewayv2_api_mapping" "drs-service-api-map" {
  api_id          = aws_apigatewayv2_api.drs-service-api.id
  domain_name     = var.api_domain_name
  stage           = aws_apigatewayv2_stage.drs-service-gateway-stage.id
}

resource "aws_apigatewayv2_stage" "drs-service-gateway-stage" {
  api_id = aws_apigatewayv2_api.drs-service-api.id

  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.drs_service_api_gateway_log_group.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
    }
    )
  }
}

resource "aws_apigatewayv2_integration" "int" {
  api_id           = aws_apigatewayv2_api.drs-service-api.id
  integration_type = "AWS_PROXY"
  connection_type = "INTERNET"
  integration_method = "POST"
  integration_uri = aws_lambda_function.drs_service_api_lambda.invoke_arn
}

resource "aws_lambda_permission" "drs-service-lambda-permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.drs_service_api_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.drs-service-api.execution_arn}/*/*"
}
