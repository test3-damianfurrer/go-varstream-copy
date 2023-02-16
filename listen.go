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
var outputbufs []*[]byte
//var tmpoutputbufs []*[]byte
var started=false
var prefix=""
var S_TMPBUF=1
var tmpbuf []byte
var readable=true //to limit read speed to fastest write

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

func goStreamWriter(c *net.Conn,ptrbuf **[]byte){
	l_buf:=make([]byte,0)
	//l_tmpbuf:=make([]byte,0)
	(*ptrbuf)=&l_buf
	for {
		//fmt.Println("writebuffer len", len(l_buf))
		if (*c != nil) && (len(l_buf)>=S_TMPBUF){
			_,err := (*c).Write(l_buf[:S_TMPBUF-1])
			readable=true
			if err != nil {
				(*c).Close()
				*c=nil
			}
			l_buf=l_buf[S_TMPBUF:]
		}
	}
}
/*func goStreamWriter(c *net.Conn, bufptr *[]byte){
	for {
		if (*c != nil) && (len(*bufptr)>S_TMPBUF){
			c.Write((*bufptr)[:S_TMPBUF])
			(*bufptr)=(*bufptr)[S_TMPBUF:]
		}
	}
}*/

func gohandleListenerMulti(l net.Listener, ptrcarr *[]net.Conn){
	for {
		conn, err := l.Accept()
		fmt.Println(prefix+"accept")
		//started=true
		if err == nil {
			fmt.Println("ptrarr len: ",len((*ptrcarr)))
			
			dobreak:=false
			var i int
			for i=0; i<len((*ptrcarr)); i++ {
				if (*ptrcarr)[i] == nil {
					(*ptrcarr)[i] = conn
					fmt.Println(prefix+"got new multi conn")
					dobreak=true
					break
				}
			}
			//dobreak=(dobreak||(len(*ptrcarr)==0))
			if !dobreak {
				outputbufs=append(outputbufs,nil)
				*ptrcarr = append(*ptrcarr,conn)
				go goStreamWriter(&(*ptrcarr)[i],&(outputbufs[i]))
				fmt.Println(prefix+"got new multi conn")
			}
			fmt.Println("after ptrarr len: ",len((*ptrcarr)))
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
		//if cinput != nil && coutput != nil {
		if cinput != nil {
			started=true
			//io.Copy(conn,cin) //maybe handle diffrently
//			tmpbuf:=make([]byte,S_TMPBUF)
			dobreak:=false
			var err error
			for {
				if readable {
					_, err = cinput.Read(tmpbuf)
					fmt.Println("after read buf")
					readable=false
					if err != nil {
						cinput.Close()
						cinput=nil
						fmt.Println(prefix+"Input Closed")
						dobreak=true
						err=nil
					}
					for i:=0; i<len(coutputs); i++ {
						fmt.Println("prc output index: ",i)
						if coutputs[i] == nil {
							continue
						}
						cbuf:=outputbufs[i]
						(*cbuf)=append((*cbuf),tmpbuf...)
						fmt.Println("add tmpbuf to output index: ",i)
						/*cout=coutputs[i]
						_, err = cout.Write(tmpbuf)
						if err != nil {
							cout.Close()
							fmt.Println(prefix+"Output Closed")
							coutputs[i]=nil
							//cinput.Close()
							//cinput=nil
							//dobreak=true
							//err=nil
						}*/
					}
					if dobreak {
						break
					}
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
	tmpbuf=make([]byte,S_TMPBUF)
	
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
