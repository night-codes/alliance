package main //aliance

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Make new script assembly
func Make(path string) {
	uris := map[string]string{}
	modules := map[string]string{}
	if _, err := os.Stat(path); err == nil {
		filepath.Walk(path, func(filename string, info os.FileInfo, err error) error {
			if !info.IsDir() && err == nil && filepath.Ext(filename) == ".js" {
				if rel, err := filepath.Rel(path, filename); err == nil {
					rel = strings.TrimSuffix(rel, ".js")
					if bfile, err := ioutil.ReadFile(filename); err == nil {
						modules[rel] = string(bfile)
						uris[rel] = filename
					}
				}
			}
			return nil
		})
	}
	t := template.Must(template.New("tpl").Parse(tpl))
	err := t.Execute(os.Stdout, map[string]interface{}{
		"uris":    uris,
		"modules": modules,
	})
	if err != nil {
		log.Println("executing template:", err)
	}

}

func main() {
	Make("test")
}
