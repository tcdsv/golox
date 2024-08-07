package loxvalue_test

import (
	loxvalue "golox/value"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValue_IsEqualTrue(t *testing.T) {
	require.True(t, loxvalue.IsEqual(loxvalue.Nil{}, loxvalue.Nil{}))
}

func TestValue_IsEqualFalse(t *testing.T) {
	// require.False(t, loxvalue.IsEqual(loxvalue.Nil{}, loxvalue.String{Value: "foo"}))
	require.False(t, loxvalue.IsEqual(loxvalue.Boolean{Value: true}, loxvalue.Boolean{Value: false}))
}

func TestValue_IsTruthyTrue(t *testing.T) {
	require.True(t, loxvalue.IsTruthy(loxvalue.Boolean{Value: true}))
	require.True(t, loxvalue.IsTruthy(loxvalue.String{Value: ""}))
}

func TestValue_IsTruthyFalse(t *testing.T) {
	require.False(t, loxvalue.IsTruthy(loxvalue.Boolean{Value: false}))
	require.False(t, loxvalue.IsTruthy(loxvalue.Nil{}))
}