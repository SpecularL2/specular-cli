package workspace

import (
	"context"
	"fmt"
	"github.com/SpecularL2/specular-cli/internal/service/config"
	"github.com/hashicorp/go-getter"
	"os"
	"os/user"
)

const defaultWorkspaceUrl = "github.com/SpecularL2/specular//config/local_devnet"

type WorkspaceHandler struct {
	cfg *config.Config
}

func (w *WorkspaceHandler) DownloadDefault() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	var dst = fmt.Sprintf("%s/.spc/workspaces/default", usr.HomeDir)

	client := &getter.Client{
		Ctx: context.Background(),
		// Define the destination to where the directory will be stored.
		// This will create the directory if it doesn't exist.
		Dst: dst,
		Dir: true,
		// The repository with a subdirectory I would like to clone only.
		Src:  defaultWorkspaceUrl,
		Mode: getter.ClientModeDir,
		// Define the type of detectors go getter should use, in this case only GitHub is needed
		Detectors: []getter.Detector{
			&getter.GitHubDetector{},
		},
		// Provide the getter needed to download the files
		Getters: map[string]getter.Getter{
			"git": &getter.GitGetter{},
		},
	}
	// Download the files
	if err := client.Get(); err != nil {
		fmt.Printf("Error getting path %s: %v", client.Src, err)
		os.Exit(1)
	}
	return nil
}

func (w *WorkspaceHandler) Cmd() error {
	if w.cfg.Workspace.Command == "download" && w.cfg.Workspace.Name == "default" {
		fmt.Printf("Downloading default workspace from: %s\n", defaultWorkspaceUrl)
		if err := w.DownloadDefault(); err != nil {
			return err
		}
	}
	return nil
}

func NewWorkspaceHandler(cfg *config.Config) *WorkspaceHandler {
	return &WorkspaceHandler{
		cfg: cfg,
	}
}
