provider "aws" {}

variable "name" {}
variable "domain" {}

module "default" {
  source = "./deployment"
  name   = var.name
  bucket = var.name
  domain = var.domain
}
