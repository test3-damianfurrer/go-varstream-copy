package main

import (
    "fmt"
    "os"
	"net"
//	"io"
)

const S_TMPBUF=1024

var outs []net.Conn //:=make(net.Conn,0) //maybe later handle with a stream copy socket, e.g. every (in/listen) output conn just (go routined) stream copies from a real output socket
//"real out" = dial it, an out connection, instead of listener
//a real out socket would be created from another program
//we would need to pass this info to our program (maybe another socket :-))
//we could also create a real out from and for this program.
//go routines -> one connection to another go routine that distributes the data to every connection. 
//this way we could just pass the out conn to the write functions. (instead of working with an array)

func handleOut(lout net.Listener){
	for {
		conn, err := lout.Accept()
		outs=append(outs,conn)
		if err != nil {
			fmt.Println("accept error:", err.Error())
		}
		//while conn alive copy data
		//nn.Close()
		//nn=nil
		//go echoServer(conn)
	}
}

func defStream(cdef net.Conn,overlconn *net.Conn) {
	for {
		if *overlconn==nil {
			tmpbuf := make([]byte, S_TMPBUF)
			_, err := cdef.Read(tmpbuf) //n
			if err != nil {
				return //prob. req. program end
			}
			for i:=0;i<len(outs);i++ {
				_ ,err = outs[i].Write(tmpbuf)
				if err != nil {
					fmt.Println("Conn write err",err)	
				}
			}
		}	
	}
 //while overlconn == nil -> copy c to out
}

func ovrlStream(covr net.Conn) {
	for {
		tmpbuf := make([]byte, S_TMPBUF)
		_, err := covr.Read(tmpbuf) //n
		if err != nil {
			covr.Close()
			covr=nil
			return
		}
		for i:=0;i<len(outs);i++ {
			_ ,err = outs[i].Write(tmpbuf)
			if err != nil {
				fmt.Println("Conn write err",err)	
			}
		}	
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

	SockAddr=mydir + "/overlay.in.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	lovr, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer lovr.Close()

	/*
	SockAddr=mydir + "/out.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	dout, err := net.Dial("unix", SockAddr)
	if err != nil {
		fmt.Println("failed to create output socket,  error:",err.Error())
		return
	}
	defer dout.Close()
	*/
	//out is a listen socket too, it just doesn't listen, but write instead.
	SockAddr=mydir + "/out.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	lout, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer lout.Close()

	conn, err := ldef.Accept()
	if err != nil {
		fmt.Println("DEFAULT IN: accept error:", err.Error())
		return
	}
	
	go handleOut(lout)
	
	defconn:=conn
	conn=nil
	go defStream(defconn, &conn)

	for {
		var conn2 net.Conn
		conn2, err = lovr.Accept()
		if err != nil {
			fmt.Println("accept error:", err.Error())
		}
		if conn != nil {
			conn.Close()
		}
		conn = conn2
		go ovrlStream(conn)
		//while conn alive copy data
		
		//go echoServer(conn)
	}
}
