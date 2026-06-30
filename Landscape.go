package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var albero int
var sentiero string

func main() {
	fmt.Println("--------------- LandScape Encryptor ---------------")

	fmt.Print("\nInserisci la PWD Iniziale: ")
	fmt.Scan(&sentiero)
	fmt.Print("\nInserisci il numero intermedio dell' albero: ")
	fmt.Scan(&albero)

	fmt.Println("Landscape Fingerprint: " + FingerPrintGenerator())
	fmt.Println()
	credentialsFaker()
	fmt.Println("Landscape Final Payload: " + finalSeed())
	var input string
	fmt.Scan(&input)
}

func tokenHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func tokenBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func randomRange(min, max int) int {
	if max < min {
		return min
	}
	bg := big.NewInt(int64(max - min + 1))
	n, _ := rand.Int(rand.Reader, bg)
	return min + int(n.Int64())
}

func FingerPrintGenerator() string {
	casa, _ := os.Hostname()
	datiMacchina := "1234567890" + runtime.GOOS + casa + runtime.GOARCH

	extra := tokenHex(8)

	possibili := []int{albero - 1, albero, albero + 1}
	scelto := possibili[randomRange(0, len(possibili)-1)]

	risultato := []byte(datiMacchina + extra + strconv.Itoa(scelto))

	for i := 0; i < 100; i++ {
		h := sha256.Sum256(risultato)
		risultato = h[:]
	}

	return fmt.Sprintf("%x", sha256.Sum256(risultato))
}

func credentialsFaker() {
	locales := []string{"en_US", "es_ES", "fr_FR", "de_DE", "it_IT", "ja_JP", "ko_KR", "zh_CN", "ru_RU"}

	emailProviders := map[string][]string{
		"en_US": {"gmail.com", "yahoo.com", "outlook.com"},
		"it_IT": {"gmail.com", "libero.it", "outlook.it", "aruba.it"},
	}

	listaComune := []string{"123456", "password", "123456789", "qwerty", "admin", "letmein", "root"}

	nomiFake := []string{"JohnDoe", "MarioRossi", "JeanDupont", "CarlosGomez"}
	locale := locales[randomRange(0, len(locales)-1)]
	nome := nomiFake[randomRange(0, len(nomiFake)-1)]

	bigProb, _ := rand.Int(rand.Reader, big.NewInt(100))
	if bigProb.Int64() < 10 && len(nome) > 2 {
		caratteri := []string{".", "_", "-"}
		car := caratteri[randomRange(0, len(caratteri)-1)]
		luogo := randomRange(1, len(nome)-1)
		nome = nome[:luogo] + car + nome[luogo:]
	}

	providers := emailProviders[locale]
	if len(providers) == 0 {
		providers = emailProviders["en_US"]
	}
	provider := providers[randomRange(0, len(providers)-1)]
	email := strings.ToLower(nome) + "@" + provider

	var password string
	scelta, _ := rand.Int(rand.Reader, big.NewInt(100))

	if scelta.Int64() < 40 {
		password = strings.ToLower(nome) + strconv.Itoa(randomRange(39, 99))
	} else if scelta.Int64() < 50 {
		caratteri := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
		length := randomRange(6, 12)
		b := make([]byte, length)
		for i := range b {
			b[i] = caratteri[randomRange(0, len(caratteri)-1)]
		}
		password = string(b)
	} else {
		if randomRange(1, 2) == 1 {
			password = listaComune[randomRange(0, len(listaComune)-1)]
		} else {
			var pwSenzaNumero []string
			for _, pwd := range listaComune {
				if !strings.ContainsAny(pwd, "0123456789") {
					pwSenzaNumero = append(pwSenzaNumero, pwd)
				}
			}
			possibili := []int{albero - 1, albero, albero + 1}
			scelto := possibili[randomRange(0, len(possibili)-1)]
			if len(pwSenzaNumero) > 0 {
				password = pwSenzaNumero[randomRange(0, len(pwSenzaNumero)-1)] + strconv.Itoa(scelto)
			} else {
				password = listaComune[randomRange(0, len(listaComune)-1)] + strconv.Itoa(scelto)
			}
		}
	}

	fmt.Printf("Landscape Credentials  Email: %s - Password: %s\n", email, password)
}

func pbkdf2HmacSha256(password, salt []byte, iter, keyLen int) []byte {
	prf := hmac.New(sha256.New, password)
	hash := make([]byte, 0, keyLen)
	numBlocks := (keyLen + sha256.Size - 1) / sha256.Size

	for block := 1; block <= numBlocks; block++ {
		prf.Reset()
		_, _ = prf.Write(salt)
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(block))
		_, _ = prf.Write(buf)
		dk := prf.Sum(nil)
		u := dk

		for i := 1; i < iter; i++ {
			prf.Reset()
			_, _ = prf.Write(u)
			u = prf.Sum(nil)
			for j := range dk {
				dk[j] ^= u[j]
			}
		}
		hash = append(hash, dk...)
	}
	return hash[:keyLen]
}

func pbkdf2HmacSha512(password, salt []byte, iter, keyLen int) []byte {
	prf := hmac.New(sha512.New, password)
	hash := make([]byte, 0, keyLen)
	blocchi := (keyLen + sha512.Size - 1) / sha512.Size

	for block := 1; block <= blocchi; block++ {
		prf.Reset()
		_, _ = prf.Write(salt)
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(block))
		_, _ = prf.Write(buf)
		dk := prf.Sum(nil)
		u := dk

		for i := 1; i < iter; i++ {
			prf.Reset()
			_, _ = prf.Write(u)
			u = prf.Sum(nil)
			for j := range dk {
				dk[j] ^= u[j]
			}
		}
		hash = append(hash, dk...)
	}
	return hash[:keyLen]
}

