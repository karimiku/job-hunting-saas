variable "project_id" {
  description = "Google Cloud project ID."
  type        = string
}

variable "region" {
  description = "Google Cloud region for regional backend resources."
  type        = string
}

variable "service_apis" {
  description = "Public Google Cloud APIs required by the backend infrastructure."
  type        = set(string)
  default = [
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "firebase.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "identitytoolkit.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com",
    "serviceusage.googleapis.com",
    "sts.googleapis.com",
  ]
}

variable "frontend_domain" {
  description = "Frontend domain hosted by Vercel."
  type        = string
}

variable "backend_domain" {
  description = "Backend API domain."
  type        = string
}

variable "artifact_repository_id" {
  description = "Artifact Registry Docker repository ID."
  type        = string
}

variable "backend_service_name" {
  description = "Cloud Run backend service name."
  type        = string
}

variable "runtime_service_account_id" {
  description = "Service account ID for Cloud Run runtime identity."
  type        = string
  default     = "entre-backend-runtime"
}

variable "github_deploy_service_account_id" {
  description = "Service account ID impersonated by GitHub Actions deploy workflows."
  type        = string
  default     = "github-deploy"
}

variable "github_repository" {
  description = "GitHub repository in owner/name form."
  type        = string
}

variable "github_repository_id" {
  description = "Immutable GitHub repository ID used in Workload Identity Federation conditions."
  type        = string
}

variable "secret_ids" {
  description = "Secret Manager secret IDs used by the backend."
  type        = set(string)
  default = [
    "database-url",
  ]
}

variable "database_url_secret_id" {
  description = "Secret Manager secret ID for DATABASE_URL."
  type        = string
  default     = "database-url"
}

variable "database_url_secret_version" {
  description = "Pinned Secret Manager version number for DATABASE_URL."
  type        = string
  default     = "1"
}

variable "enable_backend_service" {
  description = "Whether to create the Cloud Run backend service."
  type        = bool
  default     = false
}

variable "enable_domain_mapping" {
  description = "Whether to create a Cloud Run domain mapping for backend_domain."
  type        = bool
  default     = false
}

variable "enable_public_invoker" {
  description = "Whether to allow unauthenticated HTTP invocation of the Cloud Run service."
  type        = bool
  default     = true
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

variable "backend_container_image" {
  description = "Initial backend container image. Later image changes are deployed by CI/CD."
  type        = string
}

variable "container_port" {
  description = "Backend container port."
  type        = number
  default     = 8080
}

variable "cloud_run_min_instance_count" {
  description = "Minimum Cloud Run instance count."
  type        = number
  default     = 0
}

variable "cloud_run_max_instance_count" {
  description = "Maximum Cloud Run instance count."
  type        = number
  default     = 1
}

variable "cloud_run_cpu" {
  description = "Cloud Run CPU limit."
  type        = string
  default     = "1"
}

variable "cloud_run_memory" {
  description = "Cloud Run memory limit."
  type        = string
  default     = "512Mi"
}

variable "cloud_run_deletion_protection" {
  description = "Whether to protect the Cloud Run service from deletion."
  type        = bool
  default     = true
}

variable "cookie_domain" {
  description = "Session cookie domain."
  type        = string
  default     = ".entre.kamiriku.com"
}

variable "cookie_secure" {
  description = "Whether session cookies must be marked Secure."
  type        = bool
  default     = true
}

variable "cookie_same_site" {
  description = "Session cookie SameSite policy."
  type        = string
  default     = "strict"

  validation {
    condition     = contains(["strict", "lax", "none"], var.cookie_same_site)
    error_message = "cookie_same_site must be one of strict, lax, or none."
  }
}

variable "cors_allowed_origins" {
  description = "Allowed CORS origins for cookie-bearing requests."
  type        = list(string)
}
