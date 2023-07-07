package main

import (
        "go-gtp/gtpv1"
        "fmt"
        "github.com/mdlayher/ndp"
        "net"
	"time"
)

func main() {

        var m = &ndp.RouterSolicitation{}
        fmt.Println(m.Type())

        //marshal the ndp msg
        b, err := ndp.MarshalMessage(m)
        if err != nil {
                fmt.Println("Error")
        }

	//ip & port of this serving gtpv1 server
        laddr, err := net.ResolveUDPAddr("udp", IP + ":" + Port)
        UPlaneConn := gtpv1.NewUPlaneConn(laddr)

	// Send to Any IP & Port  
        srvAddr, err := net.ResolveUDPAddr("udp", "192.168.123.38:2152")

	time.Sleep(2*time.Second)
        UPlaneConn.WriteTo(b , srvAddr)

}

