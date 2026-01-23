package solana

import (
	"fmt"
	"log"
	"testing"
)

func TestCreateReferralAccount(t *testing.T) {
	// 定义推荐账户和铸币地址

	/*
		1. 推荐账户（Referral Account）
		作用：推荐账户通常用于奖励推广者或用户引导新用户使用某个应用或服务。通过引用链接或代码，新用户的交易或活动可以关联到这个推荐账户，从而让推荐者获得一定的奖励或费用返还。
		场景：在许多去中心化金融（DeFi）应用、交易平台或NFT平台中，推荐系统用于激励用户邀请新用户。完成交易后，推荐者可能会收到一定比例的费用（如手续费）作为奖励。
		使用方式：在生成费代币账户时，推荐账户的地址会成为生成地址的种子之一，确保与该账户关联的费用可以正确地分配给推荐者。

		2. 铸币地址（Mint Address）
		作用：铸币地址是某种代币的唯一标识符。所有基于该铸币的代币都会使用同一个铸币地址。该地址决定了代币的属性，如名称、符号、总供应量等。
		场景：在进行代币交易、转账或交换时，铸币地址用于定义交易中的代币类别。每种代币在区块链上都有唯一的铸币地址，确保代币的不同类型和性质可以被正确识别。
		使用方式：在生成费代币账户时，铸币地址作为一个种子之一参与计算，确保生成的费代币账户与具体的代币类型和交易相关联。

	*/
	referralAccount := "EELg1rxpMGhygQdBFQF8UofyPArvTFf5VJLGJ4NsEWB" // 例如 "3N1S1...XYZ"
	mint := "eaVWnkV45z6m1FAMHAJxCqDUe6CGGkCRK5NbiLUB44C"            // 例如 "So11111111111111111111111111111111111111112"

	// 将地址转换为 PublicKey 类型
	referralPubKey, err := PublicKeyFromBase58(referralAccount)
	if err != nil {
		log.Fatalf("无法转换推荐账户地址: %v", err)
	}

	mintPubKey, err := PublicKeyFromBase58(mint)
	if err != nil {
		log.Fatalf("无法转换铸币地址: %v", err)
	}

	// 定义种子
	seeds := [][]byte{
		[]byte("referral_ata"),
		referralPubKey[:],
		mintPubKey[:],
	}

	// 使用推荐合约的公共密钥，替换为你的推荐合约地址
	programID := MustPublicKeyFromBase58("REFER4ZgmyYx9c6He5XfaTMiGfdLwRnkV4RPp9t9iF3")

	// 生成费代币账户地址
	feeTokenAccount, _, err := FindProgramAddress(seeds, programID)
	if err != nil {
		log.Fatalf("生成费代币账户地址失败: %v", err)
	}

	// 输出生成的费代币账户地址
	fmt.Printf("生成的费代币账户地址: %s\n", feeTokenAccount.String())
}
