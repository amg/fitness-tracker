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