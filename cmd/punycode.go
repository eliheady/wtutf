package cmd

import "golang.org/x/net/idna"

// toPuny takes a string and a slice of []idna.Option rules and calls
// idna.ToASCII, returning the punycode string and error
func toPuny(s string, rules []idna.Option) (string, error) {
	punyRules := idna.New(
		rules...,
	)
	return punyRules.ToASCII(s)
}

// fromPuny takes a punycode string and returns the decoded UTF-8 string
func fromPuny(s string, rules []idna.Option) (string, error) {
	punyRules := idna.New(
		rules...,
	)
	return punyRules.ToUnicode(s)
}

// canPunyConvert takes a rune and a slice of []idna.Option rules and attempts
// the ToASCII conversion, then returns a bool indicating success
func canPunyConvert(s string, rules []idna.Option) bool {
	_, err := toPuny(s, rules)
	return err == nil
}

// enumerateErrors takes a rune, checks several punycode conversion rules and
// reports the failures as a single string
func enumerateErrors(r rune) []string {
	rules := map[string][]idna.Option{
		"CheckBidi (RFC 5893)":                       {idna.BidiRule()},
		"CheckJoiners (RFC 5892)":                    {idna.CheckJoiners(true)},
		"CheckHyphens (UTS 46)":                      {idna.CheckHyphens(true)},
		"ValidateForRegistration (RFC 5891)":         {idna.ValidateForRegistration()},
		"ValidateLabels (RFC 5891)":                  {idna.ValidateLabels(true)},
		"UseSTD3ASCIIRules (RFC 1034, 5891, UTS 46)": {idna.StrictDomainName(true), idna.ValidateLabels(true)},
	}

	var allErrors []string

	for i, ruleset := range rules {
		if !canPunyConvert(string(r), ruleset) {
			allErrors = append(allErrors, i)
		}
	}
	return allErrors
}
