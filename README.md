# Pipeflow

Pipeflow get data from Cloud Pub/Sub and exports to Google Cloud Storage.

## Use Case

- Logging Router -> Cloud Pub/Sub -> Pipeflow -> GCS

If you want to get only error log, you can filter by using Logging Router, and publish them to Cloud Pub/Sub.
Pipeflow can get theses logs from Pub/Sub and exports to GCS like the Data Flow.

## Variables

 Name | Description
 --- | ---
 GCP_PROJCET_ID | gcp project id
 PUBSUB_SUBSCRIPTION | pubsub subscription
 BUCKET_NAME | gcs bucket name
 BUCKET_PREFIX | prefix directory name of gcs objects

## Service Accounts

You need create a service accounts has Storage Write and Cloud Pub/Sub Subscriber permissions, then deploy pipeflow with that account.

## Introduction for Cloud Run

1. Create Cloud Pub/Sub topic and subscription
2. Create GCS bucket
3. Create Service Account for pipeflow(e.g.: `pipeflow@yourprojectid.iam.gserviceaccount.com`).
    - attach Pub/Sub Subscriber and Storage Object Creator roles.
4. Push pipeflow image

```
# build pipeflow docker image
$ make build/image

# <projcet_id> your gcp project id
$ docker tag pipeflow:latest asia.gcr.io/<project_id>/<repository_name>/pipeflow:latest
$ docker push asia-docker.pkg.dev/<project_id>/<repository_name>/pipeflow:latest
```

5. Deploy to Cloud Run
```
$ gcloud beta run deploy pipeflow \
		--project <project_id> \
		--image="asia-docker.pkg.dev/<project_id>/<repository_name>/pipeflow:latest" \
		--port=8080 \
		--region=asia-northeast1 \
		--platform=managed \
		--memory=512Mi \
		--allow-unauthenticated \
		--service-account=pipeflow@<project_id>.iam.gserviceaccount.com \
		--set-env-vars \
		GCP_PROJECT_ID=<project_id>,PUBSUB_SUBSCRIPTION=<your_pubsub_subscription>,BUCKET_NAME=<bucket_name>,BUCKET_PREFIX=<bucket_prefix>
```

6. Create Logging Router
    - Set the topic when you created at '1' as sync destination.
