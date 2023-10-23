package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"math/rand"
    "time"

	"google.golang.org/grpc"
)

func main() {

    personas := readArchivoNames() //Lista de los nombres de personas
	estado := generarProbabilidadesPersonas(len(personas)) //Lista del estado de la persona, la persona[i] tiene estado[i]
    rand.Seed(time.Now().UnixNano()) //Seed para el random

    conn, err := grpc.Dial("10.6.46.139:50053", grpc.WithInsecure()) //Conecta con la OMS
    if err != nil {
        log.Fatalf("No se pudo conectar: %v", err)
    }
    defer conn.Close()

    c := proto.NewMyServiceClient(conn)


    for i := 0; i < 5; i++ { //Manda las primeros 5 personas
        rd := rand.Intn(407) //Numero aleatorio para mandar la persona
        nomap := strings.Split(personas[rd], " ") //Split el nombre guardado en la lista personas
        r, err := c.SendContinentMsg(context.Background(), proto.MensajeContinente{ //Manda el mensaje
            Nombre: nomap[0],
            Apellido: nomap[1],
            Estado:   func() string {
                if estado[rd] <= 0.55 {
                    return "Infectado"
                }
                return "Fallecido"
            }(),
        })

        if err != nil {
            log.Fatalf("No se pudo enviar el mensaje: %v", err)
        }   
        log.Printf("Mensaje enviado con éxito, código número %s", r.GetMensaje())
    }

    sendMessages := func() {
        for {
            rd := rand.Intn(407)
            nomap := strings.Split(personas[rd], " ")
            r, err := c.SendContinentMsg(context.Background(), proto.MensajeContinente{
                Nombre:   nomap[0],
                Apellido: nomap[1],
                Estado: func() string {
                    if estado[rd] <= 0.55 {
                        return "Infectado"
                    }
                    return "Fallecido"
                }(),
            })

            if err != nil {
                log.Fatalf("No se pudo enviar el mensaje: %v", err)
            }
            log.Printf("Mensaje enviado con éxito, código número %s", r.GetMensaje())

            time.Sleep(3 * time.Second) // Espera 3 segundos antes de enviar el siguiente mensaje
        }
    }

    go sendMessages() //Inicia la gorutina y manda mensajes cada 3 segundos

    select {} //Mantiene el programa en ejecucion para seguir mandando mensajes

}

func generarProbabilidadesPersonas(n int) []float32 {
    fuente := rand.NewSource(time.Now().UnixNano())
    generador := rand.New(fuente)

    probabilidades := make([]float32, n)

    for i := 0; i < n; i++ {
        probabilidades[i] = generador.Float32()
    }

    return probabilidades
}

func readArchivoNames() []string {

    file, err := os.Open("names.txt") //abre el archivo names.txt

    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    reader := bufio.NewReader(file) //reader del archivo names.txt

    var nombres []string //lista de nombres vac\'ia

    for {
        name, err := reader.ReadString('\n') //lee una linea del archivo
        if err != nil {
            break
        }

		name = strings.TrimRight(name, "\n") //elimina el salto de linea del nombre le\'ido

        nombres = append(nombres, name) //agrega el nombre a la lista nombres
    }

    return nombres //retorna la lista con todos los nombres
}