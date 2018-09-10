package alliance

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dchest/jsmin"
)

// Make one script from modules directory
func Make(path string, gz bool) (string, error) {
	uris := map[string]string{}
	modules := map[string]string{}
	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	filepath.Walk(path, func(filename string, info os.FileInfo, err error) error {
		if !info.IsDir() && err == nil && filepath.Ext(filename) == ".js" {
			if rel, err := filepath.Rel(path, filename); err == nil {
				rel =  strings.Replace(strings.TrimSuffix(rel, ".js"), "-", "_", -1)
				if bfile, err := ioutil.ReadFile(filename); err == nil {
					modules[rel] = string(bfile)
					uris[rel] = filename
				}
			}
		}
		return nil
	})

	t := template.Must(template.New("tpl").Parse(tpl))
	var err error
	var scriptbufer bytes.Buffer
	var minbytes []byte
	if err = t.Execute(&scriptbufer, map[string]interface{}{"uris": uris, "modules": modules}); err == nil {
		if minbytes, err = jsmin.Minify(scriptbufer.Bytes()); err == nil {
			if gz {
				gzipbufer := bytes.Buffer{}
				gzwriter := gzip.NewWriter(&gzipbufer)
				if _, err = gzwriter.Write(minbytes); err == nil {
					if err = gzwriter.Flush(); err == nil {
						if err = gzwriter.Close(); err == nil {
							return string(gzipbufer.Bytes()), nil
						}
					}
				}
			}
			return string(minbytes), nil
		}
	}
	return "", err
}
