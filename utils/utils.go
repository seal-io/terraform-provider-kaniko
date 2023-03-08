package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfig(kubeconfigPath string) (*rest.Config, error) {
	// use the specified config.
	if kubeconfigPath != "" {
		return LoadConfig(kubeconfigPath)
	}

	// try the in-cluster config.
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	// try the recommended config.
	var loader = clientcmd.NewDefaultClientConfigLoadingRules()
	loader.Precedence = append(loader.Precedence,
		filepath.Join(home, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName))
	return loadConfig(loader)
}

func LoadConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		return nil, errors.New("blank kubeconfig path")
	}

	var loader = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	return loadConfig(loader)
}

func loadConfig(loader clientcmd.ClientConfigLoader) (*rest.Config, error) {
	var overrides = &clientcmd.ConfigOverrides{}
	return clientcmd.
		NewNonInteractiveDeferredLoadingClientConfig(loader, overrides).
		ClientConfig()
}

// list of default letters that can be used to make a random string when calling String
// function with no letters provided.
var defLetters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

// String generates a random string using only letters provided in the letters parameter
// if user omit letters parameters, this function will use defLetters instead.
func String(n int, letters ...string) string {
	var letterRunes []rune
	if len(letters) == 0 {
		letterRunes = defLetters
	} else {
		letterRunes = []rune(letters[0])
	}

	var bb bytes.Buffer
	bb.Grow(n)
	l := uint32(len(letterRunes))
	// on each loop, generate one random rune and append to output.
	for i := 0; i < n; i++ {
		bb.WriteRune(letterRunes[binary.BigEndian.Uint32(Bytes(4))%l])
	}
	return bb.String()
}

// Bytes generates n random bytes.
func Bytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
