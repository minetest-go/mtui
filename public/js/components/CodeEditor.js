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
            this.cm.setValue(this.modelValue);
        }
    },
    template: /*html*/`
    <textarea ref="textarea">{{modelValue}}</textarea>
    `
};