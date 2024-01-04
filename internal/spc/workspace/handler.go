package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/exp/maps"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

const githubUrl = "https://api.github.com/repos/%s/contents/%s"

type WorkspaceHandler struct {
	cfg *config.Config
	log *logrus.Logger
}

type ConfigFile struct {
	DownloadUrl string `json:"download_url"`
	Name        string `json:"name"`
}

func (w *WorkspaceHandler) DownloadConfig() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	dst := fmt.Sprintf("%s/.spc/workspaces/%s", usr.HomeDir, w.cfg.WorkspaceCmd.Name)
	// TODO: ask for confirmation if workspace already exists
	if err = os.RemoveAll(dst); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, os.ModePerm); err != nil {
		return err

	}
	w.log.Infof("saving workspace at: %s", dst)

	orig := fmt.Sprintf(githubUrl, w.cfg.WorkspaceCmd.ConfigRepo, w.cfg.WorkspaceCmd.ConfigPath)
	resp, err := http.Get(orig)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			w.log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var files []ConfigFile
	err = json.Unmarshal(body, &files)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := dst + "/" + file.Name
		if err = w.downloadFile(filePath, file.DownloadUrl); err != nil {
			return err
		}
	}

	return nil
}

func (w *WorkspaceHandler) Cmd() error {
	if w.cfg.WorkspaceCmd.Command == "download" {
		if err := w.DownloadConfig(); err != nil {
			return err
		}
	} else if w.cfg.WorkspaceCmd.Command == "activate" {
		if err := w.LoadWorkspaceEnvVars(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("the only supported commands are [download] and [activate]")
	}
	return nil
}

func (w *WorkspaceHandler) downloadFile(filepath string, url string) error {
	w.log.Tracef("donwloading file: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			w.log.Fatal(err)
		}
	}(resp.Body)

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			w.log.Fatal(err)
		}
	}(out)

	w.log.Tracef("saving at: %s\n", filepath)
	_, err = io.Copy(out, resp.Body)
	return err
}

func (w *WorkspaceHandler) LoadWorkspaceEnvVars() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	src := fmt.Sprintf("%s/.spc/workspaces/%s", usr.HomeDir, w.cfg.WorkspaceCmd.Name)

	items, _ := os.ReadDir(src)
	envVars := map[string]string{}

	for _, item := range items {
		fullItemPath := fmt.Sprintf("%s/%s", src, item.Name())
		isDotEnvFile := !item.IsDir() && strings.HasPrefix(item.Name(), ".") && strings.HasSuffix(item.Name(), ".env")

		isJSONFile := !item.IsDir() && strings.HasSuffix(item.Name(), ".json")
		if isDotEnvFile {
			w.log.Debugf("found env file: %s ..", fullItemPath)
			vars, err := godotenv.Read(fullItemPath)
			if err != nil {
				w.log.Warnf("failed to load %s", fullItemPath)
			}
			maps.Copy(envVars, vars)
		}

		if isJSONFile {
			w.log.Debugf("found JSON file: %s ..", fullItemPath)

			content, err := os.ReadFile(fullItemPath)
			if err != nil {
				return err
			}

			var b map[string]interface{}
			err = json.Unmarshal(content, &b)
			if err != nil {
				return err
			}
			bStr := map[string]string{}
			for k, v := range b {
				bStr[k] = fmt.Sprintf("%s", v)
			}

			maps.Copy(envVars, bStr)
		}
	}

	envPrefixVars := map[string]string{}
	for k, v := range envVars {
		key := fmt.Sprintf("SPC_%s", strings.ToUpper(k))
		envPrefixVars[key] = v
	}

	tmp, _ := json.Marshal(envPrefixVars)
	w.log.Infof("loaded vars: %s", tmp)

	return nil
}

func NewWorkspaceHandler(cfg *config.Config, log *logrus.Logger) *WorkspaceHandler {
	return &WorkspaceHandler{
		cfg: cfg,
		log: log,
	}
}
