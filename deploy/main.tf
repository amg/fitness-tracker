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
resource "google_cloud_run_v2_service" "web" {
  name     = "web"
  location = var.region
  template {
    # web is using default cloud run account
    scaling {
      max_instance_count = 2
    }
    containers {
      image = "gcr.io/${var.project_id}/staging-web"
      resources {
        limits = {
          cpu    = "1"
          memory = "1024Mi"
        }
      }
    }
  }
  deletion_protection = false
}

# Deploy api to Cloud Run
resource "google_cloud_run_v2_service" "api" {
  name     = "api"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.api.email
    scaling {
      max_instance_count = 2
    }
    containers {
      image = "gcr.io/${var.project_id}/staging-api"

      # Sets a environment variable for instance connection name
      env {
        name  = "DB_INSTANCE_CONNECTION_NAME"
        value = "${var.project_id}:${var.region}:${var.resource_db_instance_name}"
      }
    }
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = ["${var.project_id}:${var.region}:${var.resource_db_instance_name}"]
      }
    }
  }
  deletion_protection = false
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

# SQL client role
resource "google_project_iam_binding" "cloud_sql_access" {
  project = var.project_id
  role    = "roles/cloudsql.client"

  members = [
    "serviceAccount:${google_service_account.api.email}",
  ]
}

# Enable public access on Cloud Run service
resource "google_cloud_run_service_iam_policy" "web_noauth" {
  location    = google_cloud_run_v2_service.web.location
  project     = google_cloud_run_v2_service.web.project
  service     = google_cloud_run_v2_service.web.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

# Enable public access on Cloud Run service
resource "google_cloud_run_service_iam_policy" "api_noauth" {
  location    = google_cloud_run_v2_service.api.location
  project     = google_cloud_run_v2_service.api.project
  service     = google_cloud_run_v2_service.api.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

# Access to secret manager by Cloud Run service
resource "google_service_account" "api" {
  account_id   = "cloud-run-service-account"
  display_name = "Service account for Cloud Run"
}

resource "google_secret_manager_secret_iam_member" "iam_secret_google_oauth_client_secret" {
  secret_id = google_secret_manager_secret.secret_google_oauth_client_secret.id
  role      = "roles/secretmanager.secretAccessor"
  # Grant the new deployed service account access to this secret.
  member     = "serviceAccount:${google_service_account.api.email}"
  depends_on = [google_secret_manager_secret.secret_google_oauth_client_secret]
}
resource "google_secret_manager_secret_iam_member" "iam_secret_jwt_private_pem" {
  secret_id = google_secret_manager_secret.secret_jwt_private_pem.id
  role      = "roles/secretmanager.secretAccessor"
  # Grant the new deployed service account access to this secret.
  member     = "serviceAccount:${google_service_account.api.email}"
  depends_on = [google_secret_manager_secret.secret_jwt_private_pem]
}
resource "google_secret_manager_secret_iam_member" "iam_secret_jwt_public_pem" {
  secret_id = google_secret_manager_secret.secret_jwt_public_pem.id
  role      = "roles/secretmanager.secretAccessor"
  # Grant the new deployed service account access to this secret.
  member     = "serviceAccount:${google_service_account.api.email}"
  depends_on = [google_secret_manager_secret.secret_jwt_public_pem]
}

// -------- Env ---------
# GOOGLE_OAUTH_CLIENT_SECRET
resource "google_secret_manager_secret" "secret_google_oauth_client_secret" {
  secret_id = "GOOGLE_OAUTH_CLIENT_SECRET"
  project   = var.project_id

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret_google_oauth_client_secret_data" {
  secret      = google_secret_manager_secret.secret_google_oauth_client_secret.id
  secret_data = var.seed_secret_google_oauth_secret
}

# JWT_KEY_PRIVATE
resource "google_secret_manager_secret" "secret_jwt_private_pem" {
  secret_id = "JWT_KEY_PRIVATE"
  project   = var.project_id

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret_jwt_private_pem_data" {
  secret      = google_secret_manager_secret.secret_jwt_private_pem.id
  secret_data = file("${path.module}/../.secrets/jwtRSA256-private.pem")
}

# JWT_KEY_PUBLIC
resource "google_secret_manager_secret" "secret_jwt_public_pem" {
  secret_id = "JWT_KEY_PUBLIC"
  project   = var.project_id

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret_jwt_public_pem_data" {
  secret      = google_secret_manager_secret.secret_jwt_public_pem.id
  secret_data = file("${path.module}/../.secrets/jwtRSA256-public.pem")
}