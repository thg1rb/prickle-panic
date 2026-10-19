import { readdir, readFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const protocolRoot = join(root, "docs/protocol/v1");
const schema = JSON.parse(
  await readFile(join(protocolRoot, "message.schema.json"), "utf8"),
);
const ajv = new Ajv2020({ allErrors: true, strict: true });
addFormats(ajv);
const validate = ajv.compile(schema);

let fixtureCount = 0;
for (const [category, shouldValidate] of [
  ["valid", true],
  ["invalid", false],
]) {
  const fixtureDirectory = join(protocolRoot, "fixtures", category);
  const files = (await readdir(fixtureDirectory))
    .filter((file) => file.endsWith(".json"))
    .toSorted();

  if (files.length === 0) {
    throw new Error(`No ${category} protocol fixtures found.`);
  }

  for (const file of files) {
    const fixture = JSON.parse(
      await readFile(join(fixtureDirectory, file), "utf8"),
    );
    const actual = validate(fixture);
    fixtureCount += 1;

    if (actual !== shouldValidate) {
      const errors = (validate.errors ?? [])
        .map((error) => `${error.instancePath || "/"} ${error.message}`)
        .join("; ");
      throw new Error(
        `${category}/${file} expected ${shouldValidate ? "valid" : "invalid"}${errors ? `: ${errors}` : ""}`,
      );
    }
  }
}

console.log(`Validated ${fixtureCount} protocol fixtures (valid and invalid).`);
