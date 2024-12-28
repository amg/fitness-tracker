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
  # default     = "australia-southeast2"
}