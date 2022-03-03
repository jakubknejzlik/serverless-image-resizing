resource "aws_secretsmanager_secret" "private_key" {
  name                    = "${var.name}-private-key"
  recovery_window_in_days = 0
}

locals {
  private_key = {
    publicKeyId   = aws_cloudfront_public_key.main.id
    privateKeyPem = base64encode(tls_private_key.main.private_key_pem)
  }
}

resource "aws_secretsmanager_secret_version" "private_key" {
  secret_id     = aws_secretsmanager_secret.private_key.id
  secret_string = jsonencode(local.private_key)
}
