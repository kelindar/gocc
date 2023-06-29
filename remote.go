package gocc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/kelindar/gocc/internal/cc"
)

// WebRequest represents a request to the remote server
type WebRequest struct {
	File
	Arch    string   `uri:"arch" binding:"required"`
	Package string   `uri:"package" form:"package"`
	Options []string `uri:"options" form:"options"`
}

// File represents a file
type File struct {
	Name string `json:"name" uri:"name" binding:"required"`
	Body []byte `json:"body"`
}

// WebResult represents a result from the remote server
type WebResult struct {
	Asm File `json:"asm"`
	Go  File `json:"go"`
}

// Remote represents a remote translator
type Remote struct {
	Target  string
	Source  string
	Package string
	Output  string
	Options []string
	Client  *req.Client
}

// NewRemote creates a new translator that uses remote server
func NewRemote(target string, source, outputDir, packageName string, options ...string) (*Remote, error) {
	if packageName == "" {
		filepath.Base(outputDir)
	}

	return &Remote{
		Target:  target,
		Source:  source,
		Package: packageName,
		Options: options,
		Output:  outputDir,
		Client:  req.C().EnableDebugLog(),
	}, nil
}

// Translate translates the source file to Go assembly
func (t *Remote) Translate() error {
	if _, err := cc.Parse(t.Source); err != nil {
		return err
	}

	source, err := os.ReadFile(t.Source)
	if err != nil {
		return err
	}

	var result WebResult
	t.Client.R().
		SetBody(source).
		SetSuccessResult(&result).
		Post(t.endpoint())

	// Create output directory
	if err := os.MkdirAll(t.Output, 0755); err != nil {
		return err
	}

	// Write the stub file
	if err := os.WriteFile(outputPath(t.Output, result.Go.Name), result.Go.Body, 0644); err != nil {
		return err
	}

	// Write the assembly file
	if err := os.WriteFile(outputPath(t.Output, result.Asm.Name), result.Asm.Body, 0644); err != nil {
		return err
	}

	return nil
}

// endpoint returns the endpoint for the translator
func (t *Remote) endpoint() string {
	uri := fmt.Sprintf("https://gocc.onrender.com/compile/%s/%s", t.Target, t.Source)
	args := make([]string, 0, len(t.Options)+1)
	if t.Package != "" {
		args = append(args, fmt.Sprintf("package=%s", t.Package))
	}
	if len(t.Options) > 0 {
		for _, option := range t.Options {
			args = append(args, fmt.Sprintf("options=%s", option))
		}
		// args = append(args, fmt.Sprintf("options=%s", strings.Join(t.Options, ",")))
	}

	if len(args) > 0 {
		uri += "?" + strings.Join(args, "&")
	}
	return uri
}

// outputPath returns the output path
func outputPath(output, path string) string {
	return filepath.Join(output, filepath.Base(path))
}
