package main

import (
	"fmt"
	"net"
)

func getLocalNetWork() {
	interfaces, error := net.Interfaces()

	if error != nil {
		fmt.Println("Erro ao listar interfaces", error)
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()

		if err != nil {
			fmt.Println("Erro ao obter endereços:", err)
			continue
		}

		fmt.Println("Endereços:")
		for _, addr := range addrs {
			fmt.Println("  -", addr.String())
		}
	}
}

func main() {
	getLocalNetWork()
}
