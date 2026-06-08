locals {
  labels = {
    app         = "entre"
    environment = "prod"
    managed_by  = "terraform"
  }

  cors_allowed_origins = join(",", var.cors_allowed_origins)
  github_actions_workload_identity_pool_project_number = (
    var.github_actions_workload_identity_pool_project_number != ""
    ? var.github_actions_workload_identity_pool_project_number
    : data.google_project.main.number
  )
}

data "google_project" "main" {
  project_id = var.project_id
}

resource "google_project_service" "service" {
  for_each = var.service_apis

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}

resource "google_artifact_registry_repository" "backend" {
  project       = var.project_id
  location      = var.region
  repository_id = var.artifact_repository_id
  description   = "Docker images for the Entre backend"
  format        = "DOCKER"

  labels = local.labels

  depends_on = [google_project_service.service]
}

resource "google_service_account" "runtime" {
  project      = var.project_id
  account_id   = var.runtime_service_account_id
  display_name = "Entre backend Cloud Run runtime"

  depends_on = [google_project_service.service]
}

resource "google_service_account" "github_deploy" {
  project      = var.project_id
  account_id   = var.github_deploy_service_account_id
  display_name = "GitHub Actions deployer for Entre backend"

  depends_on = [google_project_service.service]
}

resource "google_secret_manager_secret" "secret" {
  for_each = var.secret_ids

  project   = var.project_id
  secret_id = each.value

  replication {
    auto {}
  }

  labels = local.labels

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.service]
}

resource "google_secret_manager_secret_iam_member" "runtime_secret_accessor" {
  for_each = google_secret_manager_secret.secret

  secret_id = each.value.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.runtime.email}"
}

resource "google_project_iam_member" "runtime_firebase_auth_admin" {
  project = var.project_id
  role    = "roles/firebaseauth.admin"
  member  = "serviceAccount:${google_service_account.runtime.email}"

  depends_on = [google_project_service.service]
}

resource "google_service_account_iam_member" "github_deploy_workload_identity_user" {
  count = var.enable_github_deploy_wif ? 1 : 0

  service_account_id = google_service_account.github_deploy.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/projects/${local.github_actions_workload_identity_pool_project_number}/locations/global/workloadIdentityPools/${var.github_actions_workload_identity_pool_id}/attribute.repository_id/${var.github_repository_id}"
}

resource "google_artifact_registry_repository_iam_member" "github_deploy_artifact_writer" {
  project    = var.project_id
  location   = google_artifact_registry_repository.backend.location
  repository = google_artifact_registry_repository.backend.repository_id
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${google_service_account.github_deploy.email}"
}

resource "google_project_iam_member" "github_deploy_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.github_deploy.email}"
}

resource "google_service_account_iam_member" "github_deploy_runtime_act_as" {
  service_account_id = google_service_account.runtime.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_deploy.email}"
}

resource "google_cloud_run_v2_service" "backend" {
  count = var.enable_backend_service ? 1 : 0

  project             = var.project_id
  name                = var.backend_service_name
  location            = var.region
  ingress             = "INGRESS_TRAFFIC_ALL"
  deletion_protection = var.cloud_run_deletion_protection

  labels = local.labels

  template {
    service_account = google_service_account.runtime.email

    scaling {
      min_instance_count = var.cloud_run_min_instance_count
      max_instance_count = var.cloud_run_max_instance_count
    }

    containers {
      image = var.backend_container_image

      ports {
        container_port = var.container_port
      }

      resources {
        limits = {
          cpu    = var.cloud_run_cpu
          memory = var.cloud_run_memory
        }
      }

      env {
        name  = "FIREBASE_PROJECT_ID"
        value = var.project_id
      }

      env {
        name  = "COOKIE_DOMAIN"
        value = var.cookie_domain
      }

      env {
        name  = "COOKIE_SECURE"
        value = tostring(var.cookie_secure)
      }

      env {
        name  = "COOKIE_SAME_SITE"
        value = var.cookie_same_site
      }

      env {
        name  = "CORS_ALLOWED_ORIGINS"
        value = local.cors_allowed_origins
      }

      env {
        name = "DATABASE_URL"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.secret[var.database_url_secret_id].secret_id
            version = var.database_url_secret_version
          }
        }
      }
    }
  }

  depends_on = [
    google_artifact_registry_repository.backend,
    google_project_iam_member.runtime_firebase_auth_admin,
    google_secret_manager_secret_iam_member.runtime_secret_accessor,
  ]

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
    ]
  }
}

resource "google_cloud_run_v2_service_iam_binding" "public_invoker" {
  count = var.enable_backend_service && var.enable_public_invoker ? 1 : 0

  project  = var.project_id
  location = google_cloud_run_v2_service.backend[0].location
  name     = google_cloud_run_v2_service.backend[0].name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}

resource "google_cloud_run_domain_mapping" "backend" {
  count = var.enable_backend_service && var.enable_domain_mapping ? 1 : 0

  name     = var.backend_domain
  location = google_cloud_run_v2_service.backend[0].location

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.backend[0].name
  }
}
