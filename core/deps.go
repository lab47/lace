package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func externalHttpSourceToPath(env *Env, lib string, url string) (string, error) {
	home, _ := os.LookupEnv("HOME")
	localBase := filepath.Join(home, ".laced", "deps", strings.SplitN(url, "//", 2)[1])
	libBase := filepath.Join(strings.Split(lib, ".")...) + ".clj"
	libPath := filepath.Join(localBase, libBase)
	libPathDir := filepath.Dir(libPath)

	if _, err := os.Stat(libPathDir); os.IsNotExist(err) {
		err := os.MkdirAll(libPathDir, 0777)
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		if !strings.HasSuffix(url, ".clj") {
			url = url + libBase
		}
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", env.NewError(fmt.Sprintf("Unable to retrieve: %s\nServer response: %d", url, resp.StatusCode))
		}

		out, err := os.Create(libPath)
		if err != nil {
			return "", err
		}

		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", err
		}

		err = out.Close()
		if err != nil {
			return "", err
		}
	}

	return libPath, nil
}

func externalSourceToPath(env *Env, lib string, url string) (string, error) {
	httpPath, _ := regexp.MatchString("http://|https://", url)
	if httpPath {
		return externalHttpSourceToPath(env, lib, url)
	} else {
		return filepath.Join(append([]string{url}, strings.Split(lib, ".")...)...) + ".clj", nil
	}
}
