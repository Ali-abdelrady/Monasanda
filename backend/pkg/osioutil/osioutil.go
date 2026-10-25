package osioutil

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"

	"golang.org/x/crypto/ssh"

	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63

type OSArchStruct struct {
	OSName        string
	OSType32Archs []string
	OSType64Archs []string
}

var osValidArchs []OSArchStruct

var lvChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var lvDivisor = uint64(len(lvChars))

func IsPathExists(parPath string) (bool, error) {

	_, pvError := os.Stat(parPath)
	if pvError == nil {
		return true, nil
	}

	if os.IsNotExist(pvError) {
		return false, nil
	}

	return true, pvError
}

func CreateDirIfNotExist(parPath string) (bool, error) {
	pfExists, pvError := IsPathExists(parPath)

	if !pfExists {
		pvError = os.Mkdir(parPath, 0777)
		if pvError != nil {
			return false, pvError
		}
	}

	return pfExists, pvError
}

func GetIDFromCurrnetTime() string {
	pvReturn := ""

	pvTime := time.Now()

	pvReturn = fmt.Sprintf("%04d", pvTime.Year())
	pvReturn += fmt.Sprintf("%02d", int(pvTime.Month()))
	pvReturn += fmt.Sprintf("%02d", pvTime.Day())
	pvReturn += fmt.Sprintf("%02d", pvTime.Hour())
	pvReturn += fmt.Sprintf("%02d", pvTime.Minute())
	pvReturn += fmt.Sprintf("%02d", pvTime.Second())
	pvReturn += fmt.Sprintf("%09d", pvTime.Nanosecond())

	return pvReturn
}

