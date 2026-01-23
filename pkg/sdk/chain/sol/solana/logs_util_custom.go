package solana

import "strings"

/*
InitializeMint 指令的作用:
1. 创建新的代币铸造地址：InitializeMint指令用于初始化一个新的代币铸造地址（Mint），并设置其相关参数，如小数位数和管理员地址。
2. 代币的基本信息：在执行InitializeMint时，通常会指定代币的基本信息，例如小数位数（decimals）和管理员（authority）。

需要注意的事项:
1. 上下文：虽然InitializeMint指令通常表示创建代币，但你仍然需要确保上下文是正确的。例如，确保该指令是在一个有效的交易中，并且没有其他异常情况。
2. 其他指令：在同一交易中可能会有其他指令，例如转账、销毁等，这些指令可能会影响代币的状态。因此，最好结合其他信息来确认代币的创建。
3. 日志解析：在解析日志时，确保你能够正确识别InitializeMint指令的上下文。通常，日志中会包含指令的详细信息，包括相关的账户地址和参数。
*/
func IsTokenCreationEvent(logs []string) bool {

	var existsInitializeMint bool
	//var existsPumpFun bool

	for _, x := range logs {
		if strings.Contains(x, "Program log: Instruction: InitializeMint") {
			existsInitializeMint = true
		}
		//else if strings.Contains(x, ProgramPumpFun.String()) {
		//	existsPumpFun = true
		//}
	}

	return existsInitializeMint
}

func IsTokenCreationByPumpFunEvent(logs []string) bool {

	var existsInitializeMint bool
	var existsPumpFun bool

	for _, x := range logs {
		if strings.Contains(x, "Program log: Instruction: InitializeMint") {
			existsInitializeMint = true
		} else if strings.Contains(x, ProgramPumpFun.String()) {
			existsPumpFun = true
		}
	}

	return existsInitializeMint && existsPumpFun
}

func LogsContains(logs []string, keywords []string) bool {
	for _, x := range logs {
		for _, keyword := range keywords {
			if strings.Contains(x, keyword) {
				return true
			}
		}
	}

	return false
}

func IsRadiumPoolCreationEvent(logs []string) bool {

	var existsInitialize bool
	var existsRaydiumPoolProgram bool

	for _, x := range logs {
		if strings.Contains(x, "Program log: initialize2: InitializeInstruction2") {
			existsInitialize = true
		}
		if strings.Contains(x, ProgramRaydiumLiquidityPoolV4.String()) {
			existsRaydiumPoolProgram = true
		}
	}

	return existsInitialize && existsRaydiumPoolProgram
}

func IsNFTCreationEvent(logs []string) bool {
	for _, x := range logs {
		if strings.Contains(x, ProgramNftCandyMachine.String()) {
			return true
		} else if strings.Contains(x, ProgramOrca.String()) {
			// 不是很确定是不是nft的条件, 下面是俩例子
			// 5JNobwKd3eNUhb9DoFjy3XVAkYvBD9N23rhSUP9eQTeFgyJGc8BZaDdufiey6zQRpdbBvdLqGDX3jVwwn4PdCN1X
			// 5YfMPF87bDG2sn8bSQa2Zue3kUVsLKizT2y4jBHnmXXPjcZieqhnibdPxg9XHMtuaFR2riYket3RGFUpy1qji4CQ
			return true
		}

	}

	return false
}
