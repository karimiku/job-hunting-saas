terraform {
  backend "gcs" {
    bucket = "job-hunting-saas-tfstate"
    prefix = "environments/prod"
  }
}
