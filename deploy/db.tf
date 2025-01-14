# DB secrets
# - postgres_dbname
# - postgres_user
# - postgres_password

# -------- postgres_dbname
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

# Update service account for dbuser secret
resource "google_secret_manager_secret_iam_member" "postgres_dbname_compute_secretaccess" {
  secret_id = google_secret_manager_secret.postgres_dbname.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.api.email}"
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

# Update service account for user secret
resource "google_secret_manager_secret_iam_member" "postgres_user_compute_secretaccess" {
  secret_id = google_secret_manager_secret.postgres_user.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.api.email}"
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

# Update service account for password secret
resource "google_secret_manager_secret_iam_member" "postgres_password_compute_secretaccess" {
  secret_id = google_secret_manager_secret.postgres_password.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.api.email}"
}

# ------- db itself
# Commented out for now and removed from state to speed up deployment
# resource "google_sql_database_instance" "postgres" {
#   name                = "postgres-instance-ft-staging"
#   database_version    = "POSTGRES_15"
#   deletion_protection = false

#   settings {
#     tier = "db-f1-micro"
#     user_labels = {
#       "environment" = "staging"
#     }
#     ip_configuration {
#       # Enable public IP
#       ipv4_enabled = true
#       # Require ssl but we are using auth proxy so this is not really used
#       require_ssl = true
#     }
#   }
# }
# resource "google_sql_user" "user" {
#   name     = var.postgres_user
#   instance = google_sql_database_instance.postgres.name
#   password = var.postgres_password
# }
# resource "google_sql_database" "database" {
#   name     = var.postgres_dbname
#   instance = google_sql_database_instance.postgres.name
# }
