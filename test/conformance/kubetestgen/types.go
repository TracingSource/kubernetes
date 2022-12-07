package main

// Area is a conformance area composed of a list of test suites
type Area struct {
	Area   string  `json:"area,omitempty"`
	Suites []Suite `json:"suites,omitempty"`
}

// Suite is a conformance test suite composed of a list of behaviors
type Suite struct {
	Suite       string     `json:"suite,omitempty"`
	Description string     `json:"description,omitempty"`
	Behaviors   []Behavior `json:"behaviors,omitempty"`
}

// Behavior describes the set of properties for a conformance behavior
type Behavior struct {
	ID          string `json:"id,omitempty"`
	APIObject   string `json:"apiObject,omitempty"`
	APIField    string `json:"apiField,omitempty"`
	APIType     string `json:"apiType,omitempty"`
	Description string `json:"description,omitempty"`
}
