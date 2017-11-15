package main

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	storage "google.golang.org/api/storage/v1"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"cloud.google.com/go/compute/metadata"
	"strconv"
)

type gcp struct {
	project, region, igmName, bucket, object, targetPath string
	client                                               *http.Client
}

func newGcp(project, region, igmName, bucket, object, targetPath string) *gcp {
	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope, storage.DevstorageFullControlScope)
	if err != nil {
		//...
		log.Fatal(err)
	}
	return &gcp{
		project:    project,
		region:     region,
		igmName:    igmName,
		bucket:     bucket,
		object:     object,
		client:     client,
		targetPath: targetPath,
	}
}

func (c gcp) elected() bool {
	computeService, err := compute.New(c.client)
	if err != nil {
		//...
		log.Fatal(err)
	}

	instanceGroup, err := computeService.RegionInstanceGroupManagers.ListManagedInstances(c.project, c.region, c.igmName).Do()
	if err != nil {
		//...
		log.Fatal(err)
	}

	// just select the machine with lower id
	electedId := uint64(math.MaxUint64)
	for _, i := range instanceGroup.ManagedInstances {
		if i.Id < electedId {
			electedId = i.Id
		}
	}
	fmt.Println("Elected machine with id:", electedId)
	//if it's not the current machine return false and go home
	//currentIdString, _ := metadata.InstanceID()
	//currentId, _ := strconv.ParseInt(currentIdString, 10, 64)
	//if electedId != uint64(currentId) {
	//	return false
	//}
	return true
}

func (c gcp) retrieveAssets() {
	storageService, err := storage.New(c.client)
	if err != nil {
		//...
		log.Fatal(err)
	}

	// retrieve assets
	resp, err := storageService.Objects.Get(c.bucket, c.object).Download()
	if err != nil {
		//...
		log.Fatal(err)
	}
	fmt.Printf("Object %v retrieved from bucket %v\n", c.object, c.bucket)
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//...
		log.Fatal(err)
	}

	// store assets
	err = ioutil.WriteFile(c.targetPath, buf, 0644)
	if err != nil {
		//...
		log.Fatal(err)
	}
	fmt.Println("Assets stored into:", c.targetPath)
	//storageService.Objects.Delete(bucket, object).Do()
}
