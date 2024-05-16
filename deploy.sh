#!/bin/bash
set -e

read -p "gcp project id: " gcp_project_id
read -p "cloud run region (default: asia-northeast1): " region
read -p "repository name: " repository_name
read -p "your repository region (default: asia): " repo_region
read -p "cloud pubsub subscription: " pubsub_sub
read -p "gcs bucket name: " bucket_name

region=${region:-asia-northeast1}
repo_region=${repo_region:-asia}

gcloud beta run deploy pipeflow \
  --min-instances 1 \
		--project $gcp_project_id \
		--image="$repo_region-docker.pkg.dev/$gcp_project_id/$repository_name/pipeflow:latest" \
		--port=8080 \
		--region=$region \
		--platform=managed \
		--memory=512Mi \
		--allow-unauthenticated \
		--service-account=pipeflow@${gcp_project_id}.iam.gserviceaccount.com \
		--set-env-vars \
		GCP_PROJECT_ID=${gcp_project_id},PUBSUB_SUBSCRIPTION=${pubsub_sub},BUCKET_NAME=${bucket_name},BUCKET_PREFIX=pipeflow
