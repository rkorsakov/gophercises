package urlshort

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
)

type PathURL struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if targetURL, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, targetURL, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) (parsedYaml []PathURL, err error) {
	var pathURLs []PathURL
	err = yaml.Unmarshal(yml, &pathURLs)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	return pathURLs, nil
}

func buildMap(parsedYaml []PathURL) map[string]string {
	pathURLMap := make(map[string]string)
	for _, pathURL := range parsedYaml {
		pathURLMap[pathURL.Path] = pathURL.Url
	}
	return pathURLMap
}
