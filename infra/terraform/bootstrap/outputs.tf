output "state_bucket_name" {
  description = "GCS bucket name used for Terraform remote state."
  value       = google_storage_bucket.terraform_state.name
}

output "state_bucket_url" {
  description = "GCS URL for the Terraform remote state bucket."
  value       = google_storage_bucket.terraform_state.url
}
