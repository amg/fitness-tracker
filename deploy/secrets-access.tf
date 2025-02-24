# duplicate from core/db.tf
resource "google_secret_manager_secret" "postgres_dbname" {
  secret_id = "POSTGRES_DBNAME"
  replication {
    auto {}
  }
}

# Attaches secret data for dbuser secret
resource "google_secret_manager_secret_version" "postgres_dbname_data" {
  secret      = google_secret_manager_secret.postgres_dbname.id
  secret_data = var.postgres_dbname # Stores secret as a plain txt in state
}

# -------- postgres_user
resource "google_secret_manager_secret" "postgres_user" {
  secret_id = "POSTGRES_USER"
  replication {
    auto {}
  }
}

# Attaches secret data for user secret
resource "google_secret_manager_secret_version" "postgres_user_data" {
  secret      = google_secret_manager_secret.postgres_user.id
  secret_data = var.postgres_user # Stores secret as a plain txt in state
}

# -------- postgres_password
resource "google_secret_manager_secret" "postgres_password" {
  secret_id = "POSTGRES_PASSWORD"
  replication {
    auto {}
  }
}

# Attaches secret data for password secret
resource "google_secret_manager_secret_version" "postgres_password_data" {
  secret      = google_secret_manager_secret.postgres_password.id
  secret_data = var.postgres_password # Stores secret as a plain txt in state
}

# secrets access permissions for api service
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

# Update service account for dbuser secret
resource "google_secret_manager_secret_iam_member" "postgres_dbname_compute_secretaccess" {
  secret_id = google_secret_manager_secret.postgres_dbname.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.api.email}"
}

# Update service account for user secret
resource "google_secret_manager_secret_iam_member" "postgres_user_compute_secretaccess" {
  secret_id = google_secret_manager_secret.postgres_user.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.api.email}"
}

# Update service account for password secret
resource "google_secret_manager_secret_iam_member" "postgres_password_compute_secretaccess" {
  secret_id = google_secret_manager_secret.postgres_password.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.api.email}"
}