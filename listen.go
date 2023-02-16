package main

import (
    "fmt"
    "os"
	"net"
	//"io"
)


var cinput net.Conn
var cout net.Conn
var coutputs []net.Conn

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
func gohandleListenerMulti(l net.Listener, ptrcarr *[]net.Conn){
	for {
		conn, err := l.Accept()
		//started=true
		if err == nil {
			dobreak:=false
			for i:=0; i<len(*ptrcarr); i++ {
				if *ptrcarr[i] == nil {
					*ptrcarr[i] = conn
					dobreak=true
					break
				}
			}
			if !dobreak {
				*ptrcarr = append(*ptrcarr,conn)
			}
			fmt.Println(prefix+"got new multi conn")
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
				_, err = cinput.Read(tmpbuf)
				if err != nil {
					cinput.Close()
					cinput=nil
					fmt.Println(prefix+"Input Closed")
					dobreak=true
					err=nil
				}
				for i:=0; i<len(coutputs); i++ {
					cout=coutputs[i]
					_, err = cout.Write(tmpbuf)
					if err != nil {
						cout.Close()
						fmt.Println(prefix+"Output Closed")
						coutputs[i]=nil
						//cinput.Close()
						//cinput=nil
						//dobreak=true
						//err=nil
					}
				}
				if dobreak {
					break
				}

			}
			//cin.Close()
			//return
		}
		//if cinput == nil && coutput == nil && started {
		if cinput == nil && started {
			for i:=0; i<len(coutputs); i++ {
				if coutputs[i] != nil {
					coutputs[i].Close()
				}
			}
			return
		}
	}
}



func main() {
	
	cinput=nil
	coutputs=make([]net.Conn,0)
	
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

	SockAddr:=mydir + "/"+prefix+"in2m.in.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	ldef, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer ldef.Close()

	SockAddr=mydir + "/"+prefix+"in2m.out.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		panic(err)
	}
	lout, err := net.Listen("unix", SockAddr)
	if err != nil {
		fmt.Println("listen error:",err.Error())
	}
	defer lout.Close()
	
	go gohandleListener(ldef,&cinput)
	go gohandleListenerMulti(lout,&coutputs)
	
	handleOut()
	fmt.Println("exit")
	os.RemoveAll(mydir + "/"+prefix+"in2m.in.sock")
	os.RemoveAll(mydir + "/"+prefix+"in2m.out.sock")
}
