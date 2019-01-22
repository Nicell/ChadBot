package helpers

import (
	"errors"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func YTsearch(query string, key string) (string, error) {

	client := &http.Client{
		Transport: &transport.APIKey{Key: key},
	}

	service, err := youtube.New(client)
	if err != nil {
		return "", err
	}

	call := service.Search.List("id").Q(query).MaxResults(1)

	response, err := call.Do()
	if err != nil {
		return "", err
	}

	if len(response.Items) > 0 {
		return response.Items[0].Id.VideoId, nil
	}

	return "", errors.New("YouTube query showed 0 results")
}
