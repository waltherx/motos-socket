package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/waltherx/motos-socket/config"
)

var (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
	urlPost  = ""
)

func init() {
	err := godotenv.Load(".env")
	connHost = os.Getenv("SS_HOST")
	connPort = os.Getenv("SS_PORT")
	urlPost = os.Getenv("URL_POST")

	if err != nil {
		fmt.Println("Error al cargar el archivo .env", err.Error())
	}
}

func main() {
	// Start the server and listen for incoming connections.
	fmt.Println("🤖⚡ Iniciando " + connType + "Servidor -> " + connHost + ":" + connPort)

	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error escuchando:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	// run loop forever, until exit.
	for {
		// Listen for an incoming connection.
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error de conexion:", err.Error())
			return
		}
		fmt.Println("Cliente conectado.")

		// Print client connection address.
		fmt.Println("Cliente " + c.RemoteAddr().String() + " conectado.")

		// Handle connections concurrently in a new goroutine.
		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	// Buffer client input until a newline.
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	// Close left clients.
	if err != nil {
		fmt.Println("Cliente salio.")
		conn.Close()
		return
	}

	// Print response message, stripping newline character.
	data := string(buffer[:len(buffer)-1])

	//fmt.Println(urlPost + " - " + data)
	config.SendPosition(data, urlPost)
	log.Println("💬⚡Cliente message:", data)

	// Send response message to the client.
	conn.Write(buffer)

	// Restart the process.
	handleConnection(conn)
}
