package markdown_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPreprocessor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Preprocessor Suite")
}
