package INIT

import (
	"encoding/json"
	"os"

	"github.com/Turk1shGuy/torrent/internal/global"
)

func readConf(conf_path string) global.C0nf {
	file, err := os.ReadFile(conf_path)
	if err != nil {
		panic(err.Error())
	}

	var data global.C0nf
	err = json.Unmarshal(file, &data)
	if err != nil {
		panic(err.Error())
	}

	return data
}
