package parser

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	hederaproto "github.com/hashgraph/hedera-sdk-go/proto"
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
			nodeId := string(nodeAddress.Memo)
			result[nodeId] = nodeAddress.RSA_PubKey
		}
	}

	return result, nil
}
