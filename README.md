# Pipeflow

Pipeflow get data from Cloud Pub/Sub and exports to Google Cloud Storage.

## Use Case

- Logging Router -> Cloud Pub/Sub -> Pipeflow -> GCS

If you want to get only error log, you can filter by using Logging Router, and publish them to Cloud Pub/Sub.
Pipeflow can get theses logs from Pub/Sub and exports to GCS like the Data Flow.

## Variables

 --- | ---
 GCP_PROJCET_ID | gcp project id
 PUBSUB_SUBSCRIPTION | pubsub subscription
 BUCKET_NAME | gcs bucket name
 BUCKET_PREFIX | prefix directory name of gcs objects
 HEALTHCHECK_ENDPOINT | Cloud Run endpoint e.g: https://yourservice-xxx.x.run.app

## IAM setting with Service Accounts

You need create a service accounts has Storage Write and Cloud Pub/Sub Subscriber permissions, then deploy pipeflow with that account.
