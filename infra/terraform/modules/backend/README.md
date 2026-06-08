# Backend Module

Creates the Google Cloud resources for the public beta backend:

- Required Google Cloud APIs
- Artifact Registry Docker repository
- Secret Manager secrets
- Cloud Run runtime service account
- GitHub Actions deploy service account
- Workload Identity Federation for GitHub Actions
- Optional Cloud Run backend service
- Optional Cloud Run domain mapping

Secret values are intentionally not managed by Terraform.
