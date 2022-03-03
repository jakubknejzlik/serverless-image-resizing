terraform {
  required_providers {
    aws = {
      source                = "hashicorp/aws"
      version               = "4.1.0"
      configuration_aliases = [aws.us-east-1]
    }
  }
}

# provider "aws" {}
provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

# provider "mysql" {}
# provider "random" {}
