package sol

import (
	"log"
	"os/exec"
)

func (b *Baize) Wikipedia() {
	cmd := exec.Command("open", b.script.Wikipedia())
	if cmd == nil {
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
}
