project_id = "job-hunting-saas"
region     = "asia-northeast1"

frontend_domain = "entre.kamiriku.com"
backend_domain  = "api.entre.kamiriku.com"

github_repository    = "karimiku/job-hunting-saas"
github_repository_id = "1162289724"

enable_backend_service   = true
enable_domain_mapping    = true
enable_github_deploy_wif = true

github_actions_workload_identity_pool_id = "github-actions"

backend_container_image = "asia-northeast1-docker.pkg.dev/job-hunting-saas/entre/backend:bootstrap"

database_url_secret_version = "1"

cors_allowed_origins = [
  "https://entre.kamiriku.com",
]
