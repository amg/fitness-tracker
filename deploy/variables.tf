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
variable "seed_secret_google_oauth_secret" {
  type        = string
  description = "OAuth secret to init in the Secrets manager"
}