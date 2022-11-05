/*
Copyright (c) 2022, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package pandoc_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	// Port defaults to 3030, it is the port number that pandoc-server listens on
	Port string `json:"port,omitempty"`
	// From is the doc type you are converting from, e.g. markdown
	From string `json:"from,omitempty"`
	// To is the doc type you are converting to, e.g. html5
	To string `json:"to,omitempty"`
	//
	// For the following fields see https://pandoc.org/pandoc-server.html#root-endpoint
	//
	ShiftHeadingLevel     int                    `json:"shift-heading-level-by,omitempty"`
	IdentedCodeClasses    []string               `json:"indented-code-classes,omitempty"`
	DefaultImageExtension string                 `json:"default-image-extension,omitempty"`
	Metadata              string                 `json:"metadata,omitempty"`
	TabStop               int                    `json:"tab-stop,omitempty"`
	TrackChanges          string                 `json:"track-changes,omitempty"`
	Abbreviations         []string               `json:"abbreviations,omitempty"`
	Standalone            bool                   `json:"standalone,omitempty"`
	Text                  string                 `json:"text,omitempty"`
	Template              string                 `json:"template,omitempty"`
	Variables             map[string]interface{} `json:"variables,omitempty"`
	DPI                   int                    `json:"dpi,omitemtpy"`
	Wrap                  string                 `json:"wrap,omitempty"`
	Columns               int                    `json:"columns,omitempty"`
	TableOfContents       bool                   `json:"table-of-contents,omitempty"`
	TOCDepth              int                    `json:"toc-depth,omitempty"`
	StripComments         bool                   `json:"strip-comments,omitempty"`
	HighlightStyle        string                 `json:"highlight-style,omitempty"`
	EmbedResources        string                 `json:"embed-resources,omitempty"`
	HTMLQTags             bool                   `json:"html-q-tags,omitempty"`
	Ascii                 bool                   `json:"ascii,omitempty"`
	ReferenceLinks        bool                   `json:"reference-links,omitempty"`
	ReferenceLocation     string                 `json:"reference-location,omitempty"`
	SetExtHeaders         string                 `json:"setext-headers,omitempty"`
	TopLevelDivision      string                 `json:"top-level-division,omitempty"`
	NumberSections        string                 `json:"number-sections,omitempty"`
	NumberOffset          []int                  `json:"number-offset,omitempty"`
	HTMLMathMethod        string                 `json:"html-math-method,omitempty"`
	Listings              bool                   `json:"listings,omitempty"`
	Incremental           bool                   `json:"incremental,omitempty"`
	SideLevel             int                    `json:"slide-level,omitempty"`
	SectionDivs           bool                   `json:"section-divs,omitempty"`
	EmailObfuscation      string                 `json:"email-obfuscation,omitempty"`
	IdentifierPrefix      string                 `json:"identifier-prefix,omitempty"`
	TitlePrefix           string                 `json:"title-prefix,omitempty"`
	ReferenceDoc          string                 `json:"reference-doc,omitempty"`
	EPubCoverImage        string                 `json:"epub-cover-image,omitempty"`
	EPubMetadata          string                 `json:"epub-metadata,omitempty"`
	EPubChapterLevel      int                    `json:"epub-chapter-level,omitempty"`
	EPubSubdirectory      string                 `json:"epub-subdirectory,omitempty"`
	EPubFonts             string                 `json:"epub-fonts,omitempty"`
	IpynbOutput           string                 `json:"ipynb-output,omitempty"`
	Citeproc              bool                   `json:"citeproc,omitempty"`
	Bibliography          []string               `json:"bibliography,omitempty"`
	Csl                   string                 `json:"csl,omitempty"`
	CiteMethod            string                 `json:"cite-method,omitempty"`
	Files                 []string               `json:files,omitempty"`

	// Verbose if set true then include logging on success as well as error
	Verbose bool

	// ExtTypes holds a mapping of extension to file type, e.d. ".html" to "html5"
	//ExtTypes map[string]string `json:"ext-types,omitempty"`
}

var (
	// DefaultExtTypes maps file extensions to document types. This allows the "to", "from"
	// Pandoc options to be set based on file extension. This can be overwritten by setting
	// `.ext_types` in the JSON configuraiton file.
	DefaultExtTypes = map[string]string{
		".md":   "markdown",
		".html": "html5",
	}
)

func inStringList(val string, list []string) bool {
	for _, expected := range list {
		if val == expected {
			return true
		}
	}
	return false
}

// Load will read a JSON file containing config attributes
// and return a config struct and error.
func Load(fName string) (*Config, error) {
	src, err := os.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	if err := json.Unmarshal(src, cfg); err != nil {
		return nil, err
	}
	if cfg.Port == "" {
		cfg.Port = ":3030"
	} else if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = fmt.Sprintf(":%s", cfg.Port)
	}

	if !inStringList(cfg.TrackChanges, []string{"accept", "reject", "all", ""}) {
		return cfg, fmt.Errorf("tract-changes: %q is not supported", cfg.TrackChanges)
	}
	if !inStringList(cfg.Wrap, []string{"auto", "preserve", "none", ""}) {
		return cfg, fmt.Errorf("wrap: %q is not supported", cfg.Wrap)
	}
	if !inStringList(cfg.HighlightStyle, []string{"pygments", "kate", "monochrome", "breezeDark", "espresso", "zenburn", "haddock", "tango", ""}) {
		return cfg, fmt.Errorf("highlight-style: %q is not supported", cfg.HighlightStyle)
	}
	if !inStringList(cfg.ReferenceLocation, []string{"document", "section", "block", ""}) {
		return cfg, fmt.Errorf("wrap: %q is not supported", cfg.ReferenceLocation)
	}
	if !inStringList(cfg.TopLevelDivision, []string{"default", "part", "chapter", "section", ""}) {
		return cfg, fmt.Errorf("top-level-division: %q is not supported", cfg.TopLevelDivision)
	}
	if !inStringList(cfg.HTMLMathMethod, []string{"plain", "webtex", "gladtex", "mathml", "mathjax", "katex", ""}) {
		return cfg, fmt.Errorf("html-math-method: %q is not supported", cfg.HTMLMathMethod)
	}
	if !inStringList(cfg.EmailObfuscation, []string{"none", "references", "javascript", ""}) {
		return cfg, fmt.Errorf("email-obfuscation: %q is not supported", cfg.EmailObfuscation)
	}
	if !inStringList(cfg.IpynbOutput, []string{"best", "all", "none", ""}) {
		return cfg, fmt.Errorf("ipynb-output: %q is not supported", cfg.IpynbOutput)
	}
	if !inStringList(cfg.CiteMethod, []string{"citeproc", "natbib", "biblatex", ""}) {
		return cfg, fmt.Errorf("cite-method: %q is not supported", cfg.CiteMethod)
	}
	/*
		// Make sure we have the extension mappings for document type
		if cfg.ExtTypes == nil {
			cfg.ExtTypes = map[string]string{}
		}
		for k, v := range DefaultExtTypes {
			if _, ok := cfg.ExtTypes[k]; !ok {
				cfg.ExtTypes[k] = v
			}
		}
	*/
	return cfg, nil
}

