package urlshortener

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// Return an http.HandlerFunc which maps the paths to their URLs using a map
	// and calls the fallback handler if the path is not found in the map.
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// Parse the YAML and build a map of paths to URLs, then return a MapHandler using that map.
	// If there is an error parsing the YAML, return an error.
	y, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	m := buildMap(y)
	mapHandler := MapHandler(m, fallback)
	return mapHandler, nil
}

func JSONHandler(j []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// Parse the JSON and build a map of paths to URLs, then return a MapHandler using that map.
	// If there is an error parsing the JSON, return an error.
	js, err := parseJSON(j)
	if err != nil {
		return nil, err
	}
	m := buildMap(js)
	mapHandler := MapHandler(m, fallback)
	return mapHandler, nil
}

func DBHandler(db *sql.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var url string

		// Get the url from the DB
		query := `SELECT url FROM mappings WHERE path=?`
		err := db.QueryRow(query, path).Scan(&url)
		if err == nil {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func parseYAML(yml []byte) ([]PathMapping, error) {
	// Parse the YAML into a slice of YAMLPathToURL structs.
	configs := []PathMapping{}
	err := yaml.Unmarshal(yml, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func parseJSON(j []byte) ([]PathMapping, error) {
	// Parse the JSON into a slice of JSONPathToURL structs.
	configs := []PathMapping{}
	err := json.Unmarshal(j, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func buildMap(configs []PathMapping) map[string]string {
	// Build a map from the slice of PathMapping structs.
	resultMap := map[string]string{}
	for _, c := range configs {
		resultMap[c.Path] = c.URL
	}
	return resultMap
}

type PathMapping struct {
	Path string `json:"path" yaml:"path"`
	URL  string `json:"url" yaml:"url"`
}
