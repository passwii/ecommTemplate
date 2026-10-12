import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { isSafeProductResourceUrl } from "./product-resource-url.mjs";

describe("isSafeProductResourceUrl", () => {
  it("rejects single and repeated backslashes in site paths", () => {
    assert.equal(isSafeProductResourceUrl("/\\evil.example"), false);
    assert.equal(isSafeProductResourceUrl("/\\\\evil.example"), false);
    assert.equal(isSafeProductResourceUrl("/safe\\path/spec-sheet.pdf"), false);
  });

  it("accepts paths that resolve to the storefront origin", () => {
    assert.equal(isSafeProductResourceUrl("/catalog/spec-sheet.pdf"), true);
    assert.equal(isSafeProductResourceUrl("/files/../catalog/spec-sheet.pdf"), true);
  });

  it("accepts HTTP(S) external links and rejects other protocols", () => {
    assert.equal(isSafeProductResourceUrl("https://docs.example/spec-sheet.pdf"), true);
    assert.equal(isSafeProductResourceUrl("http://docs.example/spec-sheet.pdf"), true);
    assert.equal(isSafeProductResourceUrl("javascript:alert(1)"), false);
    assert.equal(isSafeProductResourceUrl("//evil.example/spec-sheet.pdf"), false);
  });
});
