module "backend" {
  source = "../../modules/backend"

  project_id = var.project_id
  region     = var.region

  frontend_domain = var.frontend_domain
  backend_domain  = var.backend_domain

  github_repository    = var.github_repository
  github_repository_id = var.github_repository_id

  artifact_repository_id = "entre"
  backend_service_name   = "entre-backend"

  enable_backend_service  = var.enable_backend_service
  enable_domain_mapping   = var.enable_domain_mapping
  backend_container_image = var.backend_container_image

  database_url_secret_version = var.database_url_secret_version

  cors_allowed_origins = var.cors_allowed_origins
}
