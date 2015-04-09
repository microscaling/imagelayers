package registry_test

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/CenturyLinkLabs/docker-reg-client/registry"
)

func ExampleRepositoryService_ListTags_dockerHub() {
	c := registry.NewClient()
	auth, err := c.Hub.GetReadToken("ubuntu")

	tags, err := c.Repository.ListTags("ubuntu", auth)
	if err != nil {
		panic(err)
	}

	for tag, imageID := range tags {
		fmt.Printf("%s = %s", tag, imageID)
	}
}

func ExampleRepositoryService_ListTags_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	tags, err := c.Repository.ListTags("ubuntu", auth)
	if err != nil {
		panic(err)
	}

	for tag, imageID := range tags {
		fmt.Printf("%s = %s", tag, imageID)
	}
}

func ExampleRepositoryService_DeleteTag_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	err := c.Repository.DeleteTag("foo/bar", "1.0", auth)
	if err != nil {
		panic(err)
	}
}

func ExampleRepositoryService_GetImageID_dockerHub() {
	c := registry.NewClient()
	auth, err := c.Hub.GetReadToken("ubuntu")

	imageID, err := c.Repository.GetImageID("ubuntu", "latest", auth)
	if err != nil {
		panic(err)
	}

	fmt.Println(imageID)
}

func ExampleRepositoryService_GetImageID_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	imageID, err := c.Repository.GetImageID("ubuntu", "latest", auth)
	if err != nil {
		panic(err)
	}

	fmt.Println(imageID)
}

func ExampleRepositoryService_Delete_dockerHub() {
	c := registry.NewClient()
	basicAuth := registry.BasicAuth{"user", "pass"}
	tokenAuth, err := c.Hub.GetDeleteToken("foo/bar", basicAuth)

	err = c.Repository.Delete("foo/bar", tokenAuth)
	if err != nil {
		panic(err)
	}
}

func ExampleRepositoryService_Delete_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	err := c.Repository.Delete("foo/bar", auth)
	if err != nil {
		panic(err)
	}
}

func ExampleRepositoryService_SetTag_dockerHub() {
	c := registry.NewClient()
	basicAuth := registry.BasicAuth{"user", "pass"}
	tokenAuth, err := c.Hub.GetWriteToken("foo/bar", basicAuth)

	imageID := "2427658c75a1e3d0af0e7272317a8abfaee4c15729b6840e3c2fca342fe47bf1"
	err = c.Repository.SetTag("foo/bar", imageID, "1.0", tokenAuth)
	if err != nil {
		panic(err)
	}
}

func ExampleRepositoryService_SetTag_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	imageID := "2427658c75a1e3d0af0e7272317a8abfaee4c15729b6840e3c2fca342fe47bf1"
	err := c.Repository.SetTag("foo/bar", imageID, "1.0", auth)
	if err != nil {
		panic(err)
	}
}

func ExampleImageService_GetAncestry_dockerHub() {
	c := registry.NewClient()
	auth, err := c.Hub.GetReadToken("ubuntu")

	images, err := c.Image.GetAncestry("ubuntu", auth)
	if err != nil {
		panic(err)
	}

	for _, imageID := range images {
		fmt.Println(imageID)
	}
}

func ExampleImageService_GetAncestry_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	images, err := c.Image.GetAncestry("ubuntu", auth)
	if err != nil {
		panic(err)
	}

	for _, imageID := range images {
		fmt.Println(imageID)
	}
}

func ExampleImageService_GetMetadata_dockerHub() {
	c := registry.NewClient()
	auth, err := c.Hub.GetReadToken("ubuntu")

	imageID := "2427658c75a1e3d0af0e7272317a8abfaee4c15729b6840e3c2fca342fe47bf1"
	meta, err := c.Image.GetMetadata(imageID, auth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", meta)
}

func ExampleImageService_GetMetadata_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	imageID := "2427658c75a1e3d0af0e7272317a8abfaee4c15729b6840e3c2fca342fe47bf1"
	meta, err := c.Image.GetMetadata(imageID, auth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", meta)
}

func ExampleImageService_GetLayer_dockerHub() {
	c := registry.NewClient()
	auth, err := c.Hub.GetReadToken("ubuntu")

	buffer := &bytes.Buffer{}
	imageID := "2427658c75a1e3d0af0e7272317a8abfaee4c15729b6840e3c2fca342fe47bf1"
	err = c.Image.GetLayer(imageID, buffer, auth)
	if err != nil {
		panic(err)
	}

	fmt.Println(buffer)
}

func ExampleImageService_GetLayer_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")
	auth := registry.BasicAuth{"user", "pass"}

	buffer := &bytes.Buffer{}
	imageID := "2427658c75a1e3d0af0e7272317a8abfaee4c15729b6840e3c2fca342fe47bf1"
	err := c.Image.GetLayer(imageID, buffer, auth)
	if err != nil {
		panic(err)
	}

	fmt.Println(buffer)
}

func ExampleSearchService_Query_dockerHub() {
	c := registry.NewClient()

	results, err := c.Search.Query("mysql", 1, 25)
	if err != nil {
		panic(err)
	}

	for _, result := range results.Results {
		fmt.Println(result.Name)
	}
}

func ExampleSearchService_Query_privateRegistry() {
	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")

	results, err := c.Search.Query("mysql", 1, 25)
	if err != nil {
		panic(err)
	}

	for _, result := range results.Results {
		fmt.Println(result.Name)
	}
}
