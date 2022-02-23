# data "aws_acm_certificate" "main" {
#   domain      = "*.novacloud.app"
#   statuses    = ["ISSUED"]
#   most_recent = true
# }

# resource "aws_apigatewayv2_domain_name" "main" {
#   domain_name = "${var.code}.novacloud.app"

#   domain_name_configuration {
#     certificate_arn = data.aws_acm_certificate.main.arn
#     endpoint_type   = "REGIONAL"
#     security_policy = "TLS_1_2"
#   }
# }

# resource "aws_apigatewayv2_api_mapping" "default" {
#   api_id      = aws_apigatewayv2_api.main.id
#   domain_name = aws_apigatewayv2_domain_name.main.id
#   stage       = "$default"
# }
