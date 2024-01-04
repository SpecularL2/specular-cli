package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"

	"github.com/sirupsen/logrus"

	"github.com/SpecularL2/specular-cli/internal/service/config"
)

const defaultRepository = "specularL2/specular"
const defaultWorkspacePath = "config/local_devnet"
const githubUrl = "https://api.github.com/repos/%s/contents/%s"

type WorkspaceHandler struct {
	cfg *config.Config
	log *logrus.Logger
}

type ConfigFile struct {
	DownloadUrl string `json:"download_url"`
	Name        string `json:"name"`
}

func (w *WorkspaceHandler) DownloadDefault() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	dst := fmt.Sprintf("%s/.spc/workspaces/default", usr.HomeDir)
	// TODO: ask for confirmation if workspace already exists
	if err = os.RemoveAll(dst); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, os.ModePerm); err != nil {
		return err

	}
	w.log.Infof("saving workspace at: %s", dst)

	orig := fmt.Sprintf(githubUrl, defaultRepository, defaultWorkspacePath)
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
	if w.cfg.WorkspaceCmd.Command == "download" && w.cfg.WorkspaceCmd.Name == "default" {
		if err := w.DownloadDefault(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("the only supported command is:\nspc workspace download default")
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

func NewWorkspaceHandler(cfg *config.Config, log *logrus.Logger) *WorkspaceHandler {
	return &WorkspaceHandler{
		cfg: cfg,
		log: log,
	}
}
