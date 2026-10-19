package game

import "fmt"

// Variant represents a Coco card's five binary attributes. Attribute order is
// body color, quills, glasses, scarf, strawberry, most significant bit first.
type Variant uint8

const (
	AttributeCount = 5
	VariantCount   = 1 << AttributeCount
	MaxVariant     = VariantCount - 1
)

// NewVariant validates an externally supplied variant value.
func NewVariant(value uint8) (Variant, error) {
	if value > MaxVariant {
		return 0, fmt.Errorf("variant %d is outside 0..%d", value, MaxVariant)
	}
	return Variant(value), nil
}

// Attributes returns values in the documented attribute order.
func (v Variant) Attributes() [AttributeCount]bool {
	var attributes [AttributeCount]bool
	for index := range attributes {
		attributes[index] = v&(1<<(AttributeCount-1-index)) != 0
	}
	return attributes
}

// HammingDistance returns the number of different attributes.
func HammingDistance(left, right Variant) int {
	return bitCount(uint8(left ^ right))
}

// Matches reports whether a selected card differs in zero or one attributes.
func Matches(reference, selected Variant) bool {
	return HammingDistance(reference, selected) <= 1
}

func bitCount(value uint8) int {
	count := 0
	for value != 0 {
		count += int(value & 1)
		value >>= 1
	}
	return count
}
