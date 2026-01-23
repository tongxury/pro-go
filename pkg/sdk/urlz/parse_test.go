package urlz

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	url, err := ParseURL("https://www.xiaohongshu.com/user/profile/5fa1513f0000000001004dd6?xsec_token=ABlbaXdoUCfXmgGNPPNjFFvygqS4_xny-ushPhyEKxZmk=&xsec_source=pc_note")

	fmt.Println(url, err)
}
