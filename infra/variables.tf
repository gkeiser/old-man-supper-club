variable "aws_region" {
  default = "us-east-2"
}

variable "dynamo_table_name" {
  default = "SupperClub"
}

variable "google_client_id" {
  type      = string
  sensitive = true
}

variable "google_client_secret" {
  type      = string
  sensitive = true
}

variable "google_redirect_url" {
  type = string
}

variable "allowed_emails" {
  type = string
}

variable "session_secret" {
  type      = string
  sensitive = true
}
