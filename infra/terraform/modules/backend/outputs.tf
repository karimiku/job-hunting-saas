output "artifact_repository_url" {
  description = "Docker repository URL for backend images."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.backend.repository_id}"
}

output "runtime_service_account_email" {
  description = "Cloud Run runtime service account email."
  value       = google_service_account.runtime.email
}

output "github_deploy_service_account_email" {
  description = "GitHub Actions deploy service account email."
  value       = google_service_account.github_deploy.email
}

output "backend_service_url" {
  description = "Cloud Run backend service URL, if enabled."
  value       = var.enable_backend_service ? google_cloud_run_v2_service.backend[0].uri : null
}

output "backend_domain_mapping_records" {
  description = "Cloud Run domain mapping resource records, if domain mapping is enabled."
  value       = var.enable_backend_service && var.enable_domain_mapping ? google_cloud_run_domain_mapping.backend[0].status[0].resource_records : null
}
