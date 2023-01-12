package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var dockerFilePath string

func main() {
	cmd := &cobra.Command{
		Use:   "kube-builder",
		Short: "Image Builder for Kubernetes",
		Run:   kubeBuilder,
	}

	cmd.Flags().StringVarP(&dockerFilePath, "docker-file", "f", "./kube-project/Dockerfile", "Path to the Dockerfile")

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		return
	}
}

func kubeBuilder(cmd *cobra.Command, args []string) {

	kubeProjectDir := "./kube-project"
	// can be passed as a cli args to cobra
	imageName := "my-kube-project"
	imageVersion := "1.0"

	cmd = exec.Command("docker", "build", "-t", imageName+":"+imageVersion, "-f", dockerFilePath, kubeProjectDir)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to build Docker image: %s\n", err)
		return
	}

	cmd = exec.Command("docker", "push", imageName+":"+imageVersion)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to push Docker image: %s\n", err)
		return
	}

	imageTag := imageName + ":" + imageVersion
	kubeDeployFile := kubeProjectDir + "/kube-deploy.yml"

	// Update deployment
	cmd = exec.Command("sed", "-i", "''", "-e", strings.Replace("s@{{IMAGE_TAG}}@"+imageTag+"@g", "/", "\\/", -1), kubeDeployFile)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to update Kubernetes Deployment configuration: %s\n", err)
		return
	}

	// Apply deployment
	cmd = exec.Command("kubectl", "apply", "-f", kubeDeployFile)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to apply updated Kubernetes Deployment configuration: %s\n", err)
		return
	}
}
