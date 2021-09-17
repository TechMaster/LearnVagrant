package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

/*
Trả về true nếu ứng dụng đang chạy ở chế độ Debug và ngược lại
*/
func IsAppInDebugMode() bool {
	appCommand := os.Args[0]
	if strings.Contains(appCommand, "debug") || //debug ứng dụng trong vscode
		strings.Contains(appCommand, "exe") || //go run main.go
		strings.Contains(appCommand, "go-build") { //run test
		return true
	}
	return false
}

func ReadConfig(configPaths ...string) {
	var configPath string
	if len(configPaths) == 0 {
		configPath = "."
	} else {
		configPath = configPaths[0]
	}

	if IsAppInDebugMode() {
		viper.SetConfigName("config.dev") // Debug - development mode
	} else {
		viper.SetConfigName("config.product") // Product mode
	}

	viper.SetConfigType("json") // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(configPath) // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	ParseViperSettings()
}

/*
Đọc toàn bộ dữ liệu mà viper lưu.
Nếu là kiểu value là string thì kiểm tra nó có chứa 2 ký tự đặc biệt @@ ở đầu không?
Sau hai ký tự @@ là tên file lưu ở thư mục /run/secrets/

Nếu là kiểu map[string]interface{} thì đọc sâu vào
*/
func ParseViperSettings() {
	for key, value := range viper.AllSettings() {
		parseSecretConfig(key, value)
	}
}

func parseSecretConfig(key string, value interface{}) {
	if data, ok := value.(map[string]interface{}); ok {
		for subKey, subValue := range data {
			parseSecretConfig(key+"."+subKey, subValue)
		}
	}

	if str, ok := value.(string); ok {
		if len(str) > 2 && str[:2] == "@@" { //Nếu value chứa dấu hiệu trỏ tới file secret ở thư mục /run/secrets
			secret := readSecretFile(str[2:])
			viper.Set(key, secret)
		}
	}
}

/*
truyền vào key, đọc giá trị của file lưu ở /run/secrets/key
Nếu không tìm được, đọc lỗi thì trả về empty string
*/
func readSecretFile(fileName string) string {
	if body, err := ioutil.ReadFile("/run/secrets/" + fileName); err != nil {
		return ""
	} else {
		temp := string(body)
		if len(temp) > 1 && temp[len(temp)-1:] == "\n" { //loại bỏ ký tự \n cuối cùng
			return temp[:len(temp)-1]
		} else {
			return temp
		}
	}
}
