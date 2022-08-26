package ipfs_crud

import (
	"context"
	"fmt"
	"github.com/ipfs/kubo/plugin/loader"
	"io"
	"os"
	"path/filepath"
	"sync"

	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/repo/fsrepo"
)

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

var loadPluginsOnce sync.Once

// InitRepoAndGetApi get ipfs core api entry
func InitRepoAndGetApi(ctx context.Context, isLocal bool) (icore.CoreAPI, error) {
	var onceErr error
	loadPluginsOnce.Do(func() {
		onceErr = setupPlugins("")
	})
	if onceErr != nil {
		return nil, onceErr
	}

	cfg, err := config.Init(io.Discard, 2048)
	if err != nil {
		return nil, err
	}

	if isLocal {
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#ipfs-filestore
		cfg.Experimental.FilestoreEnabled = true
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#ipfs-urlstore
		cfg.Experimental.UrlstoreEnabled = true
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#ipfs-p2p
		cfg.Experimental.Libp2pStreamMounting = true
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#p2p-http-proxy
		cfg.Experimental.P2pHttpProxy = true
		// See also: https://github.com/ipfs/kubo/blob/master/docs/config.md
		// And: https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md
	}

	// Create the repo with the config
	repoPath, err := os.MkdirTemp("", "ipfs-shell")
	if err != nil {
		return nil, fmt.Errorf("failed to get temp dir: %s", err)
	}

	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init ephemeral node: %s", err.Error())
	}

	node, err := createNode(ctx, repoPath)
	if err != nil {
		return nil, err
	}
	api, err := coreapi.NewCoreAPI(node)
	return api, err
}


func createNode(ctx context.Context, filepath string) (*core.IpfsNode, error) {
	repo, err := fsrepo.Open(filepath)
	if err != nil {
		return nil, err
	}
	nodeOptions := &core.BuildCfg{
		Online: true,
		Routing: libp2p.DHTOption,
		Repo: repo,
	}
	return core.NewNode(ctx, nodeOptions)
}
