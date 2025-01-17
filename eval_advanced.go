package envsubst

import (
	"bytes"
	"io"

	"github.com/logandavies181/envsubst/parse"
)

type NodeInfo struct {
	node parse.Node
	args []string
	name string
}

// Orig returns the original text of the substitution template,
// before it was parsed. This can be used to provide full context
// for custom mapping functions or leave expressions un-evaluated
func (n NodeInfo) Orig() string {
	return parse.FormatNode(n.node)
}

// Args returns the arguments to the shell-style substitution function
func (n NodeInfo) Args() []string {
	return n.args
}

// Fn returns the string representing the shell-style substitution function
// e.g. `:-`
func (n NodeInfo) Fn() string {
	return n.name
}

// Result returns the value that will be set by the substitution function
// if it runs
func (n NodeInfo) Result(mapResult string) string {
	fn := lookupFunc(n.Fn(), len(n.Args()))

	return fn(mapResult, n.Args()...)
}

// AdvancedMapping is a function that takes a variable name and
// representation of the full shell variable string and returns the substituted
// string and whether or not to continue processing
type AdvancedMapping func(string, NodeInfo) (mapped string, shouldContinue bool)

// EvalAdvanced allows the caller to control how ${var} is mapped and how its
// nested parameters are evaluated.
//
// If mapping returns false, processing stops and the returned string is used.
// If mapping returns true, this behaves the same as EvalEnv
func EvalAdvanced(s string, mapping AdvancedMapping) (string, error) {
	t, err := Parse(s)
	if err != nil {
		return s, err
	}
	return t.ExecuteAdvanced(mapping)
}

// ExecuteAdvanced applies a parsed template to the specified data mapping,
// allowing greater control over execution
func (t *Template) ExecuteAdvanced(mapping AdvancedMapping) (str string, err error) {
	b := new(bytes.Buffer)
	s := new(state)
	s.node = t.tree.Root
	s.advMapper = mapping
	s.writer = b
	err = t.evalAdvanced(s)
	if err != nil {
		return
	}
	return b.String(), nil
}

func (t *Template) evalAdvanced(s *state) (err error) {
	switch node := s.node.(type) {
	case *parse.TextNode:
		err = t.evalText(s, node)
	case *parse.FuncNode:
		err = t.evalAdvancedFunc(s, node)
	case *parse.ListNode:
		err = t.evalAdvancedList(s, node)
	}
	return err
}

func (t *Template) evalAdvancedList(s *state, node *parse.ListNode) (err error) {
	for _, n := range node.Nodes {
		s.node = n
		err = t.evalAdvanced(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Template) evalAdvancedFunc(s *state, node *parse.FuncNode) error {
	var w = s.writer
	var buf bytes.Buffer
	var args []string

	for _, n := range node.Args {
		buf.Reset()
		s.writer = &buf
		s.node = n
		err := t.evalAdvanced(s)
		if err != nil {
			return err
		}
		args = append(args, buf.String())
	}

	// restore the origin writer
	s.writer = w
	s.node = node

	v, shouldContinue := s.advMapper(node.Param, NodeInfo{node, args, node.Name})
	if !shouldContinue {
		_, err := io.WriteString(s.writer, v)
		return err
	}

	fn := lookupFunc(node.Name, len(args))

	_, err := io.WriteString(s.writer, fn(v, args...))
	return err
}
