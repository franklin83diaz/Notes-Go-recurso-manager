package main

import (
	"fmt"
	"log"

	"github.com/mindprince/gonvml"
)

// No se ven los processos que estan utilizando la GPU
func main() {
	// Inicializa NVML
	err := gonvml.Initialize()
	if err != nil {
		log.Fatalf("Error inicializando NVML: %v", err)
	}
	defer gonvml.Shutdown()

	// Obtener el número de GPUs
	deviceCount, err := gonvml.DeviceCount()
	if err != nil {
		log.Fatalf("Error obteniendo la cantidad de GPUs: %v", err)
	}

	// Iterar sobre las GPUs y mostrar información de memoria
	for i := 0; i < int(deviceCount); i++ {
		device, err := gonvml.DeviceHandleByIndex(uint(i))
		if err != nil {
			log.Fatalf("Error obteniendo el manejador de la GPU %d: %v", i, err)
		}

		name, err := device.Name()
		if err != nil {
			log.Fatalf("Error obteniendo el nombre de la GPU %d: %v", i, err)
		}
		totalMemory, usedMemory, err := device.MemoryInfo()
		if err != nil {
			log.Fatalf("Error obteniendo la información de memoria de la GPU %d: %v", i, err)
		}
		// motrar temperatura
		temperature, err := device.Temperature()
		if err != nil {
			log.Fatalf("Error obteniendo la temperatura de la GPU %d: %v", i, err)
		}

		fmt.Printf("GPU %d: %s\n", i, name)
		fmt.Printf("Memoria total: %d MiB\n", totalMemory/(1024*1024))
		fmt.Printf("Memoria usada: %d MiB\n", usedMemory/(1024*1024))
		fmt.Printf("Memoria libre: %d MiB\n", (totalMemory-usedMemory)/(1024*1024))
		fmt.Printf("Temperatura: %d C\n", temperature)
	}
}
