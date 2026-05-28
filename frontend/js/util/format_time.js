
function pad(n) {
    return n<10 ? `0${n}` : `${n}`;
}

export default function(ts) {
    const d = new Date(ts*1000);
    return `${d.getYear()+1900}-${pad(d.getMonth()+1)}-${pad(d.getDate())}` +
        ` ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}