package main

import (
    "bufio"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "os"
    "strings"
)

func main() {
    key := make([]byte, 32)
    _, err := rand.Read(key)
    if err != nil {
        fmt.Println("Error generating key:", err)
        os.Exit(1)
    }

    encodedKey := base64.URLEncoding.EncodeToString(key)

    file, err := os.Open(".env")
    if err != nil && !os.IsNotExist(err) {
        fmt.Println("Error opening .env file:", err)
        os.Exit(1)
    }
    defer file.Close()
    var lines []string
    updated := false
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "JWT_SECRET_KEY=") {
            line = fmt.Sprintf("JWT_SECRET_KEY=%s", encodedKey)
            updated = true
        }
        lines = append(lines, line)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading .env file:", err)
        os.Exit(1)
    }
    if !updated {
        lines = append(lines, fmt.Sprintf("JWT_SECRET_KEY=%s", encodedKey))
    }

    file, err = os.Create(".env")
    if err != nil {
        fmt.Println("Error creating .env file:", err)
        os.Exit(1)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, line := range lines {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            fmt.Println("Error writing to .env file:", err)
            os.Exit(1)
        }
    }
    writer.Flush()

    fmt.Println("JWT secret key updated in .env file.")
}
