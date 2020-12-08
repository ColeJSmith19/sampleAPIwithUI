function getRequest() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", "http://localhost:18080/test", false);
    xmlHttp.send(null);
    document.getElementById("demo").innerHTML = xmlHttp.responseText;
}

function postRequest() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", "http://localhost:18080/testpost", false);
    xmlHttp.setRequestHeader('Content-Type', 'application/json');
    xmlHttp.send(JSON.stringify({
        value: "value"
    }));
}