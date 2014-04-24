package main

import (
	"flag"
	"log"
)

import (
	"github.com/qiniu/api/resumable/io"
)

func blockNotify(blkIdx int, blkSize int, ret *io.BlkputRet) {
	log.Println("block id", blkIdx, "block size", blkSize)
}

func main() {
	token := flag.String("t", "", "upload token")
	file := flag.String("f", "", "file path")
	key := flag.String("k", "", "file key")
	flag.Parse()
	if token == nil || file == nil || key == nil {
		log.Fatalln("invalid args")
		return
	}

	params := map[string]string{"x:qiniuput": "put"}
	extra := &io.PutExtra{
		ChunkSize: 1024,
		MimeType:  "text/plain",
		Notify:    blockNotify,
		Params:    params,
	}
	var ret io.PutRet
	err := io.PutFile(nil, &ret, *token, *key, *file, extra)
	if err != nil {
		log.Fatalln(err)
	}
}
