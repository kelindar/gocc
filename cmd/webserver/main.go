package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kelindar/gocc"
	"github.com/kelindar/gocc/internal/config"
)

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
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		input := gocc.WebRequest{
			File: gocc.File{Body: body},
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
func translate(req *gocc.WebRequest) (*gocc.WebResult, error) {
	arch, err := config.For(req.Arch)
	if err != nil {
		return nil, err
	}

	// Create a temporary directory for the build process
	tmp, err := os.MkdirTemp("", "*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)

	// Write the input file
	input := filepath.Join(tmp, req.Name)
	if err := os.WriteFile(input, req.Body, 0644); err != nil {
		return nil, err
	}

	// Create a new translator
	if len(req.Options) == 0 {
		req.Options = []string{"-O3"}
	}

	// Create a new translator
	translator, err := gocc.NewLocal(arch, input, tmp, req.Package, req.Options...)
	if err != nil {
		return nil, err
	}

	// Translate the file and read the output
	if err := translator.Translate(); err != nil {
		return nil, err
	}
	return translator.Output()
}
