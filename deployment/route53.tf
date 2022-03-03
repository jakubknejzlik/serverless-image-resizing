data "aws_route53_zone" "main" {
  name = "${var.domain}."
}
resource "aws_route53_record" "thumbnails" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "${var.subdomain}.${var.domain}"
  type    = "CNAME"
  ttl     = "300"
  records = [aws_cloudfront_distribution.main.domain_name]
}
