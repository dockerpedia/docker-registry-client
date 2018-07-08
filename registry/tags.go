package registry

type tagsResponse struct {
	Tags []string `json:"tags"`
}

// Tags returns repository tags, ie: latest, 1.0.0, 2.0.0, etc..
func (registry *Registry) Tags(repository string) (tags []string, err error) {
	url := registry.url("/v2/%s/tags/list", repository)

	var response tagsResponse
	for {
		registry.Logf("registry.tags url=%s repository=%s", url, repository)
		url, err = registry.getPaginatedJson(url, &response)
		switch err {
		case ErrNoMorePages:
			tags = append(tags, response.Tags...)
			return tags, nil
		case nil:
			tags = append(tags, response.Tags...)
			continue
		default:
			return nil, err
		}
	}
}

func (registry *Registry) TagSize(repository, reference string) (size int64, err error) {
	deserialized, err := registry.ManifestV2(repository, reference)
	if err != nil {
		return -1, err
	}
	size = int64(0)

	for _, layer := range deserialized.Layers {
		size += layer.Size
	}
	return size, nil
}
