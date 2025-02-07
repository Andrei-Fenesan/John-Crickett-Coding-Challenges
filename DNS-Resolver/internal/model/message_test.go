package model

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDnsMessage(t *testing.T) {
	assert := assert.New(t)
	encodedMessage := NewQuestion(22, "dns.google.com")

	assert.Equal("00160000000100000000000003646e7306676f6f676c6503636f6d0000010001", fmt.Sprintf("%x", encodedMessage.Encode()))
}

func TestParseMessageResponse(t *testing.T) {
	assert := assert.New(t)
	response, err := hex.DecodeString("00168180000100020000000003646e7306676f6f676c6503636f6d0000010001c00c0001000100000081000408080404c00c000100010000008100040808080800")
	if err != nil {
		t.FailNow()
	}
	message := ParseResponse(response)

	assert.NotNil(message.header, message.question, message.answer, message.authority, message.additional)
	assert.Equal(1, len(message.question))
	assert.Equal(2, len(message.answer))
	assert.Empty(message.authority, message.additional)
}

func TestParseResponseWithResourceHavingMultiplePonters(t *testing.T) {
	assert := assert.New(t)

	response, err := hex.DecodeString("0016820000010000000d000b03646e7306676f6f676c6503636f6d0000010001c017000200010002a3000014016c0c67746c642d73657276657273036e657400c017000200010002a3000004016ac02ec017000200010002a30000040168c02ec017000200010002a30000040164c02ec017000200010002a30000040162c02ec017000200010002a30000040166c02ec017000200010002a3000004016bc02ec017000200010002a3000004016dc02ec017000200010002a30000040169c02ec017000200010002a30000040167c02ec017000200010002a30000040161c02ec017000200010002a30000040163c02ec017000200010002a30000040165c02ec02c000100010002a3000004c029a21ec02c001c00010002a300001020010500d93700000000000000000030c04c000100010002a3000004c0304f1ec04c001c00010002a300001020010502709400000000000000000030c05c000100010002a3000004c036701ec05c001c00010002a30000102001050208cc00000000000000000030c06c000100010002a3000004c01f501ec06c001c00010002a300001020010500856e00000000000000000030c07c000100010002a3000004c0210e1ec07c001c00010002a300001020010503231d00000000000000020030c08c000100010002a3000004c023331e")
	if err != nil {
		t.FailNow()
	}
	message := ParseResponse(response)

	additional := message.additional[1]
	assert.Equal("l.gtld-servers.net", additional.Name)
}
