terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.69.0"
    }
  }
}

variable "target-region" {
  type = string
  default = "us-east-1"
}

variable "target-user-profile" {
  type = string
  default = "default"
}

provider "aws" {
  region  = var.target-region
  profile = var.target-user-profile
}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

#############################
# Variables for customization
#############################

variable "billing-mode" {
  type    = string
  default = "PAY_PER_REQUEST" # Possible values are "PAY_PER_REQUEST" and "PROVISIONED"
                              # Pay per request for PoC to save some money
}

variable "read-capacity" {
  type    = number
  default = 5
}

variable "write-capacity" {
  type    = number
  default = 5
}

#################
# Dynamo DB setup
#################

resource "aws_dynamodb_table" "requests-dynamodb-table" {
  name           = "RequestRecords"
  billing_mode   = var.billing-mode
  read_capacity  = var.read-capacity
  write_capacity = var.write-capacity
  hash_key       = "Id"
  range_key      = "Timestamp"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "S"
  }

  ttl {
    attribute_name = "TimeToExist"
    enabled        = false
  }

  tags = {
    Name        = "dynamodb-table-request-records"
    Environment = "poc"
  }
}
