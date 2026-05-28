import globals from "globals";
import pluginJs from "@eslint/js";
import pluginVue from "eslint-plugin-vue";
import stylisticJs from '@stylistic/eslint-plugin';

export default [
  {
    ignores: ["js/bundle.js"]
  },{
    files: ["**/*.js"]
  },{
    languageOptions: {
      globals: globals.browser
    }
  },{
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