var today = new Date();
var thisYear = today.getFullYear();
var thisMonth = ('0' + (today.getMonth()+1)).slice(-2);
var thisDay = ("0"+today.getDate()).slice(-2);
var thisHour = ("0"+today.getHours()).slice(-2);
var thisMinute = ("0"+today.getMinutes()).slice(-2);
var thisSecond = today.getSeconds()

function toDate() {
    return (thisYear + "-" + thisMonth + "-" + thisDay);
}

function invoiceNumber() {
    return (thisYear.toString().slice(-2) + thisMonth + thisDay + thisHour + thisMinute + thisSecond);
}
