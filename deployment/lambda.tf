resource "aws_lambda_function" "main" {
  function_name    = var.name
  filename         = "lambda.zip"
  source_code_hash = filebase64sha256("lambda.zip")

  role        = aws_iam_role.lambda_exec.arn
  handler     = "main"
  runtime     = "go1.x"
  publish     = true
  timeout     = 30
  memory_size = 4096

  environment {
    variables = {
      BUCKET          = var.bucket
      REDIRECT_URL    = "https://${var.subdomain}.${var.domain}/"
      PUBLIC_KEY_ID   = aws_cloudfront_public_key.main.id
      PRIVATE_KEY_PEM = base64encode(tls_private_key.main.private_key_pem)
    }
  }
}