func GetOSArchs() []OSArchStruct {
	if osValidArchs == nil {
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "darwin", OSType32Archs: []string{"386"}, OSType64Archs: []string{"amd64"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "dragonfly", OSType32Archs: []string{}, OSType64Archs: []string{"amd64"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "freebsd", OSType32Archs: []string{"386", "arm"}, OSType64Archs: []string{"amd64"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "linux", OSType32Archs: []string{"386", "arm", "mips", "mipsle"}, OSType64Archs: []string{"amd64", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "nacl", OSType32Archs: []string{"386", "amd64p32", "arm"}, OSType64Archs: []string{}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "netbsd", OSType32Archs: []string{"386", "arm"}, OSType64Archs: []string{"amd64"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "openbsd", OSType32Archs: []string{"386", "arm"}, OSType64Archs: []string{"amd64"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "plan9", OSType32Archs: []string{"386", "arm"}, OSType64Archs: []string{"amd64"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "solaris", OSType32Archs: []string{}, OSType64Archs: []string{"amd64"}})
		osValidArchs = append(osValidArchs, OSArchStruct{OSName: "windows", OSType32Archs: []string{"386"}, OSType64Archs: []string{"amd64"}})
	}

	return osValidArchs
}

func IsValidOSArch(parOSName, parOSType, parOSArch string) bool {

	_ = GetOSArchs()

	parOSName = strings.ToLower(parOSName)
	parOSArch = strings.ToLower(parOSArch)

	for _, pvArch := range osValidArchs {
		if parOSName == pvArch.OSName {

			if parOSType == "64" {
				for _, pvOSType := range pvArch.OSType64Archs {
					if parOSArch == pvOSType {
						return true
					}
				}
			}

			if parOSType == "32" {
				for _, pvOSType := range pvArch.OSType32Archs {
					if parOSArch == pvOSType {
						return true
					}
				}
			}
		}
	}

	return false
}

func GetRefError(parRefID string, parError error) error {
	if parError == nil {
		return nil
	}

	return GetRefErrorTxt(parRefID, parError.Error())
}

func GetRefErrorFromResponse(parRefID string, parError interface{}) error {
	if parError == nil {
		return nil
	}

	pmError, pfOk := parError.(map[string]interface{})
	if !pfOk {
		return nil
	}

	return GetRefError(parRefID, GetRefErrorTxt(fmt.Sprint(pmError["ErrorRef"]), fmt.Sprint(pmError["ErrorMessage"])))
}

func GetRefErrorTxt(parRefID, parError string) error {

	return errors.New("[" + parRefID + "] " + parError)
}

func GetDirSize(parPath string) (int64, error) {
	var pvSize int64

	pvError := filepath.Walk(parPath, func(_ string, parInfo os.FileInfo, parError error) error {
		if parError != nil {
			return parError
		}

		if !parInfo.IsDir() {
			pvSize += parInfo.Size()
		}

		return parError
	})

	return pvSize, pvError
}

func GetUniqueID() string {
	b := make([]byte, 16)
	pvUUID := uuid.New()

	for i, pvByte := range pvUUID {
		b[i] = pvByte
	}

	pvReturn, _ := Encode16Bytes(b)
	return pvReturn
}

func Encode16Bytes(parValue []byte) (string, error) {

	if len(parValue) != 16 {
		return "", GetRefError("ENC16B0001", errors.New("invalid value length"))
	}

	pvReturn := ""

	paSegments := make([]uint64, 2)
	paSegments[0] = binary.BigEndian.Uint64(parValue[:8])
	paSegments[1] = binary.BigEndian.Uint64(parValue[8:])

	pvIndex := 0
	for _, pvSegment := range paSegments {
		for pvSegment > 0 {
			pvIndex = int(pvSegment % lvDivisor)
			pvSegment = pvSegment / lvDivisor

			pvReturn += string(lvChars[pvIndex])
		}
	}

	return pvReturn, nil
}

func GetSignature(parData []byte) []byte {
	pvHasher := md5.New()
	pvHasher.Write(parData)

	paHash := make([]byte, hex.EncodedLen(len(pvHasher.Sum(nil))))
	hex.Encode(paHash, pvHasher.Sum(nil))

	return paHash
}

func IsNumeric(parSTValue string) bool {
	_, pvError := strconv.ParseFloat(parSTValue, 64)
	return pvError == nil
}

func IsUUID(parUUID string) bool {
	_, pvError := uuid.Parse(parUUID)
	return pvError == nil
}

func IsTextAllUpper(parText string) bool {

	for _, pvRune := range parText {
		if pvRune > 127 {
			return false
		}
	}

	return !IsTextTime(parText) && strings.ToUpper(parText) == parText
}

func IsTextTime(parText string) bool {
	if strings.TrimSpace(parText) == "" {
		return false
	}

	if !(strings.Contains(parText, "-") || strings.Contains(parText, "/") || strings.Contains(parText, "\\")) {
		return false
	}

	_, pvError := dateparse.ParseAny(parText)
	return pvError == nil
}

func GetQRCode(parData string) ([]byte, error) {
	return qrcode.Encode(parData, qrcode.Medium, 256)
}

func GetQRCodeEncoded(parData string) (string, error) {
	pvPNG, pvError := GetQRCode(parData)

	if pvError != nil {
		return "", pvError
	}

	pvPMGBase64 := base64.StdEncoding.EncodeToString(pvPNG)

	return "data:image/png;base64," + pvPMGBase64, pvError
}

func GetCurrentDir() string {

	pvPath, _ := os.Getwd()

	return pvPath
}

func SSHGetConfig(parUser, parPassword string) (*ssh.ClientConfig, error) {
	// https://www.pugetsystems.com/labs/hpc/How-To-Use-SSH-Client-and-Server-on-Windows-10-1470/
	pvHomeDir, pvError := os.UserHomeDir()
	if pvError != nil {
		return nil, fmt.Errorf("[OSIO0001] %v", pvError)
	}

	pvHomeDir += "/.ssh/id_rsa"

	b, pvError := os.ReadFile(pvHomeDir)
	if pvError != nil {
		return nil, fmt.Errorf("[OSIO0002] %v", pvError)
	}

	pvSigner, pvError := ssh.ParsePrivateKey(b)
	if pvError != nil {
		return nil, fmt.Errorf("[OSIO0003] %v", pvError)
	}

	pvConfig := ssh.ClientConfig{
		User:            parUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pvSigner),
			ssh.Password(parPassword),
		},
	}

	return &pvConfig, nil
}

func Int64ToHex(parNumber int64) ([]byte, error) {
	pvBuffer := new(bytes.Buffer)
	pvError := binary.Write(pvBuffer, binary.BigEndian, parNumber)
	if pvError != nil {
		return []byte{}, pvError
	}

	return pvBuffer.Bytes(), nil
}

func ReadFolderRecursively(parFolder string, parTrigger func(parFilePath string, parIsDir bool)) error {

	paFiles, pvError := os.ReadDir(parFolder)
	if pvError != nil {
		return pvError
	}

	for _, pvFile := range paFiles {
		pvFullPath := parFolder + string(os.PathSeparator) + pvFile.Name()

		if pvFile.IsDir() {
			parTrigger(pvFullPath, true)
			pvError = ReadFolderRecursively(pvFullPath, parTrigger)
			if pvError != nil {
				return pvError
			}
		} else {
			parTrigger(pvFullPath, false)
		}
	}

	return nil
}

func GetRepoRootDir() (string, error) {
	psCmd := exec.Command("go", "env", "GOMOD")
	paOutput, pvError := psCmd.Output()
	if pvError != nil {
		return "", pvError
	}

	gomodPath := strings.TrimSpace(string(paOutput))
	if gomodPath == "" {
		return "", errors.New("go mod path not found")
	}

	repoRootBackSlash := filepath.Dir(gomodPath)
	repoRootForwardSlash := filepath.ToSlash(repoRootBackSlash)

	return repoRootForwardSlash, nil
}
