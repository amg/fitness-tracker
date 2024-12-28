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

resource "google_cloud_run_domain_mapping" "web-domain-mapping" {
  location = var.region
  name     = "fitnesstracker.alexlearningcloud.dev/web"

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_service.web.name
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

resource "google_cloud_run_domain_mapping" "api-domain-mapping" {
  location = var.region
  name     = "fitnesstracker.alexlearningcloud.dev/api"

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_service.api.name
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
resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.web.location
  project     = google_cloud_run_service.web.project
  service     = google_cloud_run_service.web.name
  policy_data = data.google_iam_policy.noauth.policy_data
}