// Copyright 2017 Jeff Foley. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package sources

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/caffix/amass/amass/internal/utils"
)

const (
	CrtshSourceString string = "crt.sh"
)

func CrtshQuery(domain, sub string, l *log.Logger) []string {
	var unique []string

	if domain != sub {
		return unique
	}

	// Pull the page that lists all certs for this domain
	url := "https://crt.sh/?q=%25." + domain
	page, err := utils.GetWebPage(url, nil)
	if err != nil {
		l.Printf("Crtsh error: %s: %v", url, err)
		return unique
	}
	// Get the subdomain name the cert was issued to, and
	// the Subject Alternative Name list from each cert
	results := crtshGetSubmatches(page)
	for _, rel := range results {
		// Do not go too fast
		time.Sleep(50 * time.Millisecond)
		// Pull the certificate web page
		url = "https://crt.sh/" + rel
		cert, err := utils.GetWebPage(url, nil)
		if err != nil {
			l.Printf("Crtsh error: %s: %v", url, err)
			continue
		}
		// Get all names off the certificate
		unique = utils.UniqueAppend(unique, crtshGetMatches(cert, domain)...)
	}
	return unique
}

func crtshGetMatches(content, domain string) []string {
	var results []string

	re := utils.SubdomainRegex(domain)
	for _, s := range re.FindAllString(content, -1) {
		results = append(results, s)
	}
	return results
}

func crtshGetSubmatches(content string) []string {
	var results []string

	re := regexp.MustCompile("<TD style=\"text-align:center\"><A href=\"([?]id=[a-zA-Z0-9]*)\">[a-zA-Z0-9]*</A></TD>")
	for _, subs := range re.FindAllStringSubmatch(content, -1) {
		results = append(results, strings.TrimSpace(subs[1]))
	}
	return results
}
