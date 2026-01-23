package tos

import "strings"

func ChangeMany(urls []string) []string {

	//Endpoint: "tos-cn-shanghai.volces.com",
	//Endpoint: "tos-accelerate.volces.com",

	var res []string
	for _, url := range urls {
		//res = append(res, strings.ReplaceAll(url, "tos-cn-shanghai.volces.com", "tos-accelerate.volces.com"))
		//res = append(res, strings.ReplaceAll(url, "tos-accelerate.volces.com", "tos-cn-shanghai.volces.com"))

		res = append(res, Change(url))
	}

	return res
}

func Change(url string) string {

	//Endpoint: "tos-cn-shanghai.volces.com",
	//Endpoint: "tos-accelerate.volces.com",

	return strings.ReplaceAll(url, "tos-accelerate.volces.com", "tos-cn-shanghai.volces.com")
}
