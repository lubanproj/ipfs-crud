package ipfs

import (
	"context"
	"fmt"
	"io"

	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/repo/fsrepo"
)

// InitRepoAndGetApi get ipfs core api entry
func InitRepoAndGetApi(ctx context.Context, repoPath string, isLocal bool) (icore.CoreAPI, error) {
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
