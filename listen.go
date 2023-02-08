package main

import (
    "fmt"
    "os"
	"net"
	"io"
)

func handleOut(lout net.Listener,cin net.Conn){
	for {
		conn, err := lout.Accept()
		fmt.Println("got out conn")
		if err != nil {
			panic(err)
		}
		io.Copy(conn,cin) //maybe handle diffrently
		//cin.Close()
		return
	}
}



func main() {
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

	
	for{
	fmt.Println("Expecting default Socket to get stream")
	conn, err := ldef.Accept()
	if err != nil {
		fmt.Println("DEFAULT IN: accept error:", err.Error())
		return
	}
	//go
	fmt.Println("handle out")
	handleOut(lout,conn)
	}
	fmt.Println("exit")
	os.RemoveAll(mydir + "/default.in.sock")
	os.RemoveAll(mydir + "/out.sock")
}
