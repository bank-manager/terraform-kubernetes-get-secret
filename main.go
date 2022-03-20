package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
)

type Query struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Context   string `json:"context"`
}

type Result struct {
	Value map[string][]byte `json:"value"`
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	// Get terraform input
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal("cannot read stdin input: %v", err)
	}

	var q Query
	if err := json.Unmarshal(input, &q); err != nil {
		fatal("cannot json umarshal input: %v", err)
	}

	if q.Namespace == "" {
		fatal("missing or empty namespace parameter")
	}

	if q.Name == "" {
		fatal("missing or empty name parameter")
	}

	if q.Context == "" {
		fatal("missing or empty context parameter")
	}

	// Build kubernetes configuration
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{CurrentContext: q.Context}).ClientConfig()

	if err != nil {
		fatal("cannot load Kubernetes configuration: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fatal("cannot create Kubernetes configuration: %v", err)
	}

	secret, err := clientset.CoreV1().Secrets(q.Namespace).Get(context.TODO(), q.Name, metav1.GetOptions{})
	if err != nil {
		fatal("cannot get the %q secret in %q namespace: %v", q.Name, q.Namespace, err)
	}

	val := secret.Data
	// if !ok {
	// 	fatal("cannot found the key %q for %q secret in %q namespace", q.Key, q.Name, q.Namespace)
	// }

	// r := &Result{Value: val}

	o, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		fatal("cannot marshal indent result: %v", err)
	}

	os.Stdout.Write(o)
	return
}
