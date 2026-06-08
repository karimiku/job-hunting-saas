locals {
  labels = {
    app         = "entre"
    managed_by  = "terraform"
    environment = "bootstrap"
  }
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

resource "google_storage_bucket" "terraform_state" {
  name     = var.state_bucket_name
  project  = var.project_id
  location = var.state_bucket_location

  force_destroy               = false
  public_access_prevention    = "enforced"
  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }

  labels = local.labels

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.service]
}

resource "google_service_account" "terraform" {
  project      = var.project_id
  account_id   = var.terraform_service_account_id
  display_name = "GitHub Actions Terraform runner"

  depends_on = [google_project_service.service]
}

resource "google_storage_bucket_iam_member" "terraform_state_object_admin" {
  bucket = google_storage_bucket.terraform_state.name
  role   = var.terraform_state_bucket_role
  member = "serviceAccount:${google_service_account.terraform.email}"
}

resource "google_project_iam_member" "terraform_project_role" {
  for_each = var.terraform_project_roles

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.terraform.email}"
}

resource "google_iam_workload_identity_pool" "github" {
  project                   = var.project_id
  workload_identity_pool_id = var.workload_identity_pool_id
  display_name              = "GitHub Actions"
  description               = "Federates GitHub Actions OIDC tokens for ${var.github_repository}"

  depends_on = [google_project_service.service]
}

resource "google_iam_workload_identity_pool_provider" "github" {
  project                            = var.project_id
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = var.workload_identity_provider_id
  display_name                       = "GitHub"
  description                        = "Trusts GitHub Actions tokens for ${var.github_repository}"

  attribute_mapping = {
    "google.subject"          = "assertion.sub"
    "attribute.repository"    = "assertion.repository"
    "attribute.repository_id" = "assertion.repository_id"
    "attribute.ref"           = "assertion.ref"
    "attribute.workflow_ref"  = "assertion.workflow_ref"
  }

  attribute_condition = "assertion.repository_id == '${var.github_repository_id}'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account_iam_member" "terraform_workload_identity_user" {
  service_account_id = google_service_account.terraform.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/projects/${data.google_project.main.number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.github.workload_identity_pool_id}/attribute.repository_id/${var.github_repository_id}"
}
