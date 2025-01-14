# Domain mapping (for now it's manual)

# IMPORTANT: No longer required with load balancer on
# resource "google_cloud_run_domain_mapping" "web_domain_mapping" {
#   name     = "web.fitnesstracker.alexlearningcloud.dev"
#   location = var.region
#   metadata {
#     namespace = var.project_id
#   }
#   spec {
#     route_name = var.resource_service_web_name
#   }
#   lifecycle {
#     prevent_destroy = true
#   }
# }

# resource "google_cloud_run_domain_mapping" "api_domain_mapping" {
#   name     = "api.fitnesstracker.alexlearningcloud.dev"
#   location = var.region
#   metadata {
#     namespace = var.project_id
#   }
#   spec {
#     route_name = var.resource_service_api_name
#   }
#   lifecycle {
#     prevent_destroy = true
#   }
# }

# DNS A records required for load balancer
# IMPORTANT: LB IP address required to map those, or do manually
# resource "google_dns_record_set" "dns_web_cname" {
#   name         = "web.fitnesstracker.alexlearningcloud.dev."
#   managed_zone = "fitnesstracker-alexlearningcloud-dev"
#   type         = "A"
#   ttl          = 300
#   rrdatas      = ["<LB-IP-ADDRESS>"]
#   lifecycle {
#     prevent_destroy = true
#   }
# }

# resource "google_dns_record_set" "dns_api_cname" {
#   name         = "api.fitnesstracker.alexlearningcloud.dev."
#   managed_zone = "fitnesstracker-alexlearningcloud-dev"
#   type         = "CNAME"
#   ttl          = 300
#   rrdatas      = ["<LB-IP-ADDRESS>"]
#   lifecycle {
#     prevent_destroy = true
#   }
# }