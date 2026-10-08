<script>
import CodeMirror from "codemirror";
import "codemirror/mode/lua/lua.js";
import "codemirror/mode/javascript/javascript.js";
import "codemirror/mode/htmlmixed/htmlmixed.js";
import "codemirror/mode/xml/xml.js";
import "codemirror/mode/css/css.js";
import "codemirror/mode/toml/toml.js";

export default {
    props: ["mode", "modelValue"],
    emits: ['update:modelValue'],
    mounted: function() {
        this.cm = CodeMirror.fromTextArea(this.$refs.textarea, {
            lineNumbers: true,
            viewportMargin: Infinity,
            mode: {
                name: this.mode
            }
        });
        this.cm.on("change", () => {
            this.$emit('update:modelValue', this.cm.getValue());
        });
    },
    watch: {
        "modelValue": function() {
            if (this.cm.getValue() != this.modelValue){
                // only update if value has changed from outside (cursor move)
                this.cm.setValue(this.modelValue);
            }
        }
    }
};
</script>

<template>
  <textarea ref="textarea">{{modelValue}}</textarea>
</template>
