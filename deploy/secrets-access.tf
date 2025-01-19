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