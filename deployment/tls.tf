resource "tls_private_key" "main" {
  algorithm   = "RSA"
  ecdsa_curve = "P384"
}
