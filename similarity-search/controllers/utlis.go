package controllers

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

const (
	TYPE_PDF  = "application/pdf"
	TYPE_JPEG = "image/jpeg"
	TYPE_PNG  = "image/jpg"
)

func isPDF(content []byte) bool {
	return http.DetectContentType(content) == "application/pdf"
}
func isImage(content []byte) bool {
	contentType := http.DetectContentType(content)
	return contentType == "image/png" || contentType == "image/jpeg"
}
func getContentType(content []byte) string {
	return http.DetectContentType(content)
}

func saveDoc(content []byte) {
	file, err := os.Create("test.pdf")
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return
	}
	_, err = file.Write(content)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return
	}
}
func cleanContent(word string) string {
	patterns := []string{
		`[.]{2,}\s*\d+`, // Matches at least 10 consecutive dots followed by a number
		`[-\s]{3,}`,     // Matches at least 3 consecutive dashes or spaces
		`[_\s]{3,}`,     // Matches at least 3 consecutive underscore or spaces
		`(?m)^[\s]*[(\[]?[\s]*[0-9]*[\s]*[\])][\s]*[.]?`, // Matches numbered list ex: 12), (12), 12). 12)
		// ^\s*\(?[0-9]+\)?\.?$
		`(?m)^[\s]*[0-9]+[.]`,                     //Matches numbered list ex: 12.
		`(?m)^\s*[(\[]?\s*[a-zA-Z]\s*[)\]]\s*\.?`, // Matches single alphabet list ex:  a)., (a), (a).
		`(?m)^\s*[a-zA-Z]\s*\.`,
		`\(\s*[A-Za-z\s–-]+\s*-\s*[A-Za-z\s–-]+\s*\)`, // Matches text within round brackets with a dash (A - B)
		`(?i)\(\x60*\s?in\s?Lakhs\)`,                  // Matches (` in Lakhs)"
		`\x60*`,                                       //matches one or more backticks

		`\(\s*\)`, // Matches empty parentheses
		`/\s*-`,   //matches text like /-, / -
		`(?m)^\s*[(\[]?[IVXLCDMivxlcdm]+[)\]]\.?`, //matches the ROMAN number enclosed by/without parantheses.
		`(?m)^\s*[IVXLCDMivxlcdm]+\s*\.`,          //matches the ROMAN number enclosed by/without parantheses.
	}

	for _, pattern := range patterns {
		reg := regexp.MustCompile(pattern)
		word = reg.ReplaceAllString(word, "")
	}
	return word
}
