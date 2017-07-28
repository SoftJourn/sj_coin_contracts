package main
import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"

    "encoding/base64"
    "encoding/pem"
    "crypto/x509"
    "crypto/ecdsa"
    //"github.com/golang/protobuf/proto"
    //"github.com/hyperledger/fabric/msp"
    //mspproto "github.com/hyperledger/fabric/protos/msp"    
)
/*
./msp/mspimpl.go:
func (msp *bccspmsp) getCertFromPem(idBytes []byte) (*x509.Certificate, error) {
func (msp *bccspmsp) DeserializeIdentity(serializedID []byte) (Identity, error) {
*/

func main() {
	fmt.Println("Hash helper start...")
    path := os.Args[1]

    certPEM, err := ioutil.ReadFile(path)
    if err != nil {
         fmt.Printf("Cannot read file: %s\n\n", err)
         return
    }
    //fmt.Printf("certPEM = %s", certPEM)

    var dat map[string]interface{}
    if err := json.Unmarshal(certPEM, &dat); err != nil {
        panic(err)
    }
    //fmt.Println(dat)
    
    mspid := dat["mspid"].(string)
    fmt.Printf("mspid=%s\n", mspid)

    enrollment := dat["enrollment"].(map[string]interface {})
    //fmt.Println(enrollment)

    signingIdentity := enrollment["signingIdentity"].(string)
    fmt.Printf("signingIdentity=%s\n", signingIdentity)
    
    identity := enrollment["identity"].(map[string]interface {})
    certificate := identity["certificate"].(string)
    fmt.Println(certificate)

    block, _ := pem.Decode([]byte(certificate))
    var cert* x509.Certificate
    cert, _ = x509.ParseCertificate(block.Bytes)
    ecdsaPublicKey := cert.PublicKey.(*ecdsa.PublicKey)
    fmt.Println(ecdsaPublicKey.X)
    fmt.Println(ecdsaPublicKey.Y) 

    raw, err := x509.MarshalPKIXPublicKey(ecdsaPublicKey)
    if err !=nil {
    	panic("failed to marshal public key")
    }

    fmt.Printf("\n\n%x\n\n", raw)
    encoded := base64.StdEncoding.EncodeToString(raw)
    fmt.Println(encoded)
    fmt.Println("\n\n")

	var pubPEMData = []byte(`
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEBPipWZtkmd8mBWqzrqVhDJgKVbQE
P2ucIy9tVhQncM6hh7BpBkaEDZ3sYV5I3WZ2K1e5aDXajjut1SpSaCS8YQ==
-----END PUBLIC KEY-----
and some more`)

	block, rest := pem.Decode(pubPEMData)
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

    ecdsaPublicKey = pub.(*ecdsa.PublicKey)
    fmt.Println(ecdsaPublicKey.X)
    fmt.Println(ecdsaPublicKey.Y)

	fmt.Printf("Got a %T, with remaining data: %q\n\n", pub, rest)

}
