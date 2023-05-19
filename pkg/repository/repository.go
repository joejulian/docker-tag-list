package repository

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/net/idna"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

func QueryImageRepositoryTags(repo string) ([]string, error) {
	r, err := remote.NewRepository(NormalizeRepository(repo))
	if err != nil {
		return nil, fmt.Errorf("unable to open repository '%s': %v", repo, err)
	}

	ctx := context.Background()
	tags, err := registry.Tags(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("unable to query tags for repository '%s': %v", repo, err)
	}

	return tags, nil
}

func NormalizeRepository(repository string) string {
	parts := strings.Split(repository, "/")
	if len(parts) >= 3 {
		return repository
	}
	if strings.Contains(repository, ".") && !strings.Contains(parts[0], "docker.io") {
		_, err := idna.Lookup.ToASCII(parts[0])
		if err == nil {
			return repository
		}
	}
	if len(parts) == 2 {
		if strings.Contains(parts[0], "docker.io") {
			parts = append([]string{"docker.io", "library"}, parts[1:]...)
			return strings.Join(parts, "/")
		} else {
			parts = append([]string{"docker.io"}, parts...)
			return strings.Join(parts, "/")
		}
	}
	if len(parts) == 1 {
		parts = append([]string{"docker.io", "library"}, parts...)
		return strings.Join(parts, "/")
	}
	return repository
}
