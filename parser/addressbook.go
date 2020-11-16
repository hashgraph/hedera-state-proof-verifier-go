package parser

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	hederaproto "github.com/hashgraph/hedera-sdk-go/proto"
	"github.com/limechain/hedera-state-proof-verifier-go/entityid"
)

func ParseAddressBooks(addressBooks []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, ab := range addressBooks {
		bytes, err := base64.StdEncoding.DecodeString(ab)
		if err != nil {
			return nil, err
		}
		var ab hederaproto.NodeAddressBook
		err = proto.Unmarshal(bytes, &ab)
		if err != nil {
			return nil, err
		}

		for _, nodeAddress := range ab.NodeAddress {
			var nodeId string
			// For some address books node id does nt contain node id. In those cases retrieve id from memo field
			if nodeAddress.NodeId == 0 {
				nodeId = string(nodeAddress.Memo)
			} else {
				res, err := entityid.Decode(nodeAddress.NodeId)
				if err != nil {
					return nil, err
				}
				nodeId = res
			}
			result[nodeId] = nodeAddress.RSA_PubKey
		}
	}

	return result, nil
}
