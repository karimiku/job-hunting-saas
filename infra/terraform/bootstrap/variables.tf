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
