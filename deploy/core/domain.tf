# Domain mapping
# Keep around as once off but we do not want to destroy it once it's created to avoid DNS and cert provisioning delay
resource "google_cloud_run_domain_mapping" "web_domain_mapping" {
  name     = "web.fitnesstracker.alexlearningcloud.dev"
  location = var.region
  metadata {
    namespace = var.project_id
  }
  spec {
    route_name = var.resource_service_web_name
  }
}

resource "google_cloud_run_domain_mapping" "api_domain_mapping" {
  name     = "api.fitnesstracker.alexlearningcloud.dev"
  location = var.region
  metadata {
    namespace = var.project_id
  }
  spec {
    route_name = var.resource_service_api_name
  }
}
# Re-add after testing this works well manually from console
# NOTE: managed_zone is hardcoded, need to move to the other script?
resource "google_dns_record_set" "dns_web_cname" {
  name         = "web.fitnesstracker.alexlearningcloud.dev."
  managed_zone = "fitnesstracker-alexlearningcloud-dev"
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["ghs.googlehosted.com."]
  lifecycle {
    prevent_destroy = true
  }
}

resource "google_dns_record_set" "dns_api_cname" {
  name         = "api.fitnesstracker.alexlearningcloud.dev."
  managed_zone = "fitnesstracker-alexlearningcloud-dev"
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["ghs.googlehosted.com."]
  lifecycle {
    prevent_destroy = true
  }
}