// RootEndpoint takes content type and sends the request to the Pandoc Server
// Root end point based on the state of configuration struct used.
func (cfg *Config) RootEndpoint(contentType string) ([]byte, error) {
	// NOTE: Pandoc Server API want JSON in POST not urlencoded form data
	if cfg.Text == "" {
		return nil, fmt.Errorf("expected to have a source text to convert, %+v", cfg)
	}
	src, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if len(src) == 0 {
		log.Printf("Nothing to convert")
		return nil, fmt.Errorf("nothing to convert")
	}
	// Setup out our JSON post request.
	u := fmt.Sprintf("http://localhost%s/", cfg.Port)
	body := bytes.NewReader(src)
	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%s POST failed, %s", u, err)
		return nil, err
	}
	defer resp.Body.Close()
	// Process response
	src, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s POST read body failed, %s", u, err)
		return nil, err
	}
	if len(src) == 0 {
		log.Printf("zero bytes returned from Root Endpoint")
		return nil, fmt.Errorf("zero bytes returned by pandoc")
	}
	if cfg.Verbose {
		log.Printf("%d bytes returned successful from Root Endpoint", len(src))
	}
	return src, nil
}

// Pandoc a takes the configuration settings and sends a request
// to the Pandoc server with contents read from the io.Reader
// and returns a slice of bytes and error.
//
// ```
//
//	 // Setup our client configuration
//		cfg := pandoc_client.Config{
//			Standalone: true,
//			From: "markdown",
//			To: "html5",
//		}
//		src, err := os.ReadFile("htdocs/index.md")
//		// ... handle error
//		txt, err :=  cfg.Convert(bytes.NewReader(src), "text/plain"))
//		if err := os.WriteFile("htdocs/index.html", src, 0664); err != nil {
//		    // ... handle error
//		}
//
// ```
func (cfg *Config) Convert(input io.Reader, contentType string) ([]byte, error) {
	// Check to make sure we have one of the supported mimetypes
	if !inStringList(contentType, []string{"text/plain", "application/json", "application/octet-stream"}) {
		return nil, fmt.Errorf("Only text/plain, application/json, application/octet-stream mimetypes are supported.")
	}
	var src []byte

	src, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}
	// See base64 encoding is needed.
	if contentType == "application/octet-stream" {
		//FIXME: confirm that we have a base 64 encoding binary stream
	}
	cfg.Text = fmt.Sprintf("%s", src)
	defer func() {
		cfg.Text = ""
	}()
	src, err = cfg.RootEndpoint(contentType)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// Walk takes a path and walks the directories converting the files that map
// to the From values in the configuration.
func (cfg *Config) Walk(startPath string, fromExt string, toExt string) error {
	return fmt.Errorf("Walk(startPath string) error not implemented.")
}
