import { fromTextArea } from 'codemirror';

export default {
    props: ["mode", "modelValue"],
    emits: ['update:modelValue'],
    mounted: function() {
        // TODO: languages, https://codemirror.net/examples/basic/
        this.cm = fromTextArea(this.$refs.textarea, {
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
    },
    template: /*html*/`
    <textarea ref="textarea">{{modelValue}}</textarea>
    `
};