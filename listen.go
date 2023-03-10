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
var prefix=""
var S_TMPBUF=1
var tmpbuf []byte

func nilnetconn(cptr *net.Conn){
	if *cptr != nil {
		fmt.Println("manual nil")
		*cptr=nil
		fmt.Println("manual nil aftr")
	}
}

func gohandleListener(l net.Listener, ptrc *net.Conn){
	for {
		if *ptrc == nil {
			conn, err := l.Accept()
			//started=true
			if err == nil {
				*ptrc = conn
				fmt.Println(prefix+"got conn")
			}
		}
	}
}

func gohandleReplaceListener(l net.Listener, ptrc *net.Conn){
	for {
		conn, err := l.Accept()
		if err == nil {
			if *ptrc != nil {
				(*ptrc).Close()
			}
			*ptrc = conn
			fmt.Println(prefix+"got new conn")
		}
	}
}

func handleOut(){
	for {
		if cinput != nil && coutput != nil {
			started=true
			//io.Copy(conn,cin) //maybe handle diffrently
//			tmpbuf:=make([]byte,S_TMPBUF)
			dobreak:=false
			var err error
			for {
				if coverride != nil{
					_, err = cinput.Read(tmpbuf) //continue default stream read 
					if err != nil {
						cinput.Close()
						nilnetconn(&cinput)
						fmt.Println(prefix+"Input Closed")
						coutput.Close()
						nilnetconn(&coutput)
						fmt.Println(prefix+"Output Closed")
						if coverride != nil {
							coverride.Close()
							nilnetconn(&coverride)
						}
						dobreak=true
						err=nil
					}
					_, err = coverride.Read(tmpbuf)
					if err != nil {
						coverride.Close()
						nilnetconn(&coverride)
						fmt.Println(prefix+"Override Closed")
						dobreak=true
						err=nil
					}
				} else {
					_, err = cinput.Read(tmpbuf)
					if err != nil {
						cinput.Close()
						nilnetconn(&cinput)
						fmt.Println(prefix+"Input Closed")
						dobreak=true
						err=nil
					}
				}
				_, err = coutput.Write(tmpbuf)
				if err != nil {
					coutput.Close()
					//coutput=nil //errors
					nilnetconn(&coutput) //test
					fmt.Println(prefix+"Output Closed")
					if coverride != nil {
						coverride.Close()
						nilnetconn(&coverride)
					}
					cinput.Close()
					nilnetconn(&cinput)
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
	
	
	if len(os.Args) >= 2 {
		if os.Args[1] != "" {
			prefix = os.Args[1]+"-"
			fmt.Println("custom prefix:",prefix)
		}
		
		if len(os.Args) >= 3 { //optional var buffer size
			_, err := fmt.Sscanf(os.Args[2],"%d",&S_TMPBUF)
			if err != nil {
				panic(err)
			}
		}
	}
	tmpbuf=make([]byte,S_TMPBUF) //only alloc once
	
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println("Can't get Current Directory",err.Error())
		return
	}

	SockAddr:=mydir + "/"+prefix+"default.in.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	ldef, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer ldef.Close()

	SockAddr=mydir + "/"+prefix+"out.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	lout, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer lout.Close()
	SockAddr=mydir + "/"+prefix+"overlay.in.sock"
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
	go gohandleReplaceListener(lovr,&coverride)
	
	
	handleOut()
	fmt.Println("exit")
	os.RemoveAll(mydir + "/"+prefix+"default.in.sock")
	os.RemoveAll(mydir + "/"+prefix+"overlay.in.sock")
	os.RemoveAll(mydir + "/"+prefix+"out.sock")
}
