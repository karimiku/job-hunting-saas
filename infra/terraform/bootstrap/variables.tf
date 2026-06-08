variable "project_id" {
  description = "Google Cloud project ID that owns the Terraform state bucket."
  type        = string
}

variable "state_bucket_name" {
  description = "Globally unique GCS bucket name for Terraform state."
  type        = string
}

variable "state_bucket_location" {
  description = "GCS location for the Terraform state bucket."
  type        = string
  default     = "ASIA"
}

variable "service_apis" {
  description = "Google Cloud APIs required by the bootstrap foundation."
  type        = set(string)
  default = [
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "serviceusage.googleapis.com",
    "storage.googleapis.com",
    "sts.googleapis.com",
  ]
}

variable "github_repository" {
  description = "GitHub repository in owner/name form."
  type        = string
}

variable "github_repository_id" {
  description = "Immutable GitHub repository ID used in Workload Identity Federation conditions."
  type        = string
}

variable "terraform_service_account_id" {
  description = "Service account ID impersonated by GitHub Actions Terraform workflows."
  type        = string
  default     = "github-terraform"
}

variable "workload_identity_pool_id" {
  description = "Workload Identity Pool ID for GitHub Actions."
  type        = string
  default     = "github-actions"
}

variable "workload_identity_provider_id" {
  description = "Workload Identity Pool Provider ID for GitHub Actions."
  type        = string
  default     = "github"
}

variable "terraform_state_bucket_role" {
  description = "Bucket-level IAM role granted to the Terraform service account for remote state access."
  type        = string
  default     = "roles/storage.objectAdmin"
}

variable "terraform_project_roles" {
  description = "Project-level IAM roles granted to the Terraform service account for managing production infrastructure."
  type        = set(string)
  default = [
    "roles/artifactregistry.admin",
    "roles/iam.serviceAccountAdmin",
    "roles/iam.serviceAccountUser",
    "roles/resourcemanager.projectIamAdmin",
    "roles/run.admin",
    "roles/secretmanager.admin",
    "roles/serviceusage.serviceUsageAdmin",
  ]
}
