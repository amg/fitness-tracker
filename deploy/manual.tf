// ------- DNS A records required for load balancer
// In prod env you would reserve the external IP Address so it doesn't change
//  and these steps will be simpler

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

//------ LB rotutes

# For now done manually, figure out how to attach this to the LB
# resource "google_compute_url_map" "main" {
#   name            = "url-map-main"
#   default_service = google_compute_region_network_endpoint_group.web_neg.name

#   host_rule {
#     hosts        = ["web.fitnesstracker.alexlearningcloud.dev"]
#     path_matcher = "web"
#   }

#   host_rule {
#     hosts        = ["api.fitnesstracker.alexlearningcloud.dev"]
#     path_matcher = "api"
#   }

#   path_matcher {
#     name            = "web"
#     default_service = google_compute_region_network_endpoint_group.web_neg.name
#   }

#   path_matcher {
#     name            = "api"
#     default_service = google_compute_region_network_endpoint_group.api_neg.name
#   }
# }
