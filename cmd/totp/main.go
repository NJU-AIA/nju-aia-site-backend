package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"ArticleServer/internal/auth"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	genKey := flag.Bool("generate-key", false, "生成新的 ECC P-256 密钥")
	flag.Parse()

	if *genKey {
		key, err := auth.GenerateECCKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "生成密钥失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("将以下内容添加到 .env 文件:")
		fmt.Println("TOTP_ECC_KEY=" + key)
		return
	}

	eccKeyB64 := os.Getenv("TOTP_ECC_KEY")
	if eccKeyB64 == "" {
		fmt.Fprintln(os.Stderr, "请设置环境变量 TOTP_ECC_KEY，或使用 -generate-key 生成密钥")
		os.Exit(1)
	}

	key, err := auth.ParseECCKey(eccKeyB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "密钥解析失败: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	code := auth.GenerateTOTP(key, now)
	remaining := auth.TotpPeriod - (now.Unix() % auth.TotpPeriod)
	fmt.Printf("当前动态码: %s (剩余 %d 秒)\n", code, remaining)
}
