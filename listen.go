package main

import (
    "fmt"
    "os"
	"net"
	//"io"
)


var cinput net.Conn
var coutput net.Conn
var started=false

func gohandleListener(l net.Listener, ptrc *net.Conn){
	for {
		if *ptrc == nil {
			conn, err := l.Accept()
			//started=true
			if err == nil {
				*ptrc = conn
				fmt.Println("got conn")
			}
		}
	}
}

func handleOut(){
	for {
		if cinput != nil && coutput != nil {
			started=true
			//io.Copy(conn,cin) //maybe handle diffrently
			tmpbuf:=make([]byte,1)
			dobreak:=false
			for {
				_, err := cinput.Read(tmpbuf)
				if err != nil {
					cinput.Close()
					cinput=nil
					fmt.Println("Input Closed")
					dobreak=true
				}
				_, err = coutput.Write(tmpbuf)
				if err != nil {
					coutput.Close()
					coutput=nil
					fmt.Println("Output Closed")
					dobreak=true
				}
				if dobreak {
					break
				}

			}
			//cin.Close()
			//return
		}
		if cinput == nil && coutput == nil && started {
			return
		}
	}
}



func main() {
	
	cinput=nil
	coutput=nil
	
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println("Can't get Current Directory",err.Error())
		return
	}

	SockAddr:=mydir + "/default.in.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	ldef, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer ldef.Close()

	SockAddr=mydir + "/out.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	lout, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer lout.Close()


	go gohandleListener(ldef,&cinput)
	go gohandleListener(lout,&coutput)
	
	handleOut()
	fmt.Println("exit")
	os.RemoveAll(mydir + "/default.in.sock")
	os.RemoveAll(mydir + "/out.sock")
}
