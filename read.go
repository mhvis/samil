package samil

import (
	"io/ioutil"
)

func (s Samil) read() {
	defer s.Close()
	for {
		msg, err:=ioutil.ReadAll(s)
		if err != nil {
			//Report error
			return
		}
	}
}
