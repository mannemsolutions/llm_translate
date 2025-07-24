package markdown

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Internal/Preprocessor/Main", func() {
	When("loading an MD file", func() {
		Context("with proper contents", func() {
			It("should split up properly", func() {
				md := []string{
					"# My header",
					"## My subheader",
					"My text\nwith a newline",
					"My list\n- item1\n-item2",
				}
				var myDoc = strings.Join(md, "\n\n")
				parts := NewParts(myDoc)
				Expect(parts.document).To(Equal(Part(myDoc)))
				Expect(parts.parts).To(HaveLen(len(md)))
			})
		})
	})
	When("working with a Part", func() {
		Context("using the String function", func() {
			It("should return data as expected", func() {
				for _, test := range []struct {
					part     Part
					expected string
				}{
					{part: Part("# my header"), expected: "# my header"},
					{part: Part("   "), expected: "   "},
				} {
					Expect(test.part.String()).To(Equal(test.expected))
				}
			})
		})
		Context("using the Cleansed function", func() {
			It("should return data as expected", func() {
				for _, test := range []struct {
					part     Part
					expected string
				}{
					{part: Part("# my header"), expected: "# my header"},
					{part: Part("this remains(Note: this is added by ai.)"),
						expected: "this remains"},
				} {
					Expect(test.part.Cleansed()).To(Equal(Part(test.expected)))
				}
			})
		})
		Context("using the ContainsText function", func() {
			It("should return data as expected", func() {
				for _, test := range []struct {
					part     Part
					expected bool
				}{
					{part: Part("# my header"), expected: true},
					{part: Part("## &#!@#$"), expected: false},
				} {
					Expect(test.part.ContainsText()).To(Equal(test.expected))
				}
			})
		})
		Context("using the IsHeader function", func() {
			It("should return data as expected", func() {
				for _, test := range []struct {
					part     Part
					expected bool
				}{
					{part: Part("# my header"), expected: true},
					{part: Part("whatever text"), expected: false},
					{part: Part("whatever text\nwith a linebreak"),
						expected: false},
					{part: Part("# This is broken\nbecause of the break"),
						expected: false},
				} {
					Expect(test.part.IsHeader()).To(Equal(test.expected))
				}
			})
		})
		Context("using the IsURL function", func() {
			It("should return data as expected", func() {
				for _, test := range []struct {
					part     Part
					expected bool
				}{
					{part: Part("http://whatever.com/somepath/"),
						expected: true},
					{part: Part("https://whatever.com/somepath&one=two"),
						expected: true},
					{part: Part("whatever text"), expected: false},
					{part: Part("something https://whatever.com/somepath?one=two"),
						expected: false},
				} {
					fmt.Fprintf(GinkgoWriter, "DEBUG - Test: %v", test)
					Expect(test.part.IsURL()).To(Equal(test.expected))
				}
			})
		})
		Context("using the IsPath function", func() {
			It("should return data as expected", func() {
				for _, test := range []struct {
					part     Part
					expected bool
				}{
					{part: Part("some_path"), expected: true},
					{part: Part("/some_path"), expected: true},
					{part: Part("./some_path"), expected: true},
					{part: Part("../some_path"), expected: true},
					{part: Part("~/some_path"), expected: true},
					{part: Part("~me/some_path"), expected: true},
					{part: Part("./some_path/../sub"), expected: true},
					{part: Part("../../some_path/../sub"), expected: true},
					{part: Part("../../some/folder/../sub"), expected: true},
					{part: Part("../../some/folder/../../something"), expected: true},
					{part: Part("https://whatever.com/somepath&one=two")},
					{part: Part("./.../path")},
				} {
					fmt.Fprintf(GinkgoWriter, "DEBUG - Test: %v", test)
					Expect(test.part.IsPath()).To(Equal(test.expected))
				}
			})
		})
		Context("using the WordCount function", func() {
			It("should return data as expected", func() {
				for _, test := range []struct {
					part     Part
					expected int
				}{
					{part: Part("99 bottles of beer on the wall"), expected: 7},
					{part: Part("# # #"), expected: 0},
					{part: Part("#1 # 2 #3"), expected: 3},
				} {
					fmt.Fprintf(GinkgoWriter, "DEBUG - Test: %v", test)
					Expect(test.part.WordCount()).To(Equal(test.expected))
				}
			})
		})
	})
})
