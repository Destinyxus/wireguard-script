package userInt

import (
	"bufio"
	"fmt"
	"image/png"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/Destinyxus/wireguardScript/utils"
)

func User() error {
	reader := bufio.NewReader(os.Stdin)

	// Prompt the user for the client name
	fmt.Print("Enter the client name: ")
	clientName, _ := reader.ReadString('\n')
	clientName = clientName[:len(clientName)-1]

	// Generate private key
	privateKeyClient := exec.Command("wg", "genkey")
	privateKeyOut, err := privateKeyClient.Output()
	if err != nil {
		panic(err)
	}
	privateKeyClientReady := string(privateKeyOut)

	// Generate public key from private key
	publicKeyClient := exec.Command("wg", "pubkey")
	publicKeyClient.Stdin = strings.NewReader(privateKeyClientReady)
	publicKeyOut, err := publicKeyClient.Output()
	if err != nil {
		panic(fmt.Errorf("error running command: %v\noutput: %s", err, publicKeyOut))
	}
	publicKeyClientReady := string(publicKeyOut)

	// Write the private and public keys to disk
	privateKeyFileClient, err := os.Create(fmt.Sprintf("/etc/wireguard/%s_privatekey", clientName))
	if err != nil {
		panic(err)
	}
	defer privateKeyFileClient.Close()

	_, err = privateKeyFileClient.WriteString(privateKeyClientReady)

	publicKeyFile, err := os.Create(fmt.Sprintf("/etc/wireguard/%s_publickey", clientName))
	if err != nil {
		panic(err)
	}
	defer publicKeyFile.Close()
	_, err = publicKeyFile.WriteString(publicKeyClientReady)

	// Open the config file for reading
	confBytes, err := os.ReadFile("/etc/wireguard/wg0.conf")
	if err != nil {
		panic(err)
	}

	// Iterate over each line in the config file
	var lastIP net.IP
	confLines := strings.Split(string(confBytes), "\n")
	for i := range confLines {
		if strings.HasPrefix(confLines[i], "[Peer]") {
			for j := i + 1; j < len(confLines); j++ {
				if strings.HasPrefix(confLines[j], "AllowedIPs = ") {
					// Get the IP address from the AllowedIPs line
					ipString := strings.Split(confLines[j], "AllowedIPs = ")[1]
					lastIP = net.ParseIP(strings.Split(ipString, "/")[0]).To4()

				}
			}
		}
	}

	// If no previous peers were found, use the default IP address
	if lastIP == nil {
		lastIP = net.IPv4(10, 0, 0, 1).To4()
	}

	if lastIP != nil {
		lastIP[3]++
	}

	// Print the last IP address used
	fmt.Printf("Last IP address used: %s\n", lastIP.String())

	clientConf := fmt.Sprintf("\n[Peer]\n" +
		fmt.Sprintf("PublicKey = %s\n", publicKeyClientReady) +
		fmt.Sprintf("AllowedIPs = %s/32\n", lastIP.String()))

	f, err := os.OpenFile("/etc/wireguard/wg0.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(clientConf)
	if err != nil {
		panic(err)
	}

	pbk, err := os.ReadFile("/etc/wireguard/publickey")
	if err != nil {
		log.Fatal(err)
	}
	intf, err := utils.FindNetworkInterface()
	if err != nil {
		log.Fatal(err)
	}
	serverIP, err := utils.FindServerIP(intf)
	if err != nil {
		log.Fatal(err)
	}

	if err := createUserConfig(privateKeyClientReady, lastIP.String(), string(pbk), serverIP, clientName); err != nil {
		log.Fatal(err)
	}

	fmt.Println("WireGuard interface created for client!")
	return nil
}

func createUserConfig(privateKey, clientIP, servPubKey, servIP, clientName string) error {
	config := fmt.Sprintf("[Interface]\nPrivateKey = %s\nAddress = %s/32\nDNS = 8.8.8.8\n\n[Peer]\nPublicKey = %s\nEndpoint = %s:51830\nAllowedIPs = 0.0.0.0/0\nPersistentKeepalive = 20\n",
		privateKey, clientIP, servPubKey, servIP)

	err := os.WriteFile(fmt.Sprintf("/etc/wireguard/userConfigs/config_%s.conf", clientName), []byte(config), 0600)
	// Create QR code image
	qrCode, err := qrcode.New(config, qrcode.Medium)
	if err != nil {
		return err
	}
	fmt.Println(qrCode.ToSmallString(false))
	// Encode QR code image to PNG format
	qrCodePNG := qrCode.Image(256)
	if err != nil {
		return err
	}

	// Save QR code image to file
	f, err := os.Create(fmt.Sprintf("/etc/wireguard/qrcodes/%s_qrcode.png", clientName))
	if err != nil {
		return err
	}
	defer f.Close()

	png.Encode(f, qrCodePNG)

	fmt.Println("QR code also saved to /etc/wireguard/qrcodes/ path")
	return nil

}
