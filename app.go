package main

import (
	"fmt"
	"log"

	"github.com/go-routeros/routeros"
	"github.com/spf13/viper"
)

type Configurations struct {
	App AppConfiguration
}

type AppConfiguration struct {
	Address string
	Username string
	Password string
}

func main() {

	ConfigRouter()
}
func ConfigRouter() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read configuration file: %v", err)
	}

	var config Configurations
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// เชื่อมต่อรายละเอียด MikroTik
	Address := config.App.Address
	Username := config.App.Username
	Password := config.App.Password

	// เชื่อมต่อ MikroTik
	client, err := routeros.Dial(Address, Username, Password)
	if err != nil {
		log.Fatalf("Unable to connect to the router: %v", err)
	}
	defer client.Close()

	// บล็อคการส่งคำข้อเข้ามา
	blockHTTPandHTTPS(client)

	fmt.Println("Blocking HTTP and HTTPS requests completed.")
}

func blockHTTPandHTTPS(client *routeros.Client) {
	// บล็อคการส่งคำข้อเข้ามา
	_, err := client.Run("/ip/firewall/filter/add", "=chain=input", "=protocol=tcp", "=dst-port=80", "=action=drop")
	if err != nil {
		log.Fatalf("Unable to block HTTP requests: %v", err)
	}

	_, err = client.Run("/ip/firewall/filter/add", "=chain=input", "=protocol=tcp", "=dst-port=443", "=action=drop")
	if err != nil {
		log.Fatalf("Unable to block HTTPS requests: %v", err)
	}
}