variable "project_id" {
  type        = string
  description = "Identifier of the project to reuse"
  default     = "learning-gcloud-444623"
}
variable "region" {
  type        = string
  description = "Global region to use in all services"
  # using us-west so custom domain mapping works
  default = "us-west1"
}
// from common/.variables_resources
variable "resource_db_instance_name" {
  type        = string
  description = "Cloud SQL instance name"
}
variable "resource_service_api_name" {
  type        = string
  description = "service name"
}
variable "resource_service_web_name" {
  type        = string
  description = "service name"
}
// from .secrets/.deploy.env
variable "postgres_dbname" {
  type        = string
  description = "PostgreSQL dbname to init in the Secrets manager"
  sensitive   = true
}
// from .secrets/.deploy.env
variable "postgres_user" {
  type        = string
  description = "PostgreSQL user to init in the Secrets manager"
  sensitive   = true
}
// from .secrets/.deploy.env
variable "postgres_password" {
  type        = string
  description = "PostgreSQL password to init in the Secrets manager"
  sensitive   = true
}