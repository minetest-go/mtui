export default function(seconds) {
    if (!seconds){
        return "";
    }
    const minutes = Math.floor(+seconds / 60);
    const rest_seconds = Math.floor(+seconds - (minutes * 60));
  
    return minutes + ":" + (rest_seconds > 9 ? rest_seconds : "0" + rest_seconds);  
}
