variable "project_id" {
  description = "Production Google Cloud project ID."
  type        = string
}

variable "region" {
  description = "Primary Google Cloud region for production backend resources."
  type        = string
}

variable "frontend_domain" {
  description = "Production frontend domain hosted by Vercel."
  type        = string
}

variable "backend_domain" {
  description = "Production backend API domain."
  type        = string
}

variable "github_repository" {
  description = "GitHub repository in owner/name form."
  type        = string
}

variable "github_repository_id" {
  description = "Immutable GitHub repository ID used for Workload Identity Federation conditions."
  type        = string
}

variable "enable_backend_service" {
  description = "Whether to create the Cloud Run backend service. Enable after image and secret versions exist."
  type        = bool
  default     = false
}

variable "enable_domain_mapping" {
  description = "Whether to create Cloud Run domain mapping for backend_domain."
  type        = bool
  default     = false
}

variable "backend_container_image" {
  description = "Initial backend container image used when creating Cloud Run. Later image changes are deployed by CI/CD."
  type        = string
}

variable "database_url_secret_version" {
  description = "Pinned Secret Manager version number for DATABASE_URL."
  type        = string
  default     = "1"
}

variable "enable_github_deploy_wif" {
  description = "Whether to allow GitHub Actions to impersonate the deploy service account through the bootstrap WIF pool."
  type        = bool
  default     = true
}

variable "github_actions_workload_identity_pool_id" {
  description = "Existing GitHub Actions Workload Identity Pool ID created by bootstrap."
  type        = string
  default     = "github-actions"
}

variable "github_actions_workload_identity_pool_project_number" {
  description = "Project number that owns the GitHub Actions Workload Identity Pool. Empty means the current project."
  type        = string
  default     = ""
}

variable "cors_allowed_origins" {
  description = "Allowed CORS origins for cookie-bearing API requests."
  type        = list(string)
}
