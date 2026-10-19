package game

import "testing"

func TestVariantAttributeOrder(t *testing.T) {
	tests := []struct {
		value    uint8
		wantBits [AttributeCount]bool
	}{
		{value: 0, wantBits: [AttributeCount]bool{false, false, false, false, false}},
		{value: 16, wantBits: [AttributeCount]bool{true, false, false, false, false}},
		{value: 8, wantBits: [AttributeCount]bool{false, true, false, false, false}},
		{value: 4, wantBits: [AttributeCount]bool{false, false, true, false, false}},
		{value: 2, wantBits: [AttributeCount]bool{false, false, false, true, false}},
		{value: 1, wantBits: [AttributeCount]bool{false, false, false, false, true}},
		{value: 31, wantBits: [AttributeCount]bool{true, true, true, true, true}},
	}
	for _, test := range tests {
		variant, err := NewVariant(test.value)
		if err != nil {
			t.Fatalf("NewVariant(%d): %v", test.value, err)
		}
		if got := variant.Attributes(); got != test.wantBits {
			t.Errorf("Variant(%d).Attributes() = %v, want %v", test.value, got, test.wantBits)
		}
	}
	if _, err := NewVariant(32); err == nil {
		t.Fatal("NewVariant(32) succeeded; want out-of-range error")
	}
}

func TestMatchingExhaustivelyChecksAllOrderedPairs(t *testing.T) {
	pairs := 0
	for reference := 0; reference < VariantCount; reference++ {
		for selected := 0; selected < VariantCount; selected++ {
			wantDistance := naiveDistance(uint8(reference), uint8(selected))
			gotDistance := HammingDistance(Variant(reference), Variant(selected))
			if gotDistance != wantDistance {
				t.Fatalf("distance(%05b, %05b) = %d, want %d", reference, selected, gotDistance, wantDistance)
			}
			if got, want := Matches(Variant(reference), Variant(selected)), wantDistance <= 1; got != want {
				t.Fatalf("matches(%05b, %05b) = %t, want %t", reference, selected, got, want)
			}
			pairs++
		}
	}
	if pairs != 1024 {
		t.Fatalf("checked %d ordered pairs, want 1024", pairs)
	}
}

func naiveDistance(left, right uint8) int {
	count := 0
	for bit := 0; bit < AttributeCount; bit++ {
		mask := uint8(1 << bit)
		if left&mask != right&mask {
			count++
		}
	}
	return count
}
