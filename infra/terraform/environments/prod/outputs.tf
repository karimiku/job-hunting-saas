output "artifact_repository_url" {
  description = "Docker repository URL for backend images."
  value       = module.backend.artifact_repository_url
}

output "backend_service_url" {
  description = "Cloud Run backend service URL, if the service is enabled."
  value       = module.backend.backend_service_url
}

output "github_deploy_service_account_email" {
  description = "Service account email impersonated by GitHub Actions deploy workflows."
  value       = module.backend.github_deploy_service_account_email
}

output "runtime_service_account_email" {
  description = "Cloud Run runtime service account email."
  value       = module.backend.runtime_service_account_email
}

output "workload_identity_provider_name" {
  description = "Full Workload Identity provider resource name for google-github-actions/auth."
  value       = module.backend.workload_identity_provider_name
}
