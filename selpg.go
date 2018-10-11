package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
	flag "github.com/spf13/pflag"
)

var progName string
var sa selpg_args

type selpg_args struct {
	startPage  int
	endPage    int
	inFile string
	pageLen    int
	pageType   bool
	printDest string
}

func Init() {
	flag.IntVarP(&(sa.startPage), "startPage", "s", -1, "-s int or -startPage int")
	flag.IntVarP(&(sa.endPage), "endPage", "e", -1, "-e int or -endPage int")
	flag.IntVarP(&(sa.pageLen), "length", "l", 72, "-length int or -l int")
	flag.BoolVarP(&(sa.pageType), "formfeed", "f", false, "-formfeed bool or -f bool")
	flag.StringVarP(&(sa.printDest), "dest", "d", "", "-dest string or -d string")
}

func main() {
	Init()
	flag.Parse()
	process_args()
	process_input()
}

/
func process_args() {
	if flag.NArg() > 1 {
		panic_and_print_usage()
	}
	if flag.NArg() == 1 {
		sa.inFile = flag.Args()[0]
	}
	if sa.startPage == -1 || sa.endPage == -1 {
		fmt.Fprint(os.Stderr, "you must set your start page and end page with positive integers")
		panic_and_print_usage()
	}
	if sa.startPage > sa.endPage {
		fmt.Println("start large than end page")
		os.Exit(1)
	}
}
func usage() string {
	return progName + " -sNumber -eNumber [-lNumber] [-f] [-dDestination] [output file name]"

}

func process_input() {
	var file *os.File
	var cmd *exec.Cmd
	var cmdin io.WriteCloser
	defer file.Close()
	if sa.inFile == "" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(sa.inFile)
		if err == nil {
			//fmt.Println("file open successfully")
		} else {
			fmt.Println(err)
		}
	}
	reader := bufio.NewReader(file)
	var writer *bufio.Writer
	if sa.printDest == "" {
		writer = bufio.NewWriter(os.Stdout)
	} else {
		var err error
		cmd = exec.Command("lp", "-d"+sa.printDest)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println("failed open cmd ", cmd.Args, "stdinpipe")
		}
	}
	chanb := make(chan byte)
	chanel := make(chan []byte, 10)
	quit := make(chan int)
	defer close(chanb)
	defer close(chanel)
	defer close(quit)

	go func() {

		if sa.pageType == false {
			pgctr := 1
			lctr := 0
			for {
				line, _, err := reader.ReadLine()
				if err == nil {
					lctr++
					if lctr >= sa.pageLen {
						pgctr++
						lctr = 0
					}
					if pgctr < sa.startPage {
						continue
					}
					if pgctr > sa.endPage {
						chanel <- line
						time.Sleep(time.Millisecond * 100)
						quit <- 0
						break
					}
					chanel <- line
				} else if err == io.EOF {
					time.Sleep(time.Millisecond * 100)
					quit <- 0
					break
				} else {
					time.Sleep(100 * time.Millisecond)
					quit <- 0
					break
				}

			}
		} else {
			pgctr := 1
			for {
				b, err := reader.ReadByte()
				if err == nil {
					if b == '\f' {
						pgctr++
					}
					if pgctr < sa.startPage {
						continue
					} else if pgctr > sa.endPage {
						time.Sleep(100 * time.Millisecond)
						quit <- 0
						break
					}
					//work
					chanb <- b
				} else if err == io.EOF {
					time.Sleep(100 * time.Millisecond)
					quit <- 0
					break
				} else {
					time.Sleep(time.Millisecond * 100)
					quit <- 0
					break
				}
			}
		}
	}()

	func() {
		for {
			select {
			case line := <-chanel:
				if sa.printDest == "" {
					writer.Write(line)
					writer.WriteByte('\n')
					writer.Flush()
				} else {
					fmt.Fprint(cmdin, (string)(line)+"\n")
				}
			case <-quit:
				if sa.printDest != "" {
					cmd.Start()
					cmd.Wait()
				}
				break
			case bt := <-chanb:
				if sa.printDest == "" {
					writer.WriteByte(bt)
					writer.Flush()
				} else {
					fmt.Fprint(cmdin, string(bt))
				}
			}
		}

	}()
}
func panic_and_print_usage() {
	panic(errors.New(usage()))
}
