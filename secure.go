package main

// The self-signed certificate and the key are created using this command: openssl req -x509 -newkey rsa:4096 -nodes -keyout key.pem -out cert.pem -sha256 -days 365
// from https://github.com/dunglas/vulcain/issues/13 and from https://stackoverflow.com/questions/10175812/how-to-generate-a-self-signed-ssl-certificate-using-openssl
// or they are created on the server using lets encrypt and used here

// To open the Oracle Cloud instance ports (like port 8443) follow https://cleavr.io/cleavr-slice/opening-port-80-and-443-for-oracle-servers/ for OCI interface then
// sudo firewall-cmd --zone=public --add-port=8443/tcp --permanent  #  or --add-service=http
// sudo firewall-cmd --reload (from https://stackoverflow.com/questions/62326988/cant-access-oracle-cloud-always-free-compute-http-port)

import (
	"crypto/tls"
	"net/http"
)

func hello(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello, world\n"))
}
func main() {
	config := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
	http.HandleFunc("/", hello)
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: config,
	}
	cert := "/etc/letsencrypt/live/microcitest.info/cert.pem"
	key := "/etc/letsencrypt/live/microcitest.info/privkey.pem"
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		panic(err)
	}
}
