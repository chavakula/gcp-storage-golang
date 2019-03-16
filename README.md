# Upload file to Google Cloud Storage using  GOLANG

This example will help you to upload files & create a folder structure in google cloud storage

## Installation

Use the package manager [dep](https://golang.github.io/dep/) to install dependencies for this project.

You can refer [here](https://cloud.google.com/storage/docs/reference/libraries#client-libraries-usage-go) to create service account key for GCP storage

## Usage

```bash
cd src/gcp-storage-golang
go run main.go
```

Output   :
```bash
Upload files to GCP storage in standard (regional) bucket
2019/03/16 23:38:28 URL: https://storage.googleapis.com/smart-bucket1/2019March16/testlogo.png
2019/03/16 23:38:28 Size: 696279
2019/03/16 23:38:28 MD5: 53053354e7d0307abb9e2f9d039b59b2
```

#### notes
1) Make sure that .env file parameters are filled up correctly.
2) Make sure you have correct $GOPATH, refer [here](https://github.com/golang/go/wiki/SettingGOPATH) how to set it.