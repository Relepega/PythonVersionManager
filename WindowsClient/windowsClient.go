package windowsClient

import (
	pythonVersion "PythonVersion"
	utils "Utils"

	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const SymlinkDest = "%localappdata%\\Python"
const PythonRootContainer = ".\\Python"

type PythonVersions struct {
	All          []string
	Stable       []string
	Unstable     []string
	Classes      map[string]*pythonVersion.PythonVersion
	CreationDate int64
}

type Client struct {
	AppRoot          string
	Arch             string
	HttpClient       http.Client
	PythonVersions   PythonVersions
	CachedDataExists bool
	ExpiryDate       int64
}

func NewClient() *Client {
	appRoot := utils.GetWorkingDir()

	arch := "amd64"

	if runtime.GOARCH == "386" {
		arch = "win32"
	}

	httpClient := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Transport: &http.Transport{
			IdleConnTimeout: 60 * time.Second,
		},
	}

	return &Client{
		AppRoot:    appRoot,
		Arch:       arch,
		HttpClient: httpClient,
		PythonVersions: PythonVersions{
			All:          []string{},
			Stable:       []string{},
			Unstable:     []string{},
			Classes:      make(map[string]*pythonVersion.PythonVersion),
			CreationDate: 0,
		},
		CachedDataExists: false,
		ExpiryDate:       0,
	}
}

func (client *Client) clientInfo() string {
	return fmt.Sprintf("Windows client (%s)", client.Arch)
}

func (client *Client) MakeSymlink(version string, srcPath string) bool {
	// check if path exists and is a directory
	fileInfo, err := os.Stat(srcPath)

	if err != nil && !fileInfo.IsDir() {
		log.Fatalf(`Python "%s" is not installed. Try installing it first...`, version)
		return false
	}

	// run command
	symlinkCommand := fmt.Sprintf("New-Item -Force -ItemType SymbolicLink -Path '%s' -Target '%s'", SymlinkDest, srcPath)
	command := []string{"powershell.exe", "-noprofile", `Start-Process -WindowStyle hidden -Verb RunAs -Wait powershell.exe -Args "` + symlinkCommand + `"`}

	fmt.Print("Making symlink... ")

	_, err = utils.RunCmd(strings.Join(command, " "))

	if err != nil {
		print("Couldn't create the symlink, exiting ...")
		log.Fatal(err)
	}

	fmt.Println("Done!")

	return true
}
