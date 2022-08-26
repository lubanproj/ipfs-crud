package ipfs_crud

import (
	"context"

	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

// GetFile returns a files.Node by a given filepath
func GetFile(ctx context.Context, api icore.CoreAPI, filepath string) (files.Node, error) {
	fileNode, err := api.Unixfs().Get(ctx, path.New(filepath))
	if err != nil {
		return nil,  err
	}
	return fileNode, nil
}

func GetFileByPath(ctx context.Context, api icore.CoreAPI, filepath string, outputPath string) error {
	fileNode, err := api.Unixfs().Get(ctx, path.New(filepath))
	if err != nil {
		return err
	}
	err = files.WriteTo(fileNode, outputPath)
	return err
}
