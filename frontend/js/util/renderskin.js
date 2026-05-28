
// converts a raw skin url to a properly rendered data-url
export default function(src) {
    return new Promise(resolve => {
        let img = new Image();
        img.onload = function() {
            const canvas = document.createElement("canvas");
            const ctx = canvas.getContext("2d");
            canvas.width = 16;
            canvas.height = 32;

            // head
            ctx.drawImage(img, 8, 8, 8, 8, 4, 0, 8, 8);
            // chest
            ctx.drawImage(img, 20, 20, 8, 12, 4, 8, 8, 12);
            // leg left
            ctx.drawImage(img, 4, 20, 4, 12, 4, 20, 4, 12);
            // leg right
            if (img.height === 64) {
                ctx.drawImage(img, 20, 52, 4, 12, 8, 20, 4, 12);
            } else {
                ctx.drawImage(img, 4, 20, 4, 12, 8, 20, 4, 12);
            }
            // arm left
            ctx.drawImage(img, 44, 20, 4, 12, 0, 8, 4, 12);
            // arm right
            if (img.height === 64) {
                ctx.drawImage(img, 36, 52, 4, 12, 12, 8, 4, 12);
            } else {
                ctx.drawImage(img, 44, 20, 4, 12, 12, 8, 4, 12);
            }

            resolve(canvas.toDataURL("image/png"));
        };

        // trigger source image load
        img.src = src;
    });
}