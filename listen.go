package main

import (
    "fmt"
    "os"
	"net"
//	"io"
)

func defStream(c net.Conn,out net.Conn,overlconn *net.Conn) {
 //while overlconn == nil -> copy c to out
}

func echoServer(c net.Conn) {
    fmt.Printf("Client connected [%s]\n", c.RemoteAddr().Network())
    fmt.Println("addr",c.RemoteAddr())
    //io.Copy(c, c)
	for{
		databuf := make([]byte,0)
		tmpbuf := make([]byte, 1)
		
		for {
			_, err := c.Read(tmpbuf)
			if err != nil {
				if err.Error() != "EOF"{
					fmt.Println("READ ERR",err.Error())
				} else {
					fmt.Println("Client Connection Closed",err.Error())
				}
				c.Close()
				return
					
			}
			//fmt.Println("byte",tmpbuf[0])
			if tmpbuf[0] == '\n' {
				databuf = append(databuf,'\n')
				break
			}
			if tmpbuf[0] == 0 {
				databuf = append(databuf,'\n')
				break
			}
			if tmpbuf[0] == 10 {
				databuf = append(databuf,'\n')
				break
			}
			databuf = append(databuf,tmpbuf[0])
		}
		fmt.Printf("Received: %s",databuf)
		c.Write([]byte{'Y','o','u',' ','s','e','n','t',':',' '})
		c.Write(databuf)
	}
    c.Close()
    fmt.Println("Connection Closed")
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

    SockAddr:=mydir + "/overlay.in.sock"
    if err := os.RemoveAll(SockAddr); err != nil {
        panic(err)
    }
    lovr, err := net.Listen("unix", SockAddr)
    if err != nil {
        fmt.Println("listen error:",err.Error())
    }
    defer lovr.Close()

    SockAddr:=mydir + "/out.sock"
    if err := os.RemoveAll(SockAddr); err != nil {
        panic(err)
    }
    dout, err := net.Dial("unix", SockAddr)
    if err != nil {
        fmt.Println("failed to create output socket,  error:",err.Error())
	return
    }
    defer dout.Close()

    conn, err := ldef.Accept()
    if err != nil {
       fmt.Println("DEFAULT IN: accept error:", err.Error())
	return
    }
    defconn:=conn
    conn=nil
    go defStream(defconn, dout, &conn)

    for {
        conn, err = lovr.Accept()
        if err != nil {
            fmt.Println("accept error:", err.Error())
        }
	//while conn alive copy data
	conn.Close()
	conn=nil
        //go echoServer(conn)
    }
}
