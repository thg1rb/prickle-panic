import { cardAttributes } from "./cocoCardVariant";

type LocalCardArtProps = {
  variant: number;
  className?: string;
};

export function LocalCardArt({ variant, className }: LocalCardArtProps) {
  const [bodyCream, longQuills, glasses, scarf, strawberry] =
    cardAttributes(variant);
  const body = bodyCream ? "#f1d9b8" : "#79543d";
  const bodyShade = bodyCream ? "#d5b28b" : "#553b2d";

  return (
    <svg
      aria-hidden="true"
      className={className}
      viewBox="0 0 160 156"
      focusable="false"
    >
      <g
        stroke="#48352a"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="4"
      >
        <path
          d={
            longQuills
              ? "M44 66 37 28l24 20 3-38 19 35 20-33 1 40 31-20-15 39"
              : "M46 68 44 43l18 15 4-27 16 27 18-25 2 30 22-13-14 29"
          }
          fill={bodyShade}
        />
        <ellipse cx="82" cy="89" fill={body} rx="50" ry="47" />
        <path
          d="M48 111c14 12 48 15 68 1"
          fill="none"
          stroke={bodyShade}
          strokeWidth="5"
        />
        {scarf ? (
          <path
            d="M36 107c27 13 62 13 92-1l-2 15c-31 14-62 14-87 0z"
            fill="#e57570"
          />
        ) : null}
        {glasses ? (
          <g fill="none" stroke="#42382f" strokeWidth="4">
            <circle cx="65" cy="83" r="12" />
            <circle cx="99" cy="83" r="12" />
            <path d="M77 82h10M53 79l-8-4M111 79l8-4" />
          </g>
        ) : (
          <g fill="#302820" stroke="none">
            <circle cx="65" cy="83" r="4" />
            <circle cx="99" cy="83" r="4" />
          </g>
        )}
        <path
          d="M78 97q5 5 10 0"
          fill="none"
          stroke="#654736"
          strokeWidth="3"
        />
        {strawberry ? (
          <g>
            <path
              d="M117 119c-2-13 20-17 22-2 1 9-8 19-13 22-5-4-9-12-9-20z"
              fill="#e96962"
            />
            <path d="m119 119 3-8 6 5 5-7 4 9" fill="#6d965f" strokeWidth="2" />
            <circle cx="125" cy="124" r="1.5" fill="#fff4d9" stroke="none" />
            <circle cx="132" cy="127" r="1.5" fill="#fff4d9" stroke="none" />
          </g>
        ) : null}
      </g>
    </svg>
  );
}
