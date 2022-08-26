package ipfs_crud

import (
	"context"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

// DelFile calls pin rm, you should call gc to remove the files actually
func DelFile(ctx context.Context, api icore.CoreAPI, filepath string) error {
	err := api.Pin().Rm(ctx, path.New(filepath))
	return err
}