func layer1(pwd []byte) []byte {
	salt := make([]byte, 8)
	binary.BigEndian.PutUint64(salt, uint64(time.Now().UnixNano()))
	return pbkdf2HmacSha256(pwd, salt, 20000, 32)
}

func layer2(data []byte) []byte {
	fp := []byte(runtime.GOOS + "release-sim" + runtime.GOARCH)
	h := sha256.New()
	h.Write(append(data, fp...))
	return h.Sum(nil)
}

func layer3(data []byte) []byte {
	noise := tokenBytes(32)
	h := sha512.New()
	h.Write(append(data, noise...))
	return h.Sum(nil)
}

func layer4(data []byte) []byte {
	t1 := make([]byte, 8)
	binary.BigEndian.PutUint64(t1, uint64(time.Now().UnixNano()))
	t2 := make([]byte, 8)
	binary.BigEndian.PutUint64(t2, uint64(time.Now().Unix()*1e6))

	h := sha256.New()
	h.Write(append(data, append(t1, t2...)...))
	return h.Sum(nil)
}

func layer5(data []byte) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[i] = ((b << 1) | (b >> 7)) & 0xFF
	}
	return res
}

func layer6(data []byte) []byte {
	h512 := sha512.Sum512(data)
	h256 := sha256.Sum256(h512[:])
	return h256[:]
}

func layer7(data []byte) []byte {
	key := sha256.Sum256(data)
	res := make([]byte, len(data))
	for i, b := range data {
		res[i] = b ^ key[i%len(key)]
	}
	return res
}

func layer8(data []byte) []byte {
	h := sha512.New()
	h.Write(append(data, tokenBytes(64)...))
	return h.Sum(nil)
}

func layer9(data []byte) []byte {
	block := 4
	var chunks [][]byte
	for i := 0; i < len(data); i += block {
		end := i + block
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	for i, j := 0, len(chunks)-1; i < j; i, j = i+1, j-1 {
		chunks[i], chunks[j] = chunks[j], chunks[i]
	}
	var res []byte
	for _, chunk := range chunks {
		res = append(res, chunk...)
	}
	return res
}

func layer10(data []byte) []byte {
	h256 := sha256.Sum256(data)
	h := sha256.New()
	h.Write(append(data, h256[:]...))
	return h.Sum(nil)
}

func layer11(data []byte) []byte {
	saltLen := 16
	if len(data) < saltLen {
		saltLen = len(data)
	}
	return pbkdf2HmacSha512(data, data[:saltLen], 5000, 64)
}

func layer12(data []byte) []byte {
	h256 := sha256.Sum256(data)
	h := sha512.New()
	h.Write(append(data, h256[:]...))
	return h.Sum(nil)
}

func layer13(data []byte) []byte {
	res := make([]byte, len(data))
	l := len(data)
	for i := 0; i < l; i++ {
		a := data[i]
		b := data[l-1-i]
		res[i] = (a ^ (b + byte(i))) & 0xFF
	}
	return res
}

func layer14(data []byte) []byte {
	tStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha256.New()
	h.Write(append(data, append([]byte(tStr), tokenBytes(16)...)...))
	return h.Sum(nil)
}

func layer15(data []byte) []byte {
	mid := len(data) / 2
	left := sha256.Sum256(data[:mid])
	right := sha256.Sum256(data[mid:])
	h := sha512.New()
	h.Write(append(left[:], right[:]...))
	return h.Sum(nil)
}

func layer16(data []byte) []byte {
	h256 := sha256.Sum256(data)
	seed := binary.BigEndian.Uint64(h256[:8])

	lst := make([]byte, len(data))
	copy(lst, data)

	l := uint64(len(lst))

	for i := range lst {
		j := int((seed + uint64(i)*31) % l)
		lst[i], lst[j] = lst[j], lst[i]
	}

	return lst
}

func layer17(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha512.Sum512(data)
	res := make([]byte, 32)
	for i := 0; i < 32; i++ {
		res[i] = h1[i] ^ h2[i]
	}
	return res
}

func layer18(data []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(buf, data)
	return buf
}

func layer19(data []byte) []byte {
	h := sha512.New()
	h.Write(append(data, tokenBytes(32)...))
	return h.Sum(nil)
}

func layer20(data []byte) string {
	h := sha256.Sum256(data)
	return strings.ToUpper(fmt.Sprintf("%x", h))
}

func finalSeed() string {
	var rounds int

	fmt.Print("Metti il numero di rounds: ")
	fmt.Scan(&rounds)

	data := layer1([]byte(sentiero))

	for i := 0; i < rounds; i++ {
		data = layer2(data)
		data = layer3(data)
		data = layer4(data)
		data = layer5(data)
		data = layer6(data)
		data = layer7(data)
		data = layer8(data)
		data = layer9(data)
		data = layer10(data)
		data = layer11(data)
		data = layer12(data)
		data = layer13(data)
		data = layer14(data)
		data = layer15(data)
		data = layer16(data)
		data = layer17(data)
		data = layer18(data)
		data = layer19(data)
	}

	return layer20(data)
}
