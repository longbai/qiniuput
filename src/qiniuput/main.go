package main

import (
	"flag"
	"log"
	"net/url"
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
	custom := flag.String("x", "", "custom args")
	flag.Parse()
	if token == nil || file == nil || key == nil {
		log.Fatalln("invalid args")
		return
	}

	params := map[string]string{}
	extra := &io.PutExtra{
		ChunkSize: 1024,
		Notify:    blockNotify,
		Params:    params,
	}
	if custom != nil && *custom != "" {
		values, err := url.ParseQuery(*custom)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		for k, v := range values {
			params["x:"+k] = v[0]
		}
		extra.Params = params
	}
	var ret io.PutRet
	err := io.PutFile(nil, &ret, *token, *key, *file, extra)
	if err != nil {
		log.Fatalln(err)
	}
}
