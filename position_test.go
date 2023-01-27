package bfldb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToOrder(t *testing.T) {
	tests := []struct {
		name string
		p    Position
		want Order
	}{
		{
			name: "opened position",
			p:    Position{Direction: Long, Amount: 1, PrevAmount: 0, Type: Opened},
			want: Order{Direction: Long, Amount: 1, ReduceOnly: false},
		},
		{
			name: "closed position",
			p:    Position{Direction: Long, Amount: 0, PrevAmount: 1, Type: Closed},
			want: Order{Direction: Short, Amount: 1, ReduceOnly: true},
		},
		{
			name: "added to position",
			p:    Position{Direction: Long, Amount: 1, PrevAmount: 0.5, Type: AddedTo},
			want: Order{Direction: Long, Amount: 0.5, ReduceOnly: false},
		},
		{
			name: "partially closed position",
			p:    Position{Direction: Long, Amount: 0.1, PrevAmount: 1, Type: PartiallyClosed},
			want: Order{Direction: Short, Amount: 0.9, ReduceOnly: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.ToOrder()
			require.Equalf(t, tt.want, got, "Position.ToOrder() = %v, want %v", got, tt.want)
		})
	}
}
