pixiv-notifications-to-slack
==============================

Your notifications on pixiv to Slack with Cloud Functions

## Deploy

```sh
cp .env.yaml.example .env.yaml
vi .env.yaml  # Put env variables
gcloud functions deploy PixivNotificationsToSlack --runtime go111 --trigger-http --env-vars-file .env.yaml --region asia-northeast1
```

