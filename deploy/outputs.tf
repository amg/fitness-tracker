# Return service URL
output "url" {
  value = google_cloud_run_v2_service.web.urls
}