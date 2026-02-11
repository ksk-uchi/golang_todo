import eslint from "@eslint/js";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";
import reactPlugin from "eslint-plugin-react";
import { defineConfig, globalIgnores } from "eslint/config";
import tslint from "typescript-eslint";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
  ]),
  reactPlugin.configs.flat.recommended,
  {
    files: ["src/**/*.ts", "src/**/*.tsx"],
    extends: [eslint.configs.recommended, tslint.configs.recommended],
    rules: {
      "@typescript-eslint/array-type": "error",
      "react/react-in-jsx-scope": "off",
      "react/no-danger": "error",
    },
  },
]);

export default eslintConfig;
