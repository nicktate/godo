package godo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Create(t *testing.T) {
	setup()
	defer teardown()

	createdAt, err := time.Parse(time.RFC3339, "2020-01-24T20:24:31Z")
	require.NoError(t, err)
	want := &Registry{
		Name:      "foo",
		CreatedAt: createdAt,
	}

	createRequest := &RegistryCreateRequest{
		Name: want.Name,
	}

	createResponseJSON := `
{
	"registry": {
		"name": "foo",
        "created_at": "2020-01-24T20:24:31Z"
	}
}`

	mux.HandleFunc("/v2/registry", func(w http.ResponseWriter, r *http.Request) {
		v := new(RegistryCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		require.Equal(t, v, createRequest)
		fmt.Fprint(w, createResponseJSON)
	})

	got, _, err := client.Registry.Create(ctx, createRequest)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestRegistry_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Registry{
		Name: "foo",
	}

	getResponseJSON := `
{
	"registry": {
		"name": "foo"
	}
}`

	mux.HandleFunc("/v2/registry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, getResponseJSON)
	})
	got, _, err := client.Registry.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestRegistry_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/registry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Registry.Delete(ctx)
	require.NoError(t, err)
}

func TestRegistry_DockerCredentials(t *testing.T) {
	returnedConfig := "this could be a docker config"
	tests := []struct {
		name              string
		params            *RegistryDockerCredentialsRequest
		expectedReadWrite string
	}{
		{
			name:              "read-only (default)",
			params:            &RegistryDockerCredentialsRequest{},
			expectedReadWrite: "false",
		},
		{
			name:              "read/write",
			params:            &RegistryDockerCredentialsRequest{ReadWrite: true},
			expectedReadWrite: "true",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup()
			defer teardown()

			mux.HandleFunc("/v2/registry/docker-credentials", func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, test.expectedReadWrite, r.URL.Query().Get("read_write"))
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, returnedConfig)
			})

			got, _, err := client.Registry.DockerCredentials(ctx, test.params)
			require.NoError(t, err)
			require.Equal(t, []byte(returnedConfig), got.DockerConfigJSON)
		})
	}
}

func TestRepository_ListTags(t *testing.T) {
	setup()
	defer teardown()

	wantTags := []*RepositoryTag{
		{
			RegistryName:        "foo",
			Repository:          "repo-name",
			Tag:                 "latest",
			ManifestDigest:      "sha256:acd3ca9941a85e8ed16515bfc5328e4e2f8c128caa72959a58a127b7801ee01f",
			CompressedSizeBytes: 2789669,
			SizeBytes:           5843968,
			UpdatedAt:           time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	getResponseJSON := `
{
	"tags": [
		{
			"registry_name": "foo",
			"repository": "repo-name",
			"tag": "latest",
			"manifest_digest": "sha256:acd3ca9941a85e8ed16515bfc5328e4e2f8c128caa72959a58a127b7801ee01f",
			"compressed_size_bytes": 2789669,
			"size_bytes": 5843968,
			"updated_at": "2020-04-01T00:00:00Z"
		}
	],
	"links": {
	    "pages": {
			"next": "https://api.digitalocean.com/v2/registry/foo/repositories/repo-name/tags?page=2",
			"last": "https://api.digitalocean.com/v2/registry/foo/repositories/repo-name/tags?page=2"
		}
	},
	"meta": {
	    "total": 2
	}
}`

	mux.HandleFunc("/v2/registry/foo/repositories/repo-name/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, map[string]string{"page": "1", "per_page": "1"})
		fmt.Fprint(w, getResponseJSON)
	})
	got, response, err := client.Registry.ListRepositoryTags(ctx, &RepositoryListTagsRequest{RegistryName: "foo", Repository: "repo-name"}, &ListOptions{Page: 1, PerPage: 1})
	require.NoError(t, err)
	require.Equal(t, wantTags, got)

	gotRespLinks := response.Links
	wantRespLinks := &Links{
		Pages: &Pages{
			Next: "https://api.digitalocean.com/v2/registry/foo/repositories/repo-name/tags?page=2",
			Last: "https://api.digitalocean.com/v2/registry/foo/repositories/repo-name/tags?page=2",
		},
	}
	assert.Equal(t, wantRespLinks, gotRespLinks)

	gotRespMeta := response.Meta
	wantRespMeta := &Meta{
		Total: 2,
	}
	assert.Equal(t, wantRespMeta, gotRespMeta)
}
