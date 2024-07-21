import globals from "globals";
import pluginJs from "@eslint/js";
import pluginVue from "eslint-plugin-vue";
import stylisticJs from '@stylistic/eslint-plugin-js'

export default [
  {
    ignores: ["js/bundle.js"]
  },
  {
    files: ["**/*.js"]
  },
  {
    languageOptions: {
      globals: {
        Vue: "readonly",
        VueRouter: "readonly",
        CodeMirror: "readonly",
        Chart: "readonly",
        VueDatePicker: "readonly",
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
];