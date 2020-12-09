# Automated Tag Creator
The backend for Automated Tag Creator

## Setup the Github App
1. Add `Webhook URL` [https://github.com/settings/apps/automated-tag-creator](https://github.com/settings/apps/automated-tag-creator)
2. Generate a private key

## Deploy the backend
### Add pem data to KMS
Check that the kms api is enabled: [cloudkms.googleapis.com](https://console.developers.google.com/apis/library/cloudkms.googleapis.com).
1. Create a keyring
```shell script
gcloud kms keyrings create atc-secrets --location=global
```
2. Create a key
```shell script
gcloud kms keys create gh-pem-secret \
    --location=global \
    --keyring atc-secrets \
    --purpose encryption
```
### Encrypting a Github Private key
```shell script
gcloud kms encrypt \
    --plaintext-file=gh.pem \
    --ciphertext-file=ghpem.enc.txt \
    --location=global \
    --keyring=atc-secrets \
    --key=gh-pem-secret

base64 ghpem.enc.txt -w 0 > ghpem.enc.64.txt
```
### Grant permissions
Add to `serviceAccount` (1007563553609@cloudbuild.gserviceaccount.com) permissions:
```
Cloud Build Service Account
Cloud KMS CryptoKey Decrypter
Storage Object Viewer
```
### Deployment
```shell script
gcloud builds submit --config cloudbuild.yaml
```