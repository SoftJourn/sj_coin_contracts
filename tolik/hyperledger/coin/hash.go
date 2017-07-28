/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * Helper tool to hash signing identity (in development)
 * author: aokhotnikov@softjourn.com
 */

/*
$ go build hash.go
$ ./hash /tmp/fabric-client-kvs_peerOrg1/Jim
*/

package main
import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"

    //"github.com/golang/protobuf/proto"
    //"github.com/hyperledger/fabric/bccsp"
    //"github.com/hyperledger/fabric/bccsp/factory"
    "github.com/hyperledger/fabric/msp"
    //mspproto "github.com/hyperledger/fabric/protos/msp"    
)
type DeserializersManager struct {
        LocalDeserializer    msp.IdentityDeserializer
        ChannelDeserializers map[string]msp.IdentityDeserializer
}
func hashCreator(Creatorbytes []byte) ([]byte, error) {
    var Creatorhash []byte
    var err error
    if Creatorbytes == nil {
        return nil, fmt.Errorf("Creator is not found")
    }
    //Creatorhash, err = factory.GetDefault().Hash(Creatorbytes, &bccsp.SHA256Opts{})
    if err != nil {
        return nil, fmt.Errorf(fmt.Sprintf("Failed computing SHA256 on [% x]", Creatorbytes))
    }
    return Creatorhash, err
}
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

    raw, err := msp.NewSerializedIdentity(mspid, []byte{0x00,0x00,0x00})
/*
    sId := &mspproto.SerializedIdentity{Mspid: mspID, IdBytes: []byte(certificate)}
    raw, err := proto.Marshal(sId)
*/
    if err != nil {
        fmt.Printf("Error cannot marshal identity: %s\n\n", err)
        return
    }
    fmt.Printf("SerializedIdentity: \n%x\n\n", raw)

    Creatorhash, err := hashCreator(raw)
    if err != nil {
      fmt.Printf("Failed to hash creator: %s\n\n", err)
      return
    }
    fmt.Printf("Creator hash: %s\n\n", fmt.Sprintf("%x", Creatorhash))
}
