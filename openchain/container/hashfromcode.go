package container

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/openblockchain/obc-peer/openchain/util"
	pb "github.com/openblockchain/obc-peer/protos"
)

//name could be ChaincodeID.Name or ChaincodeID.Path
func generateHashFromSignature(path string, ctor string, args []string) []byte {
	fargs := ctor
	if args != nil {
		for _, str := range args {
			fargs = fargs + str
		}
	}
	cbytes := []byte(path + fargs)

	b := make([]byte, len(cbytes))
	copy(b, cbytes)
	hash := util.ComputeCryptoHash(b)
	return hash
}

//generateHashcode gets hashcode of the code under path. If path is a HTTP(s) url
//it downloads the code first to compute the hash.
//NOTE: for dev mode, user builds and runs chaincode manually. The name provided
//by the user is equivalent to the path. This method will treat the name
//as codebytes and compute the hash from it. ie, user cannot run the chaincode
//with the same (name, ctor, args)
func generateHashcode(spec *pb.ChaincodeSpec, path string) (string, error) {

	ctor := spec.CtorMsg
	if ctor == nil || ctor.Function == "" {
		return "", fmt.Errorf("Cannot generate hashcode from empty ctor")
	}

	hash := generateHashFromSignature(spec.ChaincodeID.Path, ctor.Function, ctor.Args)

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %s", err)
	}

	newSlice := make([]byte, len(hash)+len(buf))
	copy(newSlice[len(buf):], hash[:])
	//hash = md5.Sum(newSlice)
	hash = util.ComputeCryptoHash(newSlice)

	return hex.EncodeToString(hash[:]), nil
}
