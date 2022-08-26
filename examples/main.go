package main

import (
	"context"
	"fmt"

	crud "github.com/lubanproj/ipfs-crud"
)

func main() {
	ctx := context.Background()
	api, err := crud.InitRepoAndGetApi(ctx, true)
	if err != nil {
		fmt.Println("init repo err", err)
		return
	}
	cid, err := crud.AddFileWithContent(ctx, api, []byte("hello world"))
	if err != nil {
		fmt.Println("add file err",err)
		return
	}
	fmt.Println("add file succ, file cid", cid)
}

