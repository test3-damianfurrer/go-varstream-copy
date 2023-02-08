package main

import (
    "fmt"
    "os"
	"net"
	//"io"
)


var cinput net.Conn
var coutput net.Conn
var coverride net.Conn
var started=false
const S_TMPBUF=1024

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
			tmpbuf:=make([]byte,S_TMPBUF)
			dobreak:=false
			var err error
			for {
				if coverride != nil{
					_, err = cinput.Read(tmpbuf) //continue default stream read 
					if err != nil {
						cinput.Close()
						cinput=nil
						fmt.Println("Input Closed")
						dobreak=true
						err=nil
					}
					_, err = coverride.Read(tmpbuf)
					if err != nil {
						coverride.Close()
						coverride=nil
						fmt.Println("Override Closed")
						dobreak=true
						err=nil
					}
				} else {
					_, err = cinput.Read(tmpbuf)
					if err != nil {
						cinput.Close()
						cinput=nil
						fmt.Println("Input Closed")
						dobreak=true
						err=nil
					}
				}
				_, err = coutput.Write(tmpbuf)
				if err != nil {
					coutput.Close()
					coutput=nil
					fmt.Println("Output Closed")
					dobreak=true
					err=nil
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
	coverride=nil
	
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
	SockAddr=mydir + "/overlay.in.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	lovr, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer lovr.Close()


	go gohandleListener(ldef,&cinput)
	go gohandleListener(lout,&coutput)
	go gohandleListener(lovr,&coverride)
	
	
	handleOut()
	fmt.Println("exit")
	os.RemoveAll(mydir + "/default.in.sock")
	os.RemoveAll(mydir + "/overlay.in.sock")
	os.RemoveAll(mydir + "/out.sock")
}
