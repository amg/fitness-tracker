terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.13.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Deploy image to Cloud Run
resource "google_cloud_run_service" "web" {
  name     = "web"
  location = var.region
  template {
    spec {
      containers {
        image = "gcr.io/${var.project_id}/web"
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Deploy api to Cloud Run
resource "google_cloud_run_service" "api" {
  name     = "api"
  location = var.region
  template {
    spec {
      containers {
        image = "gcr.io/${var.project_id}/api"
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Create public access
data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

# Enable public access on Cloud Run service
resource "google_cloud_run_service_iam_policy" "web_noauth" {
  location    = google_cloud_run_service.web.location
  project     = google_cloud_run_service.web.project
  service     = google_cloud_run_service.web.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

# Enable public access on Cloud Run service
resource "google_cloud_run_service_iam_policy" "api_noauth" {
  location    = google_cloud_run_service.api.location
  project     = google_cloud_run_service.api.project
  service     = google_cloud_run_service.api.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

# Domain mapping
resource "google_cloud_run_domain_mapping" "web_domain_mapping" {
  name     = "web.fitnesstracker.alexlearningcloud.dev"
  location = google_cloud_run_service.web.location
  metadata {
    namespace = var.project_id
  }
  spec {
    route_name = google_cloud_run_service.web.name
  }
}

resource "google_cloud_run_domain_mapping" "api_domain_mapping" {
  name     = "api.fitnesstracker.alexlearningcloud.dev"
  location = google_cloud_run_service.api.location
  metadata {
    namespace = var.project_id
  }
  spec {
    route_name = google_cloud_run_service.api.name
  }
}

resource "google_dns_record_set" "dns_web_cname" {
  name         = "web.${google_dns_managed_zone.zone.dns_name}"
  managed_zone = google_dns_managed_zone.zone.name
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["web.fitnesstracker.alexlearningcloud.dev."]
}

resource "google_dns_record_set" "dns_api_cname" {
  name         = "api.${google_dns_managed_zone.zone.dns_name}"
  managed_zone = google_dns_managed_zone.zone.name
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["api.fitnesstracker.alexlearningcloud.dev."]
}

resource "google_dns_managed_zone" "zone" {
  name     = "dns-zone"
  dns_name = "fitnesstracker.alexlearningcloud.dev."
}
