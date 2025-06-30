package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Ingrese la dirección del nodo (host:puerto), ej. localhost:19002: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	fmt.Println("Comandos disponibles:")
	fmt.Println(" - ENTRENAR")
	fmt.Println(" - PREDECIR:estrato,dominio (ej: PREDECIR:3.0,450.0)")

	for {
		fmt.Print("Ingrese comando: ")
		comando, _ := reader.ReadString('\n')
		comando = strings.TrimSpace(comando)
		if comando == "" {
			continue
		}

		enviar(comando, address)
	}
}

func enviar(comando, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("❌ Error al conectar:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintln(conn, comando)
	fmt.Println("✅ Comando enviado:", comando)
}
