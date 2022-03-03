resource "aws_cloudfront_distribution" "main" {
  origin {
    domain_name = aws_s3_bucket_website_configuration.main.website_endpoint
    origin_id   = "S3Origin"

    custom_origin_config {
      #   origin_access_identity = "origin-access-identity/cloudfront/ABCDEFG1234567"
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1", "TLSv1.1", "TLSv1.2"]
    }
  }

  enabled = true
  #   is_ipv6_enabled     = true
  comment             = var.name
  default_root_object = "index.html"

  #   logging_config {
  #     include_cookies = false
  #     bucket          = "mylogs.s3.amazonaws.com"
  #     prefix          = "myprefix"
  #   }

  aliases = ["${var.subdomain}.${var.domain}"]

  default_cache_behavior {
    allowed_methods    = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods     = ["GET", "HEAD"]
    target_origin_id   = "S3Origin"
    trusted_key_groups = [aws_cloudfront_key_group.main.id]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  #   # Cache behavior with precedence 0
  #   ordered_cache_behavior {
  #     path_pattern     = "/content/immutable/*"
  #     allowed_methods  = ["GET", "HEAD", "OPTIONS"]
  #     cached_methods   = ["GET", "HEAD", "OPTIONS"]
  #     target_origin_id = local.s3_origin_id

  #     forwarded_values {
  #       query_string = false
  #       headers      = ["Origin"]

  #       cookies {
  #         forward = "none"
  #       }
  #     }

  #     min_ttl                = 0
  #     default_ttl            = 86400
  #     max_ttl                = 31536000
  #     compress               = true
  #     viewer_protocol_policy = "redirect-to-https"
  #   }

  #   # Cache behavior with precedence 1
  #   ordered_cache_behavior {
  #     path_pattern     = "/content/*"
  #     allowed_methods  = ["GET", "HEAD", "OPTIONS"]
  #     cached_methods   = ["GET", "HEAD"]
  #     target_origin_id = local.s3_origin_id

  #     forwarded_values {
  #       query_string = false

  #       cookies {
  #         forward = "none"
  #       }
  #     }

  #     min_ttl                = 0
  #     default_ttl            = 3600
  #     max_ttl                = 86400
  #     compress               = true
  #     viewer_protocol_policy = "redirect-to-https"
  #   }

  #   price_class = "PriceClass_200"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    # cloudfront_default_certificate = true
    acm_certificate_arn = data.aws_acm_certificate.thumbcert.arn
    ssl_support_method  = "sni-only"
  }
}

resource "aws_cloudfront_public_key" "main" {
  comment     = "image resizer public key"
  encoded_key = tls_private_key.main.public_key_pem
  name        = "image-resizer-key"
}

resource "aws_cloudfront_key_group" "main" {
  comment = "image resizer key group"
  items   = [aws_cloudfront_public_key.main.id]
  name    = "image-resizer-key-group"
}

data "aws_acm_certificate" "thumbcert" {
  provider    = aws.us-east-1
  domain      = "*.${var.domain}"
  statuses    = ["ISSUED"]
  most_recent = true
}
