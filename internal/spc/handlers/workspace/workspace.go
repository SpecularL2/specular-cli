package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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

	dst := fmt.Sprintf("%s/.spc/workspaces/%s", usr.HomeDir, w.cfg.Args.Workspace.Name)
	// TODO: ask for confirmation if workspace already exists
	if err = os.RemoveAll(dst); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, os.ModePerm); err != nil {
		return err

	}
	w.log.Infof("saving workspace at: %s", dst)

	orig := fmt.Sprintf(githubUrl, w.cfg.Args.Workspace.Download.ConfigRepo, w.cfg.Args.Workspace.Download.ConfigPath)
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
	switch {
	case w.cfg.Args.Workspace.Download != nil:
		return w.DownloadConfig()
	case w.cfg.Args.Workspace.Activate != nil:
		return w.LoadWorkspaceEnvVars()
	case w.cfg.Args.Workspace.Set != nil:
		return w.SetWorkspace()
	case w.cfg.Args.Workspace.List != nil:
		return w.ListWorkspaces()
	}

	w.log.Warn("no command found, exiting...")
	return nil
}

func (w *WorkspaceHandler) SetWorkspace() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	workspaceDir := fmt.Sprintf("%s/.spc/workspaces/", usr.HomeDir)
	selectedWorkspace := fmt.Sprintf("%s%s", workspaceDir, w.cfg.Args.Workspace.Set.Name)

	if _, err := os.Stat(selectedWorkspace); err != nil {
		w.log.Fatalf("could not find workspace with name: %s", w.cfg.Args.Workspace.Set.Name)
		return nil
	}

	activePath := fmt.Sprintf("%s%s", workspaceDir, "active_workspace")
	if _, err := os.Lstat(activePath); err == nil {
		w.log.Trace("removing existing active workspacet")
		os.Remove(activePath)
	}

	err = os.Symlink(selectedWorkspace, activePath)
	if err != nil {
		return err
	}

	w.log.Infof("set workspace %s as active", w.cfg.Args.Workspace.Set.Name)
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

	activeWorkspace := fmt.Sprintf("%s/.spc/workspaces/active_workspace", usr.HomeDir)
	_, err = os.Stat(activeWorkspace)
	if err != nil {
		w.log.Fatalf("no active workspace set")
		return nil
	}

	items, err := os.ReadDir(activeWorkspace)
	if err != nil {
		return err
	}

	envVars := map[string]string{}

	for _, item := range items {
		fullItemPath := fmt.Sprintf("%s/%s", activeWorkspace, item.Name())
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
		err := os.Setenv(key, v)
		if err != nil {
			w.log.Warnf("could not set env var: %s=%s", key, v)
		}
	}

	tmp, _ := json.Marshal(envPrefixVars)
	w.log.Infof("loaded vars: %s", tmp)

	return nil
}

func (w *WorkspaceHandler) ListWorkspaces() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	src := fmt.Sprintf("%s/.spc/workspaces/", usr.HomeDir)
	items, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			w.log.Infof("\t%s", item.Name())
		}
	}
	return nil
}

// run a string command in the context of the currently active workspace
func (w *WorkspaceHandler) RunStringCommand(cmd string) error {
	// TODO: handle case where there is no active workspace
	err := w.LoadWorkspaceEnvVars()
	if err != nil {
		return err
	}

	commandToRun := os.ExpandEnv(cmd)
	args := strings.Fields(commandToRun)

	if len(args) > 0 {
		w.log.Debugf("Running: %s %v", commandToRun, args)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
		}
	}
	return nil
}
func NewWorkspaceHandler(cfg *config.Config, log *logrus.Logger) *WorkspaceHandler {
	return &WorkspaceHandler{
		cfg: cfg,
		log: log,
	}
}
