# Automated Tag Creator
The backend for Automated Tag Creator
## Developer mode configuration
1. Create app instance in own account (Developer settings/GitHub Apps)
2. Start `ngrok`
3. Use ngrok address for webhook url like: ```https://26680d04b127.ngrok.io/api/webhook```
4. Generate and download private key
5. Configure app permissions:
    - Content: Read-only
    - Discussion: Read-Write
    - MetaData:Read-only
6. Subscribe to events:
    - Create
    - Push
    - Delete
7. Use script like below to start app
    ```bash
    export ATC_PEM_PATH=/home/andrey/.ssh/atc-local.2021-03-25.private-key.pem
    export ATC_APP_ID=106890
    bin/atcapp
    ```

## Create the GitHub App
1. Navigate to your account settings.
2. Go to `Developer settings` -> `GitGub Apps`
3. Click `NewGitHub App`
4. In `GitHub App name`, type the name of your app 
5. In `HomepageURL`, type the full URL to your app's website
6. Cancel select `Webhook -> Active`

More informations you can see in [Creating a GitHub App](https://docs.github.com/en/developers/apps/building-github-apps/creating-a-github-app)

## Configurate and install the GitHub App
1. Navigate to your account settings.
2. Go to `Developer settings` -> `GitHub Apps`
3. Click `Edit` in your App
    - Select `Webhook -> Active`
    - Use ngrok address with *api/webhook* for webhook url like: ```https://26680d04b127.ngrok.io/api/webhook```
    - Click `Generate a private key` and download private key
    - Click `Save changes`
4. Go to `Permissions & events`
    - Configurate `Repository permissions`:
        * `Content`: Read & write
        * `MetaData`: Read-only
    - Select in `subscribe to events`:
        * `Create`
        * `Push`
        * `Delete`
    - Click `Save changes`
5. Go to `Install App`
    - Choose an account to install and click `Install`
    - Choose `All repositories` or  `Only select repositories` and select repositories
    - Click `Install`

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
Also you should follow the steps:
 1. Grant the Cloud Run Admin role to the Cloud Build service account:
     * In the Cloud Console, go to the Cloud Build Settings page:
     * Open the Settings page
     * Locate the row with the Cloud Run Admin role and set its Status to ENABLED.
     * In the Additional steps may be required pop-up, click Skip.

 2. Grant the IAM Service Account User role to the Cloud Build service account on the Cloud Run runtime service account:
     * In the Cloud Console, go to the Service Accounts page:
     * Open the Service Accounts page
     * In the list of members, locate and select [PROJECT_NUMBER]-compute@developer.gserviceaccount.com. This is the Cloud Run runtime service account.
     * Click SHOW INFO PANEL in the top right corner.
     * In the Permissions panel, click the Add Member button.
     * In the New member field, enter the email address of the Cloud Build service account. This is of the form [PROJECT_NUMBER]@cloudbuild.gserviceaccount.com. Note: The email address of Cloud Build service account is different from that of Cloud Run runtime service account.
     * In the Role dropdown, select Service Accounts, and then Service Account User.
     * Click Save.

See more on [Stackoverflow](https://stackoverflow.com/questions/62783869/why-am-i-seeing-this-error-error-gcloud-run-deploy-permission-denied-the-c) 
### Check a default project
```
> gcloud config get-value project
atc-sf
```
### Deployment
```shell script
gcloud builds submit --config cloudbuild.yaml
```
