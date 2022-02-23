# resource "digitalocean_record" "gateway" {
#   domain = "novacloud.app"
#   type   = "CNAME"
#   name   = var.code
#   ttl    = 300
#   value  = "${aws_apigatewayv2_domain_name.main.domain_name_configuration[0].target_domain_name}."
# }
