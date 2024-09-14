package main

/*
#cgo CFLAGS: -I/usr/local/cuda/include
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lnvidia-ml
#include <nvml.h>
*/
import "C"
import (
	"fmt"
	"log"
	"os"
)

// nvmlDeviceGetComputeRunningProcesses
// nvmlDeviceGetGraphicsRunningProcesses
// nvmlDeviceGetProcessUtilization

//install driver nvidia: sudo ubuntu-drivers autoinstall
//need to install sudo apt install nvidia-cuda-toolkit

func main() {
	// Inicializar NVML
	if res := C.nvmlInit(); res != C.NVML_SUCCESS {
		log.Fatalf("No se pudo inicializar NVML: %v", nvmlErrorString(res))
	}
	defer func() {
		if res := C.nvmlShutdown(); res != C.NVML_SUCCESS {
			log.Printf("No se pudo cerrar NVML: %v", nvmlErrorString(res))
		}
	}()

	// Obtener el dispositivo (GPU 0)
	var device C.nvmlDevice_t
	if res := C.nvmlDeviceGetHandleByIndex(0, &device); res != C.NVML_SUCCESS {
		log.Fatalf("No se pudo obtener el dispositivo: %v", nvmlErrorString(res))
	}

	// Obtener la cantidad de procesos de computación que están utilizando la GPU
	var infoCount C.uint = 0
	// Llamada inicial para obtener el número de procesos
	res := C.nvmlDeviceGetComputeRunningProcesses(device, &infoCount, nil)
	if res != C.NVML_SUCCESS && res != C.NVML_ERROR_INSUFFICIENT_SIZE {
		log.Fatalf("No se pudo obtener el número de procesos: %v", nvmlErrorString(res))
	}

	if infoCount == 0 {
		fmt.Println("No hay procesos de computación utilizando la GPU.")
	} else {
		// Crear un slice para almacenar la información de los procesos
		processInfos := make([]C.nvmlProcessInfo_t, infoCount)
		// Llamada para obtener la información de los procesos
		res = C.nvmlDeviceGetComputeRunningProcesses(device, &infoCount, &processInfos[0])
		if res != C.NVML_SUCCESS {
			log.Fatalf("No se pudo obtener la información de los procesos: %v", nvmlErrorString(res))
		}

		fmt.Printf("Número de procesos de computación utilizando la GPU: %d\n", infoCount)

		// Mostrar información de cada proceso
		for i := uint(0); i < uint(infoCount); i++ {
			pid := uint32(processInfos[i].pid)
			memUsed := uint64(processInfos[i].usedGpuMemory) // En bytes

			// Obtener el nombre del proceso a partir del PID
			procName := getProcessName(pid)

			fmt.Printf("Proceso %d: PID %d, Memoria GPU usada: %d MB, Nombre: %s\n", i+1, pid, memUsed/(1024*1024), procName)
		}
	}

	// Repetir el proceso para los procesos gráficos
	var graphicsInfoCount C.uint = 0
	res = C.nvmlDeviceGetGraphicsRunningProcesses(device, &graphicsInfoCount, nil)
	if res != C.NVML_SUCCESS && res != C.NVML_ERROR_INSUFFICIENT_SIZE {
		log.Fatalf("No se pudo obtener el número de procesos gráficos: %v", nvmlErrorString(res))
	}

	if graphicsInfoCount == 0 {
		fmt.Println("No hay procesos gráficos utilizando la GPU.")
	} else {
		graphicsProcessInfos := make([]C.nvmlProcessInfo_t, graphicsInfoCount)
		res = C.nvmlDeviceGetGraphicsRunningProcesses(device, &graphicsInfoCount, &graphicsProcessInfos[0])
		if res != C.NVML_SUCCESS {
			log.Fatalf("No se pudo obtener la información de los procesos gráficos: %v", nvmlErrorString(res))
		}

		fmt.Printf("Número de procesos gráficos utilizando la GPU: %d\n", graphicsInfoCount)

		for i := uint(0); i < uint(graphicsInfoCount); i++ {
			pid := uint32(graphicsProcessInfos[i].pid)
			memUsed := uint64(graphicsProcessInfos[i].usedGpuMemory) // En bytes

			// Obtener el nombre del proceso a partir del PID
			procName := getProcessName(pid)

			fmt.Printf("Proceso Gráfico %d: PID %d, Memoria GPU usada: %d MB, Nombre: %s\n", i+1, pid, memUsed/(1024*1024), procName)
		}
	}
}

// Función auxiliar para convertir códigos de error NVML a cadenas
func nvmlErrorString(res C.nvmlReturn_t) string {
	return C.GoString(C.nvmlErrorString(res))
}

// Función para obtener el nombre del proceso dado su PID
func getProcessName(pid uint32) string {
	procPath := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdline, err := os.ReadFile(procPath)
	if err != nil {
		return "Desconocido"
	}
	// Los argumentos del comando están separados por bytes nulos
	for i := 0; i < len(cmdline); i++ {
		if cmdline[i] == 0 {
			cmdline[i] = ' '
		}
	}
	return string(cmdline)
}
