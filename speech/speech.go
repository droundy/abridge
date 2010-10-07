package speech

import (
	"fmt"
	"os"
	"exec"
	"io/ioutil"
	"sync"
	"github.com/droundy/bridge"
)

type answertype struct {
	sound string
	err os.Error
}
type speechrequest struct {
	bid string
	answer chan<- answertype
}

var seek_speech = make(chan speechrequest)

var once sync.Once

func Speak(bid string) (string, os.Error) {
	once.Do(func() {
		cache := make(map[string]string)
		translations := make(map[string]string)
		for i:='1'; i<'7'; i++ {
			for sv:=0;sv<=bridge.NoTrump;sv++ {
				stringi := string([]byte{byte(i)})
				spoken := stringi + " " + bridge.SuitName[sv] + "."
				translations[stringi+bridge.SuitHTML[sv]] = spoken
				translations[stringi+bridge.SuitLetter[sv]] = spoken
				translations[stringi+bridge.SuitColorHTML[sv]] = spoken
			}
		}
		translations[" P"] = "Pass."
		translations["P"] = "Pass."
		translations[" X"] = "Double!"
		translations["X"] = "Double!"
		translations["XX"] = "Redouble!"
		espeak, err := exec.LookPath("espeak")
		if err == nil {
			go func() {
				for {
					req := <- seek_speech
					bid := req.bid
					if sp,ok := translations[bid]; ok {
						bid = sp
					}
					fmt.Println("I am speaking:", bid)
					if s,ok := cache[req.bid]; ok {
						req.answer <- answertype{ s, nil }
					} else {
						c,err := exec.Run(espeak, []string{"espeak","--stdout"}, nil, "",
							exec.Pipe, exec.Pipe, exec.PassThrough)
						if err != nil {
							req.answer <- answertype{ "", err }
							continue
						}
						defer c.Close()
						fmt.Fprintln(c.Stdin, bid)
						c.Stdin.Close()
						o,err := ioutil.ReadAll(c.Stdout)
						if err != nil {
							req.answer <- answertype{ "", err }
							continue
						}
						out := string(o)
						cache[req.bid] = out
						req.answer <- answertype{ out, nil }
					}
				}
			}()
		} else {
			go func() {
				fmt.Println("We won't be able to speak, since I can't find espeak!", err)
				for {
					req := <- seek_speech
					req.answer <- answertype { "", os.NewError("Cannot find espeak in path.") }
				}
			}()
		}
	})
	ans := make(chan answertype)
	seek_speech <- speechrequest{ bid, ans }
	x := <- ans
	return x.sound, x.err
}
