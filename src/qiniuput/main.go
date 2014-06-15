package main

import (
	"flag"
	"log"
	"net/url"
	"os"
)

import (
	"github.com/qiniu/api/io"
	rio "github.com/qiniu/api/resumable/io"
)

func resumblePut(token, key, file string, params map[string]string) (err error) {
	blockNotify := func(blkIdx int, blkSize int, ret *rio.BlkputRet) {
		log.Println("block id", blkIdx, "offset", ret.Offset)
	}
	extra := &rio.PutExtra{
		ChunkSize: 8192,
		Notify:    blockNotify,
	}
	if len(params) != 0 {
		extra.Params = params
	}

	var ret rio.PutRet
	return rio.PutFile(nil, &ret, token, key, file, extra)

}

func put(token, key, file string, params map[string]string) (err error) {
	var extra *io.PutExtra
	if len(params) != 0 {
		extra = &io.PutExtra{
			Params: params,
		}
	}

	var ret io.PutRet
	return io.PutFile(nil, &ret, token, key, file, extra)
}

func main() {
	token := flag.String("t", "", "upload token")
	file := flag.String("f", "", "file path")
	key := flag.String("k", "", "file key")
	resume := flag.Bool("c", false, "断点续传")
	custom := flag.String("x", "", "custom args foo=1&bar=2")
	flag.Parse()
	if *token == "" || *file == "" || *key == "" {
		flag.PrintDefaults()
		log.Fatalln("invalid args")
		return
	}

	f, err := os.Open(*file)
	if err != nil {
		log.Fatalln("file not exist")
		return
	}
	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		log.Fatalln("invalid file")
		return
	}
	log.Println("file size", stat.Size())

	params := map[string]string{}

	if custom != nil && *custom != "" {
		values, err := url.ParseQuery(*custom)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		for k, v := range values {
			params["x:"+k] = v[0]
		}
	}

	if *resume {
		err = resumblePut(*token, *key, *file, params)
	} else {
		err = put(*token, *key, *file, params)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
