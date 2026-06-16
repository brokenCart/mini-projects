package urlshortener

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	y, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	m := buildMapFromYAML(y)
	mapHandler := MapHandler(m, fallback)
	return mapHandler, nil
}

func JSONHandler(j []byte, fallback http.Handler) (http.HandlerFunc, error) {
	js, err := parseJSON(j)
	if err != nil {
		return nil, err
	}
	m := buildMapFromJSON(js)
	mapHandler := MapHandler(m, fallback)
	return mapHandler, nil
}

func DBHandler(rows *sql.Rows, fallback http.Handler) (http.HandlerFunc, error) {
	r, err := parseDBRows(rows)
	if err != nil {
		return nil, err
	}
	m := buildMapFromDBRows(r)
	mapHandler := MapHandler(m, fallback)
	return mapHandler, nil
}

func parseYAML(yml []byte) ([]YAMLPathToURL, error) {
	yConfigs := []YAMLPathToURL{}
	err := yaml.Unmarshal(yml, &yConfigs)
	if err != nil {
		return nil, err
	}
	return yConfigs, nil
}

func parseJSON(j []byte) ([]JSONPathToURL, error) {
	jConfigs := []JSONPathToURL{}
	err := json.Unmarshal(j, &jConfigs)
	if err != nil {
		return nil, err
	}
	return jConfigs, nil
}

func parseDBRows(rows *sql.Rows) ([]DBPathToURL, error) {
	result := []DBPathToURL{}
	for rows.Next() {
		var obj DBPathToURL
		err := rows.Scan(&obj.Path, &obj.URL)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

func buildMapFromYAML(yConfigs []YAMLPathToURL) map[string]string {
	resultMap := map[string]string{}
	for _, y := range yConfigs {
		resultMap[y.Path] = y.URL
	}
	return resultMap
}

func buildMapFromJSON(jConfigs []JSONPathToURL) map[string]string {
	resultMap := map[string]string{}
	for _, j := range jConfigs {
		resultMap[j.Path] = j.URL
	}
	return resultMap
}

func buildMapFromDBRows(dbConfigs []DBPathToURL) map[string]string {
	resultMap := map[string]string{}
	for _, r := range dbConfigs {
		resultMap[r.Path] = r.URL
	}
	return resultMap
}

type YAMLPathToURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type JSONPathToURL struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

type DBPathToURL struct {
	Path string
	URL  string
}
