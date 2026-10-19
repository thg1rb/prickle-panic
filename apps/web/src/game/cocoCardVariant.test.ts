import { describe, expect, it } from "vitest";
import {
  ATTRIBUTE_COUNT,
  VARIANT_COUNT,
  cardAttributes,
  describeCard,
  hammingDistance,
  isCardVariant,
  isValidPick,
} from "./cocoCardVariant";

describe("Coco card variants", () => {
  it("covers exactly the 32 five-attribute variations", () => {
    expect(ATTRIBUTE_COUNT).toBe(5);
    expect(VARIANT_COUNT).toBe(32);
    expect(
      new Set(Array.from({ length: VARIANT_COUNT }, (_, value) => value)).size,
    ).toBe(32);
  });

  it("evaluates every ordered pair using the one-attribute rule", () => {
    for (let reference = 0; reference < VARIANT_COUNT; reference += 1) {
      for (let selected = 0; selected < VARIANT_COUNT; selected += 1) {
        const expectedDistance = (reference ^ selected)
          .toString(2)
          .replaceAll("0", "").length;
        expect(hammingDistance(reference, selected)).toBe(expectedDistance);
        expect(isValidPick(reference, selected)).toBe(expectedDistance <= 1);
      }
    }
  });

  it("maps binary attributes in the documented order and describes them in Thai", () => {
    expect(cardAttributes(0b10110)).toEqual([true, false, true, true, false]);
    expect(describeCard(0b10110)).toContain("ตัวสีครีม");
    expect(describeCard(0b10110)).toContain("ใส่แว่น");
    expect(isCardVariant(31)).toBe(true);
    expect(isCardVariant(32)).toBe(false);
  });

  it("rejects invalid values when converting or comparing cards", () => {
    expect(() => cardAttributes(-1)).toThrow(RangeError);
    expect(() => hammingDistance(0, 32)).toThrow(RangeError);
  });
});
