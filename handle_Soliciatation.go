package main

import(
      "fmt"
      "net"
      "strconv"
      "context"
      "go-gtp/gtpv1"
      "github.com/mdlayher/ndp"
      "go-gtp/gtpv1/message"
      "time"
)

var Conn *gtpv1.UPlaneConn
var SenderAddr net.Addr

func main() {

    StartGtpv1Server()
}

func StartGtpv1Server() {

	fmt.Println("Starting Gtpv1 Server")

	var ipAddr *net.IPAddr
        var err error

	ipAddr, err = net.ResolveIPAddr("ip", "192.168.123.18")
        fmt.Println("Gtpv1 interface IP address: ",ipAddr)
        if err != nil {
                    panic(err)
        }

        port := 2155

        ctx := context.Background()
        laddr, err := net.ResolveUDPAddr("udp", ipAddr.String()+":"+strconv.Itoa(port))
        if err != nil {
                fmt.Println("Resolve UDP Addr failed", err)
                return
        }

	UPlaneConn := gtpv1.NewUPlaneConn(laddr)
        fmt.Println("Created Connection  %#v", UPlaneConn)

        go func(){
                if err := UPlaneConn.ListenAndServe(ctx); err != nil {
                        fmt.Println("Failed to ListenPacket", err)
                        return
                }
        }()

        fmt.Println("Started serving on GTPV1 %s", laddr)
        UPlaneConn.AddHandler(message.MsgTypeTPDU,handleTPDU)

	Conn = UPlaneConn
        pgwcIp ,err := net.LookupIP("192.168.6.38")
        fmt.Println(" PGWC IP %#v", ipAddr)
	fmt.Println(pgwcIp)
        if err != nil{
                fmt.Println("Not able to resolve SPGWC IP address from its fqdn : 192.168.6.38 ")
                return
        }
	addr, err := net.ResolveUDPAddr("udp", "192.168.6.38:3386")
        if err!=nil{
                fmt.Println("Resolve UDP Addr for SPGWC failed", err)
                return
        }
	SenderAddr = addr
        time.Sleep(1*time.Second)

}

func handleTPDU(c gtpv1.Conn , senderAddr net.Addr , msg message.Message) error {
        fmt.Println("Handling Router solication message")
        fmt.Println("Received message : %+v ", msg);
        tpdu := msg.(*message.TPDU)
        header := msg.(*message.TPDU).Header
        tpdu.Version();
        teid := header.TEID
        if(teid != 0){
                payload  := header.Payload
		fmt.Println("Payload:",payload)
                 //Unmarshal the payload into RouterSolicitation
                data , err := ndp.ParseMessage(payload)
                fmt.Println("Data:",data)
                if err != nil {
                        // Handle the error
                        fmt.Println("Error unmarshaling payload:", err)
                } else {
                        SendGtpRouterAdvertisement(teid)
                }
        }
                return nil
}

func SendGtpRouterAdvertisement(teid uint32) ( error) {

        var macAddr net.HardwareAddr
        interfaces, err := net.Interfaces()
        if err != nil {
                fmt.Println("Error:", err)
                return err
        }


        for _, iface := range interfaces {
                if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
                        macAddr = iface.HardwareAddr
                        fmt.Println("MAC Address:  %s \n", macAddr )
                        break
                }
        }

        if err != nil {
                fmt.Println("Error parsing MAC address:", err)
                return err
        }

        var ra ndp.Message

        ra, err  = CreateRouterAdvertisement(macAddr)

        if err != nil {
                fmt.Println(" Failed to create router Advertisement, error:[%v]",  err)
                return err
        }
        buff, err := ndp.MarshalMessage(ra)
        if err != nil {
                fmt.Println("failed to marshal message: %v", err)
        }

        if _, err := Conn.WriteToGTP(teid , buff , SenderAddr); err != nil {
                fmt.Println("Write failed %d", teid)
        }

        return nil
}

func CreateRouterAdvertisement (addr net.HardwareAddr) (*ndp.RouterAdvertisement, error) {

        var Message = &ndp.RouterAdvertisement{
                CurrentHopLimit:           10,
                ManagedConfiguration:      false,
                OtherConfiguration:        false,
                MobileIPv6HomeAgent:       false,
                RouterSelectionPreference: ndp.Medium,
                NeighborDiscoveryProxy:    false,
                RouterLifetime:            65535 * time.Second,
                ReachableTime:             0 * time.Millisecond,
                RetransmitTimer:           0 * time.Millisecond,
                Options: []ndp.Option{
                        &ndp.LinkLayerAddress{
                                Direction: ndp.Target,
                                Addr:      addr,
                        }},
                }

                return Message , nil
        }

