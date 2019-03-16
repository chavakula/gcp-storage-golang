package main

import (
	"fmt"
	"context"
	"github.com/joho/godotenv"
	"cloud.google.com/go/storage"
	"log"
	"os"
	"io"
	"path/filepath"
	"time"
)

/*
	This GOLANG script will upload file to GCP storage, if bucket is not present,
	then it will first create as per .ENV variable defined and then upload file to it.
 */

// ***** PLEASE SET ENV VARIABLES DEFINED IN .env FILE ******

// PROJECTID : project id of google cloud account ,
// GOOGLE_APPLICATION_CREDENTIALS :  JSON of service account key
// BUCKETNAME : bucketname in which you want to store files
// STORAGE_CLASS : default value: STANDARD, other possible values MULTI_REGIONAL, REGIONAL, NEARLINE, COLDLINE, DURABLE_REDUCED_AVAILABILITY
// STORAGE_LOC : location in region where you want to store files

func main(){
	fmt.Println("Upload files to GCP storage in standard (regional) bucket");

	// file to be uploaded
	source := "/Users/rajshekar/Documents/testlogo.png"

	ctx := context.Background();
	err := godotenv.Load();

	// Creates a client.
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalf("unable to create GCP storage client %v", err)
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// create bucket if it doesn't exist
	cerr := createBucket(ctx, client)

	if cerr != nil{
		log.Fatalf("Failed to create bucket: %v", cerr)
	}

	// upload file here
	_, objAttrs, err := uploadToCloud(client, ctx, source, false)

	if err != nil {
		switch err {
		case storage.ErrBucketNotExist:
			log.Fatal("Please create the bucket first e.g. with `gsutil mb`")
		default:
			log.Fatal(err)
		}
	}

	// some basic information after uploading file

	log.Printf("URL: %s", objectURL(objAttrs))
	log.Printf("Size: %d", objAttrs.Size)
	log.Printf("MD5: %x", objAttrs.MD5)
	//log.Printf("objAttrs: %+v", objAttrs)

}


// create & check bucket in GCP storage service
func createBucket(ctx context.Context , client *storage.Client) error  {

	// sets project id
	projectID := os.Getenv("PROJECTID")

	// Sets the name for the new bucket.
	bucketName := os.Getenv("BUCKETNAME")

	// check bucket existence
	exists, err := checkBucketExists(ctx, client, bucketName)

	if err != nil {
		return err
	}

	// if bucket doesn't exist
	if exists == false {

		// Creates a Bucket handle.
		bucket := client.Bucket(bucketName)

		// Creates the new bucket.
		if err := bucket.Create(ctx, projectID, &storage.BucketAttrs{StorageClass: os.Getenv("STORAGE_CLASS"), Location: os.Getenv("STORAGE_LOC")}); err != nil {
			return err
		}
	}
	return err
}


// UPLOAD FILE FUNCTIONALITY
func uploadToCloud(client *storage.Client, ctx context.Context, source string, publicflag bool) (*storage.ObjectHandle, *storage.ObjectAttrs, error){

	// get original filename
	filename := filepath.Base(source)

	// create folder structure in bucket
	tree := folderStructure()

	f, err := os.Open(source)

	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	defer f.Close()

	clientBucket := client.Bucket(os.Getenv("BUCKETNAME"))
	writeObj := clientBucket.Object(tree + "/" + filename)
	w := writeObj.NewWriter(ctx)

	if _, err = io.Copy(w, f); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	// if you want to make uploaded item public
	if publicflag {
		if err := writeObj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil{
			return nil,nil, err
		}

	}

	attrs, err := writeObj.Attrs(ctx)
	return writeObj, attrs, err

}


// FUNCTION TO CHECK IF BUCKET ALREADY EXISTS
func checkBucketExists(ctx context.Context ,client *storage.Client , bucketName string) (bool,error){
	bucket := client.Bucket(bucketName)
	_,err := bucket.Attrs(ctx)

	if err != nil {
		return false,err
	}
	return true,err
}

// FORM URL OF UPLOADED FILE
func objectURL(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}

// FORM FOLDER STRUCTURE
func folderStructure() string{
	currentTime := time.Now()
	return fmt.Sprintf("%v%v%v",currentTime.Year(), currentTime.Month(), currentTime.Day())
}