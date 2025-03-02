package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Printer struct {
	Name   string
	Class  string
	Port   string
	Status string
}

func detectUSBPrinters() []Printer {
	cmd := exec.Command("powershell", "-Command",
		"Get-PnpDevice | Format-Table Name, Status, Class -AutoSize | Out-String -Width 4096")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Erro ao listar dispositivos:", err)
		fmt.Println("Saída do comando:", string(output))
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var printers []Printer

	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		status := fields[len(fields)-2]
		class := fields[len(fields)-1]
		name := strings.Join(fields[:len(fields)-2], " ")

		lowerClass := strings.ToLower(class)
		if lowerClass == "usb" || lowerClass == "ports" {
			port := "USB001"
			lowerName := strings.ToLower(name)
			if strings.Contains(lowerName, "com") {
				start := strings.Index(lowerName, "com")
				end := strings.Index(lowerName[start:], ")")
				if end == -1 {
					end = len(lowerName) - start
				}
				port = strings.ToUpper(lowerName[start : start+end])
			}
			printers = append(printers, Printer{
				Name:   name,
				Class:  class,
				Port:   port,
				Status: status,
			})
		}
	}
	return printers
}

func installDriver(printer Printer) error {
	var scriptPath string
	lowerName := strings.ToLower(printer.Name)

	if strings.Contains(lowerName, "bematech") {
		scriptPath = "./install_bematech.exe"
	} else {
		return fmt.Errorf("Nenhum driver conhecido para %s", printer.Name)
	}

	if _, err := os.Stat(scriptPath[2:]); os.IsNotExist(err) {
		return fmt.Errorf("Script %s não encontrado na pasta.", scriptPath[2:])
	}

	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Start-Process -FilePath '%s' -Verb RunAs", scriptPath[2:]))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Erro ao executar script para %s: %v\nSaída: %s", printer.Name, err, string(output))
	}
	fmt.Printf("Driver instalado para %s!\n", printer.Name)

	time.Sleep(20 * time.Second)
	return nil
}

func main() {
	fmt.Println("Procurando impressoras USB e COM...")
	usbPrinters := detectUSBPrinters()

	if len(usbPrinters) == 0 {
		fmt.Println("Nenhuma impressora USB ou COM encontrada.")
		return
	}

	fmt.Println("Impressoras detectadas:")
	for _, printer := range usbPrinters {
		fmt.Printf("Nome: %s, Classe: %s, Porta: %s, Status: %s\n", printer.Name, printer.Class, printer.Port, printer.Status)
	}

	for _, printer := range usbPrinters {
		lowerName := strings.ToLower(printer.Name)
		if strings.Contains(lowerName, "bematech") {
			if err := installDriver(printer); err != nil {
				fmt.Println(err)
			}
		}
	}
}
