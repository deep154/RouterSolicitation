# RouterSolicitation
Router Solicitation Message ( Ipv6 support 5g)

#Ipv6 stateless address auto-configuration
1) Stateless address auto-configuration is the process that IPv6 nodes (hosts or routers) use to automatically configure IPv6 addresses for interfaces.
2) In 5G technology, SMF is given the responsibility of distributing globally unique ipv6 prefix and the link-local address (UE shall use this interface identifier provided by SMF to configure its link-local address).

Library Used :
1) go-gtp/gtpv1
2) mdlayher/ndp

Router Solicitation Message is send containing the 
1) Flags  & TunnelID
These are used as default set in library and can be configured from functions used.

In Response to Router Solicitation msg , "Router Advertisement msg" is sent over the same TunnelID containing the Ipv6 prefix for the UE.

Use : go run sol.go 
Please use Ip and ports as your specific convinince. Target IP:port are hardcoded here.
