# Terraform Bootstrap

This root module creates the Terraform foundation that must exist before the production root can run from CI:

- GCS bucket for production Terraform remote state
- GitHub Actions Workload Identity Pool / Provider
- `github-terraform` service account
- state bucket access for `github-terraform`
- project roles required for Terraform to manage the production infrastructure

Run this once with local state:

```bash
terraform init
terraform apply
```

If API enablement fails on a fresh project, enable the bootstrap APIs once:

```bash
gcloud services enable serviceusage.googleapis.com cloudresourcemanager.googleapis.com iam.googleapis.com --project=job-hunting-saas
```

After apply, copy these outputs into GitHub repository variables:

```text
GCP_WORKLOAD_IDENTITY_PROVIDER = workload_identity_provider_name
GCP_TERRAFORM_SERVICE_ACCOUNT  = terraform_service_account_email
```

Then initialize `../environments/prod`, which uses this bucket as its `gcs` backend.
