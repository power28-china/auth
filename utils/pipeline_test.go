package utils

import (
	"testing"

	"github.com/power28-china/auth/utils/logger"
)

func TestMongoPipeline(t *testing.T) {
	srt := `
	[{
		"$match":{
			"$and":
				[{"provinces":"湖北省"},{"dealstatus":{"$ne":"1"}}]
		}
	},{
			"$count":"total"
	}]`

	res := MongoPipeline(srt)
	logger.Sugar.Debugf("MarshalBsonD: %#v\n", res)

	if res == nil {
		t.Errorf("MarshalBsonD failed.\n")
	}
}
