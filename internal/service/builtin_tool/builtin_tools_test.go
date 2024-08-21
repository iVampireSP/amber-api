package builtin_tool

import (
	"context"
	"rag-new/internal/schema"
	"testing"
)

func TestToolCall(t *testing.T) {
	s := NewService()

	function, err := s.CallFunction(context.Background(), "now", schema.FunctionCallArguments{})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(function)
}

func TestToolCallWithPrefix(t *testing.T) {
	s := NewService()

	function, err := s.CallFunction(context.Background(), "now", schema.FunctionCallArguments{})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(function)
}
