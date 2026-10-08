import globals from "globals";
import pluginJs from "@eslint/js";
import pluginVue from "eslint-plugin-vue";
import stylisticJs from '@stylistic/eslint-plugin-js';

export default [
  {
    ignores: ["dist/"]
  },
  {
    files: ["**/*.js"]
  },
  {
    languageOptions: {
      globals: {
        ...globals.browser
      }
    }
  },
  {
    plugins: {
      '@stylistic/js': stylisticJs
    },
    rules: {
      '@stylistic/js/semi': "error",
    }
  },
  pluginJs.configs.recommended,
  ...pluginVue.configs["flat/essential"],
  {
    // legacy patterns in the existing components
    rules: {
      "vue/multi-word-component-names": "off",
      "vue/require-v-for-key": "off",
      "vue/valid-v-for": "off",
      "vue/no-mutating-props": "off",
      "vue/no-async-in-computed-properties": "off",
      "vue/no-textarea-mustache": "off",
      "vue/return-in-computed-property": "off",
      "vue/no-use-v-if-with-v-for": "off",
      "vue/no-unused-components": "off"
    }
  }
];