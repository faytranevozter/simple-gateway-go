package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway-api/domain"
	"io"
	"mime/multipart"
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
	var payload io.Reader
	if requestOption.Data != nil {
		payloadByte, _ := json.Marshal(requestOption.Data)
		payload = strings.NewReader(string(payloadByte))
	}

	var req *http.Request

	if requestOption.Multipart != nil {
		multipartBody := new(bytes.Buffer)
		writer := multipart.NewWriter(multipartBody)

		// loop through all uploaded files
		for fieldName, fs := range requestOption.Multipart.File {
			for _, f := range fs {
				ff, _ := writer.CreateFormFile(fieldName, f.Filename)
				x, errFile := f.Open()
				if errFile != nil {
					fmt.Println(errFile)
					return
				}
				defer x.Close()

				byteFile, errRead := io.ReadAll(x)
				if errRead != nil {
					fmt.Println(errRead)
					return
				}

				// write
				ff.Write(byteFile)
			}
		}

		// regular form (text)
		for fieldName, fs := range requestOption.Multipart.Value {
			for _, val := range fs {
				writer.WriteField(fieldName, val)
			}
		}

		// Close the multipart writer
		writer.Close()

		req, err = http.NewRequest(method, url, multipartBody)
		if err != nil {
			fmt.Println(err)
			return
		}

		// set header
		req.Header = requestOption.Headers
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest(method, url, payload)
		if err != nil {
			fmt.Println(err)
			return
		}

		// set header
		req.Header = requestOption.Headers
	}

	// parsing query
	q := req.URL.Query()
	for queryKey, queryValue := range requestOption.Query {
		q.Add(queryKey, queryValue)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
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
