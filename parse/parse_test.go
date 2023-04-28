package parse

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buf(s string) bytes.Buffer {
	return *bytes.NewBuffer([]byte(s))
}

var tests = []struct {
	Text string
	Node Node
}{

	//
	// text only
	//
	{
		Text: "text",
		Node: &TextNode{Value: "text"},
	},
	{
		Text: "}text",
		Node: &TextNode{Value: "}text"},
	},
	{
		Text: "http://github.com",
		Node: &TextNode{Value: "http://github.com"}, // should not escape double slash
	},
	{
		Text: "$${string}",
		Node: &ListNode{Nodes: []Node{&TextNode{Value: "$"}, &FuncNode{
			Param: "string",
			buf: buf("${string}"),
		}}},
	},
	{
		Text: "$$string",
		Node: &ListNode{Nodes: []Node{&TextNode{Value: "$"}, &FuncNode{
			Param: "string",
			buf: buf("$string"),
		}}},
	},
	{
		Text: `\\.\pipe\pipename`,
		Node: &TextNode{Value: `\\.\pipe\pipename`},
	},

	//
	// variable only
	//
	{
		Text: "${string}",
		Node: &FuncNode{Param: "string", buf: buf("${string}")},
	},

	//
	// text transform functions
	//
	{
		Text: "${string,}",
		Node: &FuncNode{
			Param: "string",
			Name:  ",",
			Args:  nil,
			buf: buf("${string,}"),
		},
	},
	{
		Text: "${string,,}",
		Node: &FuncNode{
			Param: "string",
			Name:  ",,",
			Args:  nil,
			buf: buf("${string,,}"),
		},
	},
	{
		Text: "${string^}",
		Node: &FuncNode{
			Param: "string",
			Name:  "^",
			Args:  nil,
			buf: buf("${string^}"),
		},
	},
	{
		Text: "${string^^}",
		Node: &FuncNode{
			Param: "string",
			Name:  "^^",
			Args:  nil,
			buf: buf("${string^^}"),
		},
	},

	//
	// substring functions
	//
	{
		Text: "${string:position}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&TextNode{Value: "position"},
			},
			buf: buf("${string:position}"),
		},
	},
	{
		Text: "${string:position:length}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&TextNode{Value: "position"},
				&TextNode{Value: "length"},
			},
			buf: buf("${string:position:length}"),
		},
	},

	//
	// string removal functions
	//
	{
		Text: "${string#substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
			buf: buf("${string#substring}"),
		},
	},
	{
		Text: "${string##substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "##",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
			buf: buf("${string##substring}"),
		},
	},
	{
		Text: "${string%substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "%",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
			buf: buf("${string%substring}"),
		},
	},
	{
		Text: "${string%%substring}",
		Node: &FuncNode{
			Param: "string",
			Name:  "%%",
			Args: []Node{
				&TextNode{Value: "substring"},
			},
			buf: buf("${string%%substring}"),
		},
	},

	//
	// string replace functions
	//
	{
		Text: "${string/substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
			buf: buf("${string/substring/replacement}"),
		},
	},
	{
		Text: "${string//substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "//",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
			buf: buf("${string//substring/replacement}"),
		},
	},
	{
		Text: "${string/#substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "/#",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
			buf: buf("${string/#substring/replacement}"),
		},
	},
	{
		Text: "${string/%substring/replacement}",
		Node: &FuncNode{
			Param: "string",
			Name:  "/%",
			Args: []Node{
				&TextNode{Value: "substring"},
				&TextNode{Value: "replacement"},
			},
			buf: buf("${string/%substring/replacement}"),
		},
	},

	//
	// default value functions
	//
	{
		Text: "${string=default}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&TextNode{Value: "default"},
			},
			buf: buf("${string=default}"),
		},
	},
	{
		Text: "${string:=default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":=",
			Args: []Node{
				&TextNode{Value: "default"},
			},
			buf: buf("${string:=default}"),
		},
	},
	{
		Text: "${string:-default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":-",
			Args: []Node{
				&TextNode{Value: "default"},
			},
			buf: buf("${string:-default}"),
		},
	},
	{
		Text: "${string:?default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":?",
			Args: []Node{
				&TextNode{Value: "default"},
			},
			buf: buf("${string:?default}"),
		},
	},
	{
		Text: "${string:+default}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":+",
			Args: []Node{
				&TextNode{Value: "default"},
			},
			buf: buf("${string:+default}"),
		},
	},

	//
	// length function
	//
	{
		Text: "${#string}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
			buf: buf("${#string}"),
		},
	},

	//
	// special characters in argument
	//
	{
		Text: "${string#$%:*{}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
			Args: []Node{
				&TextNode{Value: "$%:*{"},
			},
			buf: buf("${string#$%:*{}"),
		},
	},

	// text before and after function
	{
		Text: "hello ${#string} world",
		Node: &ListNode{
			Nodes: []Node{
				&TextNode{
					Value: "hello ",
				},
				&ListNode{
					Nodes: []Node{
						&FuncNode{
							Param: "string",
							Name:  "#",
							buf: buf("${#string}"),
						},
						&TextNode{
							Value: " world",
						},
					},
				},
			},
		},
	},
	// text before and after function with \\ outside of function
	{
		Text: `\\ hello ${#string} world \\`,
		Node: &ListNode{
			Nodes: []Node{
				&TextNode{
					Value: `\\ hello `,
				},
				&ListNode{
					Nodes: []Node{
						&FuncNode{
							Param: "string",
							Name:  "#",
							buf: buf("${#string}"),
						},
						&TextNode{
							Value: ` world \\`,
						},
					},
				},
			},
		},
	},

	// TODO
	// Tests from here down are broken

	// escaped function arguments
	{
		Text: `${string/\/position/length}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "/position",
				},
				&TextNode{
					Value: "length",
				},
			},
			buf: buf(`${string/\/position/length}`),
		},
	},
	{
		Text: `${string/\/position\\/length}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "/position\\",
				},
				&TextNode{
					Value: "length",
				},
			},
		},
	},
	{
		Text: `${string/position/\/length}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "position",
				},
				&TextNode{
					Value: "/length",
				},
			},
		},
	},
	{
		Text: `${string/position/\/length\\}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "position",
				},
				&TextNode{
					Value: "/length\\",
				},
			},
		},
	},
	{
		Text: `${string/position/\/leng\\th}`,
		Node: &FuncNode{
			Param: "string",
			Name:  "/",
			Args: []Node{
				&TextNode{
					Value: "position",
				},
				&TextNode{
					Value: "/leng\\th",
				},
			},
		},
	},

	// functions in functions
	{
		Text: "${string:${position}}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&FuncNode{
					Param: "position",
				},
			},
		},
	},
	{
		Text: "${string:${stringy:position:length}:${stringz,,}}",
		Node: &FuncNode{
			Param: "string",
			Name:  ":",
			Args: []Node{
				&FuncNode{
					Param: "stringy",
					Name:  ":",
					Args: []Node{
						&TextNode{Value: "position"},
						&TextNode{Value: "length"},
					},
				},
				&FuncNode{
					Param: "stringz",
					Name:  ",,",
				},
			},
		},
	},
	{
		Text: "${string#${stringz}}",
		Node: &FuncNode{
			Param: "string",
			Name:  "#",
			Args: []Node{
				&FuncNode{Param: "stringz"},
			},
		},
	},
	{
		Text: "${string=${stringz}}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&FuncNode{Param: "stringz"},
			},
		},
	},
	{
		Text: "${string=prefix-${var}}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&TextNode{Value: "prefix-"},
				&FuncNode{Param: "var"},
			},
		},
	},
	{
		Text: "${string=${var}-suffix}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&FuncNode{Param: "var"},
				&TextNode{Value: "-suffix"},
			},
		},
	},
	{
		Text: "${string=prefix-${var}-suffix}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&TextNode{Value: "prefix-"},
				&FuncNode{Param: "var"},
				&TextNode{Value: "-suffix"},
			},
		},
	},
	{
		Text: "${string=prefix${var} suffix}",
		Node: &FuncNode{
			Param: "string",
			Name:  "=",
			Args: []Node{
				&TextNode{Value: "prefix"},
				&FuncNode{Param: "var"},
				&TextNode{Value: " suffix"},
			},
		},
	},
	{
		Text: "${string//${stringy}/${stringz}}",
		Node: &FuncNode{
			Param: "string",
			Name:  "//",
			Args: []Node{
				&FuncNode{Param: "stringy"},
				&FuncNode{Param: "stringz"},
			},
			buf: buf("${string//${stringy}/${stringz}}"),
		},
	},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Text, func(t *testing.T) {
			got, err := Parse(test.Text)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.Node, got.Root)
		})
	}
}
