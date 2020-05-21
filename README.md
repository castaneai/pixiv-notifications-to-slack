pixiv-notifications-to-slack
==============================

Your notifications on pixiv to Slack with Cloud Functions

## Deploy

```sh
cp .env.yaml.example .env.yaml
vi .env.yaml  # Put env variables
make deploy
```

## Test

```sh
gcloud beta emulators firestore start --host-port=localhost:8812
export FIRESTORE_EMULATOR_HOST=localhost:8812
go test ./...
```