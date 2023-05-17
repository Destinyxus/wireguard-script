# WIREGUARD-SCRIPT

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/destinyxus/wireguard-script?color=green&logo=green&logoColor=violet&style=plastic)

## About

**wireguard-script** - is a mini-script in golang which implements a [wireguard](https://wireguard.com/#simple-network-interface) server setup logic and automated creation of client's config

## Installation

### First method:


SSH into the Linux server, after logging in, check if the machine is up to date by running the following command:
````
sudo apt-get update && sudo apt-get upgrade
````
Now install WireGuard by running this one:
````
sudo apt-get install wireguard
````
Install a command-line utility  'iptables' if it's not already installed on your machine:
````
sudo apt-get install iptables
````
Install git:
````
sudo apt install git
````
Install golang on machine:
````
sudo apt install golang-go
````
Go to any directory where you want to save the source code of the program and clone it from github:
````
git clone https://github.com/Destinyxus/wireguard-script.git
````
Go directly to the directory with the script code:
````
cd wireguard-script
````
Launch it:
````
go run main.go
````
***Good job! All config files you create will be saved in "/etc/wireguard/userConfigs" and QR-codes in ".../qrcodes" respectively.*** 

***Now you can share them with your users!***

## P.S.
### **Do not forget to restart wireguard system for all changes to take effect:**
````
systemctl restart wg-quick@wg0
````


## Second method:
If you are facing some difficulties with installation, you can easily download a binary file which implements the same script-logic.



* Go to the [releases page](https://github.com/Destinyxus/wireguard-script/releases) of this repository.


* Find the latest release and download the binary file for your system in the certain directory you like.


* Make this binary file executable with a following command:
````
chmod +x ./"script_name"
````
Run it with sudo privileges:
````
sudo ./"script_name"
````

## Example
Downloading:
```bash
wget https://github.com/Destinyxus/wireguard-script/releases/download/latest/"binary_name"
chmod +x ./"binary_name"
sudo ./"binary_name"
```


## P.S.S.

**_You can face an error while downloading like this one:_**
````
"Resolving github.com (github.com)... failed: Temporary failure in name resolution.
wget: unable to resolve host address ‘github.com"
````
To fix it, you should update DNS configuration on your server!

You can check the DNS settings in the /etc/resolv.conf file.
````bash
sudo vim /etc/resolv.conf
````    
Append these two lines to this configuration:
````
nameserver 8.8.8.8
nameserver 8.8.4.4
````
Save it and restart the networking service on your Ubuntu server to apply any changes:
````
sudo systemctl restart networking
````
Try to download a binary file again! `:dizzy:`

## Developers
 [Destinyxus](https://github.com/Destinyxus)

## License

Project "wireguard-script" is distributed under the MIT license.
