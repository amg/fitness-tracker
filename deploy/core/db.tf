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
      # 
      authorized_networks {
        name  = "Home"
        value = var.local_public_ip
      }
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

# seed the db with tables and some default data
# NOTE: have to manually authorize the IP address of the machine running this script in 
#  https://console.cloud.google.com/sql/instances/postgres-instance-ft-staging/connections/networking
# Then run `terraform apply --target null_resource.seed_db` (TODO: need env variables so add this to the scripts ideally)
resource "null_resource" "seed_db" {
  provisioner "local-exec" {
    command = <<EOT
      psql -h ${google_sql_database_instance.postgres.public_ip_address} -U ${var.postgres_user} -d ${var.postgres_dbname} -f ../../db/seed/user_info_seed.sql
      psql -h ${google_sql_database_instance.postgres.public_ip_address} -U ${var.postgres_user} -d ${var.postgres_dbname} -f ../../db/seed/exercises_seed.sql
    EOT
    environment = {
      PGPASSWORD = var.postgres_password
    }
  }

  depends_on = [google_sql_database.database]
}