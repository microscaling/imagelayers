/*
Package registry is an API wrapper for the Docker Registry v1 API and supports
all of the documented operations, including: getting/putting images, retrieving
repository tags, and executing searches.

Full documentation for the Docker Registry API can be found here:

https://docs.docker.com/reference/api/registry_api/

Usage

There are two different usage models for the registry client depending on
whether you are trying to interact with a private Docker registry or the Docker
Hub.

When using the client with a private Docker registry you can simply pass basic
auth credentials to the functions you are invoking and the appropriate
authentication headers will be passed along with the request:

	c := registry.NewClient()
	c.BaseURL, _ = url.Parse("http://my.registry:8080/v1/")

	auth := registry.BasicAuth{"user", "pass"}

	err := c.Repository.SetTag("foo/bar", "abc123", "1.0", auth)

When using the client to interact with the Docker Hub there is an extra call
that you need to make. You need to use your basic auth credentials to first
retrieve a token and then you can use that token when interacting with the
Registry API:

	c := registry.NewClient()

	basicAuth := registry.BasicAuth{"user", "pass"}

	tokenAuth, err := c.Hub.GetWriteToken("foo/bar", basicAuth)

	err := c.Repository.SetTag("foo/bar", "abc123", "1.0", tokenAuth)

Depending on the API you're trying to use, you may need to first retreive either
a read, write or delete token.

For more details about interacting with the Docker Hub, see the Docker Hub and
Registry Spec:

https://docs.docker.com/reference/api/hub_registry_spec/

Note: If your Docker Hub repository has been setup as an Automated Build
repositry you will not be able to make any changes to it (setting tags, etc...)
via the Registry API. You can only make changes to an Automated Build
repository through the Docker Hub UI.
*/
package registry
