# Terraform

This directory manages the public beta Google Cloud infrastructure.

Layout follows Google Cloud Terraform root module guidance:

```text
bootstrap/              # one-time remote state bucket
environments/prod/      # production root module
modules/backend/        # reusable backend service module
```

## Bootstrap

Create the GCS bucket used by the production remote state:

```bash
cd infra/terraform/bootstrap
terraform init
terraform apply
```

After the bucket exists, initialize the production environment:

```bash
cd ../environments/prod
terraform init
terraform plan
```

## Secrets

Terraform creates Secret Manager secret resources and IAM bindings, but it does not store secret values.

Add secret versions manually or from a protected CI job:

```bash
gcloud secrets versions add database-url --data-file=/path/to/database-url.txt --project=job-hunting-saas
```

Do not commit `DATABASE_URL`, API keys, service account JSON, tokens, or passwords.

## Backend Service Creation

The production environment starts with `enable_backend_service = false` so that the first apply can create APIs, Artifact Registry, service accounts, secrets, and Workload Identity Federation before a production container image and secret versions exist.

Once the backend image has been pushed and required secret versions exist, set:

```hcl
enable_backend_service = true
```

Cloud Run image changes are ignored by Terraform after creation so CI/CD can deploy new revisions without Terraform reverting the image.

## CI

Pull requests run Terraform formatting and validation without a remote backend:

```bash
terraform fmt -check -recursive infra/terraform
terraform -chdir=infra/terraform/bootstrap init -backend=false
terraform -chdir=infra/terraform/bootstrap validate
terraform -chdir=infra/terraform/environments/prod init -backend=false
terraform -chdir=infra/terraform/environments/prod validate
```

Production plans require GCP authentication and are enabled after the initial bootstrap/prod apply creates Workload Identity Federation. Configure these GitHub repository variables:

```text
GCP_WORKLOAD_IDENTITY_PROVIDER
GCP_TERRAFORM_SERVICE_ACCOUNT
```

The plan job is skipped until both variables exist. `terraform apply` remains manual.
