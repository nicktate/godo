package godo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	registryPath = "/v2/registry"
	// RegistryServer is the hostname of the DigitalOcean registry service
	RegistryServer = "registry.digitalocean.com"
)

// RegistryService is an interface for interfacing with the Registry endpoints
// of the DigitalOcean API.
// See: https://developers.digitalocean.com/documentation/v2#registry
type RegistryService interface {
	Create(context.Context, *RegistryCreateRequest) (*Registry, *Response, error)
	Get(context.Context) (*Registry, *Response, error)
	Delete(context.Context) (*Response, error)
	DockerCredentials(context.Context, *RegistryDockerCredentialsRequest) (*DockerCredentials, *Response, error)
	ListRepositoryImages(context.Context, *RepositoryListImagesRequest, *ListOptions) ([]*RepositoryImage, *Response, error)
}

var _ RegistryService = &RegistryServiceOp{}

// RegistryServiceOp handles communication with Registry methods of the DigitalOcean API.
type RegistryServiceOp struct {
	client *Client
}

// RegistryCreateRequest represents a request to create a registry.
type RegistryCreateRequest struct {
	Name string `json:"name,omitempty"`
}

// RegistryDockerCredentialsRequest represents a request to retrieve docker
// credentials for a registry.
type RegistryDockerCredentialsRequest struct {
	ReadWrite bool `json:"read_write"`
}

// RepositoryListImagesRequest represents a request to retrieve images
// for a given repository in a registry.
type RepositoryListImagesRequest struct {
	RegistryName string `json:"registry_name,omitempty"`
	Repository   string `json:"repository,omitempty"`
}

type registryRoot struct {
	Registry *Registry `json:"registry,omitempty"`
}

type repositoryImagesRoot struct {
	Images []*RepositoryImage `json:"images,omitempty"`
	Links  *Links             `json:"links,omitempty"`
	Meta   *Meta              `json:"meta"`
}

// Registry represents a registry.
type Registry struct {
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// RepositoryImage represents a repository image
type RepositoryImage struct {
	RegistryName        string    `json:"registry_name,omitempty"`
	Repository          string    `json:"repository,omitempty"`
	Tag                 string    `json:"tag,omitempty"`
	ManifestDigest      string    `json:"manifest_digest,omitempty"`
	CompressedSizeBytes uint64    `json:"compressed_size_bytes,omitempty"`
	SizeBytes           uint64    `json:"size_bytes,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}

// Get retrieves the details of a Registry.
func (svc *RegistryServiceOp) Get(ctx context.Context) (*Registry, *Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodGet, registryPath, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(registryRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Registry, resp, nil
}

// Create creates a registry.
func (svc *RegistryServiceOp) Create(ctx context.Context, create *RegistryCreateRequest) (*Registry, *Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodPost, registryPath, create)
	if err != nil {
		return nil, nil, err
	}
	root := new(registryRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Registry, resp, nil
}

// Delete deletes a registry. There is no way to recover a registry once it has
// been destroyed.
func (svc *RegistryServiceOp) Delete(ctx context.Context) (*Response, error) {
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, registryPath, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DockerCredentials is the content of a Docker config file
// that is used by the docker CLI
// See: https://docs.docker.com/engine/reference/commandline/cli/#configjson-properties
type DockerCredentials struct {
	DockerConfigJSON []byte
}

// DockerCredentials retrieves a Docker config file containing the registry's credentials.
func (svc *RegistryServiceOp) DockerCredentials(ctx context.Context, request *RegistryDockerCredentialsRequest) (*DockerCredentials, *Response, error) {
	path := fmt.Sprintf("%s/%s?read_write=%t", registryPath, "docker-credentials", request.ReadWrite)

	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var buf bytes.Buffer
	resp, err := svc.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, resp, err
	}

	dc := &DockerCredentials{
		DockerConfigJSON: buf.Bytes(),
	}
	return dc, resp, nil
}

// ListRepositoryImages returns a list of the RepositoryImages available within the given repository.
func (svc *RegistryServiceOp) ListRepositoryImages(ctx context.Context, request *RepositoryListImagesRequest, opts *ListOptions) ([]*RepositoryImage, *Response, error) {
	path := fmt.Sprintf("%s/%s/repositories/%s/images", registryPath, request.RegistryName, request.Repository)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(repositoryImagesRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Images, resp, nil
}
