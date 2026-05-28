const ranges = [{
	size: 1000*1000,
	suffix: "M"
},{
	size: 1000,
	suffix: "K"
}];

export default function(value) {
	let count = +value;
	for (let i=0; i<ranges.length; i++){
		let range = ranges[i];
		if (count > range.size){
			return (Math.floor(count / range.size * 100) / 100) + " " + range.suffix;
		}
	}

	return count;
}
