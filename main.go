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
	"golang.org/x/sys/windows"
)

func main() {
	version := "1.0.0.1"
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
		combo.Size().Width = 400

		head := widget.NewGroup("Placa de red", container.NewHBox(container.NewMax(combo), hello))

		buttonA := widget.NewButton("Schaaf1",
			func() {
				go SetIP(adaptador, "11.11.11.252", "255.255.255.128", "11.11.11.1", hello)
			})
		buttonB := widget.NewButton("Mapa",
			func() {
				go SetIP(adaptador, "12.12.12.252", "255.255.255.128", "12.12.12.1", hello)
			})
		buttonC := widget.NewButton("PC32",
			func() {
				go SetIP(adaptador, "13.13.13.252", "255.255.255.128", "13.13.13.1", hello)
			})
		buttonD := widget.NewButton("Schaaf2",
			func() {
				go SetIP(adaptador, "14.14.14.252", "255.255.255.128", "14.14.14.1", hello)
			})
		buttonD4 := widget.NewButton("Schaaf2",
			func() {
				go SetIP(adaptador, "14.14.14.252", "255.255.255.128", "14.14.14.1", hello)
			})
		buttonD1 := widget.NewButton("Schaaf2",
			func() {
				go SetIP(adaptador, "14.14.14.252", "255.255.255.128", "14.14.14.1", hello)
			})
		buttonD2 := widget.NewButton("Schaaf2",
			func() {
				go SetIP(adaptador, "14.14.14.252", "255.255.255.128", "14.14.14.1", hello)
			})
		buttonD3 := widget.NewButton("Schaaf2",
			func() {
				go SetIP(adaptador, "14.14.14.252", "255.255.255.128", "14.14.14.1", hello)
			})

		ss := container.NewVBox(head, widget.NewSeparator(), buttonA, buttonB, buttonC, buttonD, buttonD1, buttonD2, buttonD3, buttonD4)

		scroll := container.NewScroll(ss)
		scroll.SetMinSize(fyne.NewSize(800, 600))

		w.SetContent(scroll)
		w.ShowAndRun()
	}
}

func SetIP(adaptador string, IP string, mask string, gateway string, H *widget.Label) {
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

		H.SetText("Setting IP...")
		cmd := exec.Command("netsh", "interface", "ipv4", "set", "address", "name="+adaptador, "static", IP, mask, gateway)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)

		add, err := GetInterfaceIpv4Addr(adaptador)
		H.SetText("IP: " + add)
	}
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

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	return true
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
