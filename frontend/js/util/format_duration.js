
const ranges = [{
    start: 60*60*24*356, suffix: "year(s)"
},{
    start: 60*60*24*30, suffix: "month(s)"
},{
    start: 60*60*24*7, suffix: "week(s)"
},{
    start: 60*60*24, suffix: "day(s)"
},{
    start: 60*60, suffix: "hour(s)"
},{
    start: 60, suffix: "minute(s)"
},{
    start: 1, suffix: "second(s)"
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
