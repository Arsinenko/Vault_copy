package main

import (
	service_user "Vault_copy/services/user"
)

func main() {
	res1 := service_user.CreateUser("+79775509028", "12345678", "Viktor")
	println(res1)

	res2 := service_user.AuthStandard("+79775509028", "12345678")
	println(res2)
}