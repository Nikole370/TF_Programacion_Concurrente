// armando la red P2P
package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"PROJECTFINAL/Modelo"
	"PROJECTFINAL/cRedis"
)

var (
	addrs  []string // IPs de todos los nodos (incluyéndome)
	hostIP string   // IP de este nodo
)

const portHP = 9002 // puerto de escucha

// sendMessage envía msg al peer ip:portHP
func sendMessage(msg, ip string) {
	addr := fmt.Sprintf("%s:%d", ip, portHP)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("❌ Error conectando a %s: %v\n", addr, err)
		return
	}
	defer conn.Close()
	fmt.Fprintln(conn, msg)
}

// broadcastWeights manda los pesos a todos los peers excepto a mí
func broadcastWeights(pesos []float64) {
	parts := make([]string, len(pesos))
	for i, p := range pesos {
		parts[i] = strconv.FormatFloat(p, 'f', -1, 64)
	}
	msg := "PESOS:" + strings.Join(parts, ",")
	for _, ip := range addrs {
		if ip == hostIP {
			continue
		}
		go sendMessage(msg, ip)
	}
}

// doTrainAndBroadcast encapsula el entrenamiento local + difusión de pesos
func doTrainAndBroadcast() {
	idx := obtenerPartitionIdx()
	pesos, minE, maxE, minD, maxD, _ := Modelo.EntrenarModelo("datos.csv", idx, len(addrs))
	cRedis.GuardarPesos("modelo", pesos)
	cRedis.GuardarMinMax("minEstrato", minE)
	cRedis.GuardarMinMax("maxEstrato", maxE)
	cRedis.GuardarMinMax("minDominio", minD)
	cRedis.GuardarMinMax("maxDominio", maxD)
	fmt.Printf("⚙️ Nodo %s entrenó partición %d → pesos: %v\n", hostIP, idx, pesos)
	broadcastWeights(pesos)
}

func main() {
	// 1) Iniciar Redis
	cRedis.IniciarRedis()

	// 2) Descubrir IP local
	hostIP = descubrirIP()
	fmt.Printf("Mi IP es %s\n", hostIP)

	// 3) Declarar todas las IPs del cluster (incluyéndome)
	//addrs = []string{"172.20.0.2", "172.20.0.3", "172.20.0.4"}
	addrs = []string{"nodo1", "nodo2", "nodo3"}

	// 4) Iniciar servidor P2P
	go servicioHP()
	fmt.Println("🟢 Nodo en espera de comandos P2P…")

	select {} // bloqueo infinito
}

func servicioHP() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", hostIP, portHP))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err == nil {
			go handlerHP(conn)
		}
	}
}

func handlerHP(conn net.Conn) {
	defer conn.Close()
	raw, _ := bufio.NewReader(conn).ReadString('\n')
	msg := strings.TrimSpace(raw)

	switch {
	// 1) Cliente dispara "ENTRENAR"
	case msg == "ENTRENAR":
		fmt.Printf("📡 Recibido ENTRENAR CLI en %s\n", hostIP)
		// reenvío **una sola vez** a todos los peers
		for _, ip := range addrs {
			if ip != hostIP {
				go sendMessage("ENTRENAR:"+hostIP, ip)
			}
		}
		// entreno y difundo pesos
		doTrainAndBroadcast()

	// 2) Peer recibe "ENTRENAR:origin" — entreno pero **no** vuelvo a re-broadcast
	case strings.HasPrefix(msg, "ENTRENAR:"):
		origin := strings.TrimPrefix(msg, "ENTRENAR:")
		fmt.Printf("📡 Recibido ENTRENAR de %s\n", origin)
		doTrainAndBroadcast()

	// 3) Peer recibe pesos
	case strings.HasPrefix(msg, "PESOS:"):
		data := strings.TrimPrefix(msg, "PESOS:")
		parts := strings.Split(data, ",")
		rec := make([]float64, len(parts))
		for i, p := range parts {
			rec[i], _ = strconv.ParseFloat(p, 64)
		}
		fmt.Printf("📦 Nodo %s recibe PESOS: %v\n", hostIP, rec)

		local := cRedis.LeerPesos("modelo")
		if len(local) != len(rec) {
			fmt.Println("⚠️ Tamaño de pesos no coincide, omito fusión")
			return
		}
		merged := make([]float64, len(local))
		for i := range local {
			merged[i] = (local[i] + rec[i]) / 2
		}
		cRedis.GuardarPesos("modelo", merged)
		fmt.Printf("✅ Nodo %s fusionó pesos → %v\n", hostIP, merged)

	// 4) Petición de predicción
	case strings.HasPrefix(msg, "PREDECIR:"):
		data := strings.TrimPrefix(msg, "PREDECIR:")
		vals := strings.Split(data, ",")
		if len(vals) != 2 {
			fmt.Println("⚠️ Formato PREDECIR inválido. Usa PREDECIR:3.0,450.0")
			return
		}
		estrato, _ := strconv.ParseFloat(vals[0], 64)
		dominio, _ := strconv.ParseFloat(vals[1], 64)

		pesos := cRedis.LeerPesos("modelo")
		minE := cRedis.LeerMinMax("minEstrato")
		maxE := cRedis.LeerMinMax("maxEstrato")
		minD := cRedis.LeerMinMax("minDominio")
		maxD := cRedis.LeerMinMax("maxDominio")

		pred := Modelo.Predecir(estrato, dominio, pesos, minE, maxE, minD, maxD)
		fmt.Printf("🔮 Predicción nodo %s: estrato=%.1f dominio=%.1f → %.4f\n",
			hostIP, estrato, dominio, pred)

	default:
		fmt.Println("⚠️ Comando no reconocido:", msg)
	}
}

func obtenerPartitionIdx() int {
	switch hostIP {
	case "172.20.0.2":
		return 0
	case "172.20.0.3":
		return 1
	case "172.20.0.4":
		return 2
	default:
		return 0
	}
}

func descubrirIP() string {
	var dirIP = "127.0.0.1"
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "eth0") {
			addrs, _ := iface.Addrs()
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return dirIP
}
