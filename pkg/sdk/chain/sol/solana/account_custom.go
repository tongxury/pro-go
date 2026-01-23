package solana

func IsValidPublicKey(in string) bool {

	_, err := PublicKeyFromBase58(in)
	if err != nil {
		return false
	}

	return true
}
