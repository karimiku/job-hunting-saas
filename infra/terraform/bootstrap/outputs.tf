output "state_bucket_name" {
  description = "GCS bucket name used for Terraform remote state."
  value       = google_storage_bucket.terraform_state.name
}

output "state_bucket_url" {
  description = "GCS URL for the Terraform remote state bucket."
  value       = google_storage_bucket.terraform_state.url
}

output "terraform_service_account_email" {
  description = "Service account email impersonated by GitHub Actions Terraform workflows."
  value       = google_service_account.terraform.email
}

output "workload_identity_pool_name" {
  description = "Full Workload Identity Pool resource name."
  value       = google_iam_workload_identity_pool.github.name
}

output "workload_identity_provider_name" {
  description = "Full Workload Identity Provider resource name for google-github-actions/auth."
  value       = google_iam_workload_identity_pool_provider.github.name
}
