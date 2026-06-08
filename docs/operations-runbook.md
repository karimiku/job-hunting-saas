# Operations Runbook

Status: active for beta
Created: 2026-06-08

## Purpose

公開β運用で、個人開発として許容できない課金や障害が起きたときに迷わず止血するための手順をまとめる。

このRunbookは、まず人間が判断して止める運用を前提にする。Budget alertから自動停止する仕組みは、Pub/SubやCloud Functionsなど追加の運用対象と微小な課金要素が増えるため、公開β初期では採用しない。

## Cost Guardrail

### Budget Alert

Google Cloud Budget alertを設定する。Budget alert自体は通知のための設定であり、Cloud RunやCloud Functionsのような常時実行リソースではない。

推奨初期値:

```text
Billing account: 01F9CE-57B873-BD64C4
Project: job-hunting-saas
Monthly budget: 1000 JPY
Alert thresholds:
  - 50%
  - 80%
  - 100%
  - forecasted 100%
```

CLIで作成する場合:

```bash
gcloud billing budgets create \
  --billing-account=01F9CE-57B873-BD64C4 \
  --display-name="job-hunting-saas beta monthly budget" \
  --budget-amount=1000JPY \
  --calendar-period=month \
  --filter-projects=projects/job-hunting-saas \
  --threshold-rule=percent=0.50 \
  --threshold-rule=percent=0.80 \
  --threshold-rule=percent=1.00 \
  --threshold-rule=percent=1.00,basis=forecasted-spend
```

確認:

```bash
gcloud billing budgets list \
  --billing-account=01F9CE-57B873-BD64C4
```

注意:

- Budget alertは自動停止ではなく通知である。
- 100%通知が届いた時点で、実際の課金額がすでに1000円を超えている可能性がある。
- 自動停止を入れる場合は、別途「Pub/Sub + 停止用処理」の設計レビューを行う。

## Cost Alert Response

### 1. 状況確認

まず、どのサービスで費用が出ているか確認する。

```bash
gcloud billing projects describe job-hunting-saas
```

Google Cloud Console:

```text
Billing > Reports
Billing > Budgets & alerts
Logging > Logs Explorer
Cloud Run > Services
Artifact Registry > Repositories
```

確認観点:

- Cloud Runへの想定外アクセスが増えていないか
- Artifact Registryに大きなimageが溜まっていないか
- Cloud Loggingが大量に出ていないか
- Secret ManagerやGCS以外の想定外サービスが有効になっていないか
- Firebase側でFirestore / Storage / Functionsなど未使用サービスが動いていないか

### 2. まずトラフィックを止める

Cloud Run backendが公開されている場合は、外部から叩けないようにする。

```bash
gcloud run services remove-iam-policy-binding entre-backend \
  --project=job-hunting-saas \
  --region=asia-northeast1 \
  --member=allUsers \
  --role=roles/run.invoker
```

Terraformで恒久反映する場合は、`infra/terraform/environments/prod/terraform.tfvars` を以下に戻してapplyする。

```hcl
enable_backend_service = false
enable_domain_mapping  = false
```

```bash
cd /Users/kamiriku/my_projects/job-hunting-saas/infra/terraform/environments/prod
terraform plan -out prod-stop.tfplan
terraform apply prod-stop.tfplan
```

### 3. 高コストになりやすいものを削る

Artifact Registryの不要imageを確認する。

```bash
gcloud artifacts docker images list \
  asia-northeast1-docker.pkg.dev/job-hunting-saas/entre \
  --project=job-hunting-saas
```

不要なimageは削除する。

```bash
gcloud artifacts docker images delete IMAGE_URL \
  --project=job-hunting-saas \
  --quiet
```

Cloud Run serviceを緊急削除する場合:

```bash
gcloud run services delete entre-backend \
  --project=job-hunting-saas \
  --region=asia-northeast1
```

注意:

- Terraform管理対象を手動削除した場合、次回 `terraform plan` で差分が出る。
- 後で復旧するなら、Terraform stateと実リソースの差分を必ず確認する。

### 4. 最後の手段としてBillingを外す

想定外課金が止まらない、またはすぐに調査できない場合のみ、プロジェクトからBilling accountを外す。

```bash
gcloud billing projects unlink job-hunting-saas
```

影響:

- Firebase projectはBlaze課金前提の機能が使えなくなる可能性がある。
- Cloud Run、Artifact Registry、GCS、Secret Managerなどの操作に失敗する可能性がある。
- Terraform remote state bucketの読み書きに影響する可能性がある。
- 復旧時はBilling accountを再リンクしてからTerraform planで差分を確認する。

再リンク:

```bash
gcloud billing projects link job-hunting-saas \
  --billing-account=01F9CE-57B873-BD64C4
```

## Normal Recovery

課金アラート対応後、復旧前に必ず確認する。

```bash
cd /Users/kamiriku/my_projects/job-hunting-saas/infra/terraform/environments/prod
terraform plan
```

確認観点:

- Terraformが意図しない再作成をしようとしていないか
- `enable_backend_service` と `enable_domain_mapping` が意図した値か
- Secret Managerのsecret versionが存在するか
- Cloud Run runtime service accountの権限が残っているか

復旧時は、Cloud Runの `/health` を確認してからVercel/Chrome Extension側の動作確認に進む。
