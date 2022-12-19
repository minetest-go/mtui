
const ranges = [{
    start: 60*60*24*7, suffix: "weeks"
},{
    start: 60*60*24, suffix: "days"
},{
    start: 60*60, suffix: "hours"
},{
    start: 60, suffix: "minutes"
},{
    start: 0, suffix: "seconds"
}];

export default function(seconds) {
    for (let i=0; i<ranges.length; i++) {
        let range = ranges[i];
        if (seconds > range.start) {
            const units = Math.floor(seconds / range.start);
            return `${units} ${range.suffix}`;
        }
    }
}
