package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kelindar/gocc"
	"github.com/kelindar/gocc/internal/config"
)

// File represents a file
type File struct {
	Name string `json:"name" uri:"name" binding:"required"`
	Body []byte `json:"body"`
}

// Ext returns the extension of the file
func (f *File) Ext() string {
	return filepath.Ext(f.Name)
}

type Result struct {
	Asm File `json:"asm"`
	Go  File `json:"go"`
}

type Request struct {
	File
	Arch    string   `uri:"arch" binding:"required"`
	Level   int      `uri:"level" form:"level"`
	Package string   `uri:"package" form:"package"`
	Options []string `uri:"options" form:"options"`
}

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Health endpoint
	r.GET("/health", func(c *gin.Context) {
		clang, _ := config.FindClang()
		objdump, _ := config.FindObjdump()
		c.JSON(http.StatusOK, gin.H{
			"compiler":     clang,
			"disassembler": objdump,
		})
	})

	// Compile endpoint
	r.POST("/compile/:arch/:name", func(c *gin.Context) {

		// Read all body
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		input := Request{
			Level: 3,
			File: File{
				Body: body,
			},
		}

		if err := c.BindUri(&input); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := c.BindQuery(&input); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		// Translate the file
		result, err := translate(&input)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, err.Error())
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}

// Translate calls the translator to translate the file and returns a result. It also
// cleans up the temporary files.
func translate(req *Request) (*Result, error) {
	arch, err := config.For(req.Arch)
	if err != nil {
		return nil, err
	}

	// Create a temporary directory for the build process
	tmp, err := os.MkdirTemp("", "*")
	if err != nil {
		return nil, err
	}

	// Clean up the temporary directory
	defer os.RemoveAll(tmp)

	// Write the input file
	input := filepath.Join(tmp, req.Name)
	if err := os.WriteFile(input, req.Body, 0644); err != nil {
		return nil, err
	}

	// Create a new translator
	options := append(req.Options, fmt.Sprintf("-O%d", req.Level))
	translator, err := gocc.NewTranslator(arch, input, tmp, req.Package, options...)
	if err != nil {
		return nil, err
	}

	// Translate the file
	if err := translator.Translate(); err != nil {
		return nil, err
	}

	// Read the generated files
	asm, err := os.ReadFile(translator.GoAssembly)
	if err != nil {
		return nil, err
	}

	goFile, err := os.ReadFile(translator.Go)
	if err != nil {
		return nil, err
	}

	// Return the result
	return &Result{
		Asm: File{
			Name: translator.GoAssembly,
			Body: asm,
		},
		Go: File{
			Name: translator.Go,
			Body: goFile,
		},
	}, nil
}
