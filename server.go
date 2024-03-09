package main

import (
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
		return
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
	defer conn.Close()
	fmt.Println("Nueva conexión establecida")

	// Leer datos del cliente
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error al leer datos:", err)
		return
	}

	clientMessage := string(buffer[:n])
	fmt.Println("Mensaje del cliente:", clientMessage)
	config.SendPosition(clientMessage, urlPost)
	log.Println("💬⚡Cliente message:", clientMessage)

	// Enviar respuesta al cliente
	response := "ok 200!"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error al enviar respuesta:", err)
		return
	}

	fmt.Println("Respuesta enviada al cliente:", response)
}
