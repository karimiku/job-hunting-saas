# Terraform Bootstrap

This root module creates the GCS bucket used by production Terraform remote state.

Run this once with local state:

```bash
terraform init
terraform apply
```

Then initialize `../environments/prod`, which uses this bucket as its `gcs` backend.
