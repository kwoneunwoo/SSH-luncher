package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)


// Server 구조체 정의 (JSON 매핑용)
type Server struct {
	Name    string `json:"name"`
	SSHPath string `json:"ssh_path"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type Config struct {
	Servers []Server `json:"servers"`
}


// clearScreen은 OS에 맞춰 터미널 화면을 지웁니다.
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// 엔터를 치면 다음으로 넘어가는 코드.
func pause() {
	fmt.Print("\nPress Enter to exit... ")
	var dummy string
	fmt.Scanln(&dummy)
}

// loadServers는 JSON 파일에서 서버 정보를 읽어옵니다.
func loadServers(filename string) []Server {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("[Error] File '%s' not found.", filename)
		pause()
		return nil
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("[Error] An error occurred while reading the file: %v", err)
		pause()
		return nil
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("[Error] An error occurred while parsing the JSON file: %v", err)
		pause()
		return nil
	}

	return config.Servers
}

func main() {
	servers := loadServers("servers.json")

	var target Server
	var choiceIndex int

	switch len(servers) {
	case 0: // Json에 리스트가 없다면 오류 메시지 출력
		fmt.Println("No server information to load.")
		pause()
		return
	case 1: // 리스트 값이 한 개라면 자동 선택
		target = servers[0]
		choiceIndex = 0
	default: // 리스트 값이 여러 개라면 사용자가 선택
		for {
			// Print to console
			clearScreen()
			fmt.Println("========================================")
			fmt.Println("       SSH Server connecting menu       ")
			fmt.Println("========================================")

			for idx, server := range servers {
				fmt.Printf("[%d] %s\n", idx+1, server.Name)
			}

			fmt.Println("========================================")
			fmt.Println()

			fmt.Printf("Choose server number (1-%d): ", len(servers))

			// 사용자가 입력 후 숫자로 변환하여 저장.
			var inputStr string
			_, err := fmt.Scanln(&inputStr)
			if err != nil {
				continue
			}

			choice, err := strconv.Atoi(strings.TrimSpace(inputStr))
			if err != nil {
				continue
			}

			choiceIndex := choice - 1
			if choiceIndex < 0 || choiceIndex >= len(servers) {
				continue
			}

			break
		}
	}

	// 선택한 서버 정보를 정리하여 접속
	target = servers[choiceIndex]
	sshKey := target.SSHPath
	address := target.Address
	port := strconv.Itoa(target.Port)

	clearScreen()
	fmt.Printf("[%d] Connecting to %s ...\n\n", choiceIndex+1, target.Name)

	// SSH 명령어 실행 (인터랙티브 입출력 연동)
	cmd := exec.Command("ssh", "-i", sshKey, address, "-p", port)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("[Error] An error occurred while closing the SSH connection: %v", err)
	}
	pause()
}
