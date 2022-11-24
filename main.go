package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/gonutz/w32/v2"
	"golang.org/x/sys/windows"
)

func main() {
	version := "1.0.0.1"
	hideConsole()
	// if not elevated, relaunch by shellexecute with runas verb set
	if !amAdmin() {
		runMeElevated()
		time.Sleep(10 * time.Second)
		os.Exit(0)
	} else {
		a := app.New()
		a.Settings().SetTheme(theme.DarkTheme())
		w := a.NewWindow("IPs - ver: " + version)

		inter, err := net.Interfaces()
		if err != nil {
			fmt.Println("error")
		}
		placas := make([]string, 0)
		for _, v := range inter {
			fmt.Println(v.Name)
			placas = append(placas, v.Name)
		}
		hello := widget.NewLabel("IP actual")
		adaptador := ""
		combo := widget.NewSelect(placas, func(value string) {
			log.Println("Select set to", value)
			add, _ := GetInterfaceIpv4Addr(value)
			hello.SetText("IP: " + add)
			adaptador = value
		})

		SetDHCP := widget.NewButton("Automatica",
			func() {
				go SetIP(placas, adaptador, "", "", "", hello, "DHCP")
			})

		head := widget.NewGroup("Placa de red", container.NewGridWithColumns(3, combo, SetDHCP, hello))

		buttonA := widget.NewButton("SCHAAF1",
			func() {
				go SetIP(placas, adaptador, "11.11.11.252", "255.255.255.128", "11.11.11.129", hello, "SCHAAF1")
			})
		buttonB := widget.NewButton("MAPA",
			func() {
				go SetIP(placas, adaptador, "12.12.12.252", "255.255.255.128", "12.12.12.129", hello, "MAPA")
			})
		buttonC := widget.NewButton("PC32",
			func() {
				go SetIP(placas, adaptador, "13.13.13.252", "255.255.255.128", "13.13.13.129", hello, "PC32")
			})
		buttonD := widget.NewButton("SCHAAF2",
			func() {
				go SetIP(placas, adaptador, "14.14.14.252", "255.255.255.128", "14.14.14.129", hello, "SCHAAF2")
			})
		buttonD4 := widget.NewButton("DORITOS",
			func() {
				go SetIP(placas, adaptador, "15.15.15.252", "255.255.255.128", "15.15.15.129", hello, "DORITOS")
			})
		buttonD1 := widget.NewButton("PALI RBS",
			func() {
				go SetIP(placas, adaptador, "16.16.16.252", "255.255.255.128", "16.16.16.129", hello, "PALI RBS")
			})
		buttonD2 := widget.NewButton("PCFLEX",
			func() {
				go SetIP(placas, adaptador, "17.17.17.252", "255.255.255.128", "17.17.17.129", hello, "PCFLEX")
			})
		buttonD3 := widget.NewButton("PALITOS3",
			func() {
				go SetIP(placas, adaptador, "18.18.18.252", "255.255.255.128", "18.18.18.129", hello, "PALITOS3")
			})
		buttonD5 := widget.NewButton("FRYPACK",
			func() {
				go SetIP(placas, adaptador, "19.19.19.252", "255.255.255.128", "19.19.19.129", hello, "FRYPACK")
			})
		buttonD6 := widget.NewButton("HALOILA",
			func() {
				go SetIP(placas, adaptador, "24.24.24.252", "255.255.255.128", "24.24.24.129", hello, "HALOILA")
			})
		buttonD7 := widget.NewButton("EFLUENTES",
			func() {
				go SetIP(placas, adaptador, "34.34.34.252", "255.255.255.128", "34.34.34.129", hello, "EFLUENTES")
			})

		ss := container.NewVBox(head, widget.NewSeparator(), buttonA, buttonB, buttonC, buttonD, buttonD1, buttonD2, buttonD3, buttonD4, buttonD5, buttonD6, buttonD7)

		scroll := container.NewScroll(ss)
		scroll.SetMinSize(fyne.NewSize(800, 600))

		w.SetContent(scroll)
		w.ShowAndRun()
	}
}

func SetIP(placas []string, adaptador string, IP string, mask string, gateway string, H *widget.Label, linea string) {
	if adaptador == "" {
		H.SetText("")
		time.Sleep(200 * time.Millisecond)
		H.SetText("Seleccione un placa de red")
		time.Sleep(200 * time.Millisecond)
		H.SetText("")
		time.Sleep(200 * time.Millisecond)
		H.SetText("Seleccione un placa de red")
		time.Sleep(200 * time.Millisecond)
		H.SetText("")
		time.Sleep(200 * time.Millisecond)
		H.SetText("Seleccione un placa de red")
	} else {
		H.SetText("Setting IP")
		if IP != "" {
			for _, v := range placas {
				fmt.Println(v)
				if v != adaptador {
					add, _ := GetInterfaceIpv4Addr(v)
					if add == IP {
						H.SetText("La direccion ya existe en -> " + v)
						return
					}
				}
			}
		}

		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Println("Recovered. Error:\n", r)
		// 	}
		// }()

		var cmd *exec.Cmd
		if IP == "" {
			cmd = exec.Command("netsh", "interface", "ipv4", "set", "address", "name=\""+adaptador+"\"", "source=dhcp")
		} else {
			cmd = exec.Command("netsh", "interface", "ipv4", "set", "address", "\""+adaptador+"\"", "static", IP, mask, gateway, "1")
		}
		err := cmd.Run()
		if err != nil {
			//log.Fatal(err)
			fmt.Println(err)
		}

		time.Sleep(1 * time.Second)
		H.SetText("Setting IP.")
		time.Sleep(1 * time.Second)
		H.SetText("Setting IP..")
		time.Sleep(1 * time.Second)
		H.SetText("Setting IP...")
		time.Sleep(1 * time.Second)

		add, err := GetInterfaceIpv4Addr(adaptador)
		H.SetText("IP: " + add + " -> " + linea)
	}
}

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	return true
}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func hideConsole() {
	console := w32.GetConsoleWindow()
	if console == 0 {
		return // no console attached
	}
	// If this application is the process that created the console window, then
	// this program was not compiled with the -H=windowsgui flag and on start-up
	// it created a console along with the main application window. In this case
	// hide the console window.
	// See
	// http://stackoverflow.com/questions/9009333/how-to-check-if-the-program-is-run-from-a-console
	_, consoleProcID := w32.GetWindowThreadProcessId(console)
	if w32.GetCurrentProcessId() == consoleProcID {
		w32.ShowWindowAsync(console, w32.SW_HIDE)
	}
}

func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", errors.New(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
	}
	return ipv4Addr.String(), nil
}
