package helpers

import (
	"encoding/json"
	"fmt"
	"gateway-api/domain"
	"io"
	"net/http"
	"strings"
)

func ConfigToMap(config domain.Config) map[string]domain.RouteService {
	mapData := make(map[string]domain.RouteService)
	for _, service := range config.Services {
		for _, route := range service.Routes {
			for _, path := range route.Path {
				mapData[path] = domain.RouteService{
					Name:       service.Name + "." + route.Name,
					Auth:       route.Auth,
					BaseURL:    service.BaseURL,
					Path:       route.ServicePath,
					Method:     route.ServiceMethod,
					TokenField: route.TokenField,
				}
			}
		}
	}
	return mapData
}

func SimpleRequest(method, url string, requestOption domain.RequestPayload) (body []byte, status int, err error) {
	var payload *strings.Reader
	if requestOption.Data != nil {
		payloadByte, _ := json.Marshal(requestOption.Data)
		payload = strings.NewReader(string(payloadByte))
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// parsing query
	q := req.URL.Query()
	for queryKey, queryValue := range requestOption.Query {
		q.Add(queryKey, queryValue)
	}
	req.URL.RawQuery = q.Encode()

	// set header
	req.Header = requestOption.Headers

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	return body, res.StatusCode, nil
}
