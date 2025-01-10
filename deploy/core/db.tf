# ------- db itself
resource "google_sql_database_instance" "postgres" {
  name             = var.resource_db_instance_name
  database_version = "POSTGRES_15"
  // change to true in prod
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    user_labels = {
      "environment" = "staging"
    }
    ip_configuration {
      # Enable public IP
      ipv4_enabled = true
    }
  }
}
resource "google_sql_user" "user" {
  name     = var.postgres_user
  instance = google_sql_database_instance.postgres.name
  password = var.postgres_password
}
resource "google_sql_database" "database" {
  name     = var.postgres_dbname
  instance = google_sql_database_instance.postgres.name
}
