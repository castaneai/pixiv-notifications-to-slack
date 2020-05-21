deploy:
	gcloud functions deploy PixivNotificationsToSlack --runtime go113 --trigger-http --env-vars-file .env.yaml --region asia-northeast1 --allow-unauthenticated