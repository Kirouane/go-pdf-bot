package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeadlessRun(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	headless := NewHeadless("http://localhost:9222/json")
	pdf := headless.PrintPdf("testId", "file:///"+dir+"/config/test/chrome/test.html")
	assert.Equal(t, "testId.pdf", pdf.Filename, "the should be equal")
	assert.NotNil(t, pdf.Content, "Should not nil")
}
