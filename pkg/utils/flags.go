/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// RegexpVal is a wrapper for regexp command line flags
type RegexpVal struct {
	regexp.Regexp
}

// Set implements the flag.Value interface
func (a *RegexpVal) Set(val string) error {
	r, err := regexp.Compile(val)
	a.Regexp = *r
	return err
}

// UnmarshalJSON implements the Unmarshaler interface from "encoding/json"
func (a *RegexpVal) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case string:
		if r, err := regexp.Compile(string(val)); err != nil {
			return err
		} else {
			*a = RegexpVal{*r}
		}
	default:
		return fmt.Errorf("invalid regexp %s", data)
	}
	return nil
}

// StringSetVal is a Value encapsulating a set of comma-separated strings
type StringSetVal map[string]struct{}

// Set implements the flag.Value interface
func (a *StringSetVal) Set(val string) error {
	m := map[string]struct{}{}
	for _, n := range strings.Split(val, ",") {
		m[n] = struct{}{}
	}
	*a = m
	return nil
}

// String implements the flag.Value interface
func (a *StringSetVal) String() string {
	if *a == nil {
		return ""
	}

	vals := make([]string, len(*a), 0)
	for val := range *a {
		vals = append(vals, val)
	}
	sort.Strings(vals)
	return strings.Join(vals, ",")
}

// StringSliceVal is a Value encapsulating a slice of comma-separated strings
type StringSliceVal []string

// Set implements the regexp.Value interface
func (a *StringSliceVal) Set(val string) error {
	*a = strings.Split(val, ",")
	return nil
}

// String implements the regexp.Value interface
func (a *StringSliceVal) String() string {
	if *a == nil {
		return ""
	}
	return strings.Join(*a, ",")
}