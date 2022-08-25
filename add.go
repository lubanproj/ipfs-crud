package ipfs

import (
	"context"
	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"os"
)

// AddFileWithContent add an unexist file with file content
func AddFileWithContent(ctx context.Context, api icore.CoreAPI, content []byte) (string, error) {
	fileCid, err := api.Unixfs().Add(ctx, files.NewBytesFile(content))
	if err != nil {
		return "",  err
	}
	err = api.Pin().Add(ctx, fileCid)
	if err != nil {
		return "",  err
	}
	return fileCid.String(), nil
}

// AddFileWithPath add an exist file with a given file path
func AddFileWithPath(ctx context.Context, api icore.CoreAPI, filepath string) (string, error) {
	fileNode, err := getUnixfsNode(filepath)
	if err != nil {
		return "", err
	}
	fileCid, err := api.Unixfs().Add(ctx, fileNode)
	if err != nil {
		return "",  err
	}
	err = api.Pin().Add(ctx, fileCid)
	if err != nil {
		return "",  err
	}
	return fileCid.String(), nil
}


func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}



