provider "aws" {}

variable "name" {}

module "default" {
  source = "./deployment"
  name   = var.name
  bucket = var.name
}
