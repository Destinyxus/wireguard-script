package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Destinyxus/wireguardScript/userInt"
	"github.com/Destinyxus/wireguardScript/utils"
)

func main() {

	if _, err := os.Stat("/etc/wireguard/wg0.conf"); err == nil {
		userInt.User()

	} else if os.IsNotExist(err) {
		err := os.Mkdir("/etc/wireguard/qrcodes", 0755)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Directory created successfully")
		}
		// Generate private key
		privateKeyCmd := exec.Command("wg", "genkey")
		privateKeyOut, err := privateKeyCmd.Output()
		if err != nil {
			panic(err)
		}
		privateKey := string(privateKeyOut)

		// Generate public key from private key
		publicKeyCmd := exec.Command("wg", "pubkey")
		publicKeyCmd.Stdin = strings.NewReader(privateKey)
		publicKeyOut, err := publicKeyCmd.Output()
		if err != nil {
			panic(fmt.Errorf("error running command: %v\noutput: %s", err, publicKeyOut))
		}
		publicKey := string(publicKeyOut)

		err = os.Chmod("/etc/wireguard", 0700)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write private key to file
		privateFile, err := os.Create("/etc/wireguard/privatekey")
		if err != nil {
			panic(err)
		}
		defer privateFile.Close()

		_, err = privateFile.WriteString(privateKey)
		if err != nil {
			panic(err)
		}

		// Write public key to file
		publicFile, err := os.Create("/etc/wireguard/publickey")
		if err != nil {
			panic(err)
		}
		defer publicFile.Close()

		_, err = publicFile.WriteString(publicKey)
		if err != nil {
			panic(err)
		}

		networkInterface, err := utils.FindNetworkInterface()
		if err != nil {
			log.Fatal(err)
		}

		// Create WireGuard configuration file
		i := "%i"
		// Create WireGuard configuration file
		configFile := fmt.Sprintf("[Interface]\n"+
			fmt.Sprintf("PrivateKey = %s\n", privateKey)+
			fmt.Sprintf("Address = 10.0.0.1/24\n")+
			fmt.Sprintf("ListenPort = %d\n", 51830)+
			fmt.Sprintf("PostUp = iptables -A FORWARD -i %%s -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE\n", networkInterface)+
			fmt.Sprintf("PostDown = iptables -D FORWARD -i %%s -j ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE\n", networkInterface), i, i)

		err = os.WriteFile("/etc/wireguard/wg0.conf", []byte(configFile), 0600)
		if err != nil {
			panic(err)
		}
		// Enable IP forwarding
		sysctlCmd := exec.Command("sh", "-c", "echo 'net.ipv4.ip_forward=1' >> /etc/sysctl.conf && sysctl -p")
		err = sysctlCmd.Run()
		if err != nil {
			panic(err)
		}

		fmt.Println("WireGuard configuration complete")
		userInt.User()

	}

}
