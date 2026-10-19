export const ATTRIBUTE_COUNT = 5;
export const VARIANT_COUNT = 1 << ATTRIBUTE_COUNT;
export const BOARD_SIZE = 16;

export type CardAttributes = readonly [
  bodyCream: boolean,
  longQuills: boolean,
  glasses: boolean,
  scarf: boolean,
  strawberry: boolean,
];

export function isCardVariant(value: number): value is number {
  return Number.isInteger(value) && value >= 0 && value < VARIANT_COUNT;
}

export function cardAttributes(value: number): CardAttributes {
  if (!isCardVariant(value)) {
    throw new RangeError(
      `Card variant must be an integer from 0 to ${VARIANT_COUNT - 1}.`,
    );
  }

  return [
    (value & 0b10000) !== 0,
    (value & 0b01000) !== 0,
    (value & 0b00100) !== 0,
    (value & 0b00010) !== 0,
    (value & 0b00001) !== 0,
  ];
}

export function hammingDistance(reference: number, selected: number): number {
  if (!isCardVariant(reference) || !isCardVariant(selected)) {
    throw new RangeError("Card variants must be integers from 0 to 31.");
  }

  let difference = reference ^ selected;
  let distance = 0;
  while (difference !== 0) {
    distance += difference & 1;
    difference >>>= 1;
  }
  return distance;
}

export function isValidPick(reference: number, selected: number): boolean {
  return hammingDistance(reference, selected) <= 1;
}
