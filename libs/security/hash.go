package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	ArgonTime    = 1         // Iterasi (t)
	ArgonMemory  = 64 * 1024 // Memory usage (m)
	ArgonThreads = 4         // Paralelisme (p)
	ArgonKeyLen  = 32        // Panjang hash (k)
)

// HashSHA256 melakukan hashing menggunakan SHA-256
func HashSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt membuat salt baru
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16) // 16 byte salt
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// HashArgon2 membuat hash dari password
func HashArgon2(password string) (string, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, ArgonTime, ArgonMemory, ArgonThreads, ArgonKeyLen)

	// Encode hasil ke base64 agar mudah disimpan
	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	// Format yang sesuai standar
	fullHash := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", ArgonMemory, ArgonTime, ArgonThreads, saltBase64, hashBase64)

	return fullHash, nil
}

// VerifyHashArgon2 memverifikasi password dengan hash yang disimpan
func VerifyHashArgon2(password string, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		log.Println("Invalid hash format")
		return false
	}

	// Baca parameter dari hash
	var memory, time, threads int
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		log.Println("Failed to parse hash parameters:", err)
		return false
	}

	// Decode salt & hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		log.Println("Failed to decode salt:", err)
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		log.Println("Failed to decode hash:", err)
		return false
	}

	// Generate hash baru dengan parameter yang sama
	newHash := argon2.IDKey([]byte(password), salt, uint32(time), uint32(memory), uint8(threads), uint32(len(expectedHash)))

	// Bandingkan secara aman (mencegah timing attack)
	if subtle.ConstantTimeCompare(newHash, expectedHash) == 1 {
		return true
	}

	log.Println("Password verification failed")
	return false
}
