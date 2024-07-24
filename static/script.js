document.getElementById("form").addEventListener("submit", function (event) {
    event.preventDefault();

    let formData = new FormData();
    let fileInput = document.getElementById("file-upload");
    let textInput = document.getElementById("input-text").value;

    formData.append("file", fileInput.files[0]);
    formData.append("text", "abc");
    console.log(formData.get("text"));

    fetch("/api/posts", {
        method: 'POST',
        body: formData,
        headers: {'Content-Type': 'multipart/form-data'}
    }).then(function (response) {
        return response.json();
    }).then(function (data) {
        console.log(data);
    })
    location.reload();
});

fetch("/api/posts", {method: "GET"}).then(function (response) {
    return response.json();
}).then(function (data) {
    let msgList = data.data.reverse();
    msgList.forEach(function (item) {
        const sectionDiv = document.getElementById("section");
        const boxDiv = document.createElement("div");
        boxDiv.className = "box";
        sectionDiv.appendChild(boxDiv);
        const boxTextDiv = document.createElement("div");
        boxTextDiv.className = "box-text";
        boxTextDiv.id = "box-texxt";
        boxTextDiv.innerHTML = item["Text"];
        boxDiv.appendChild(boxTextDiv);
        const imgDiv = document.createElement("img");
        imgDiv.className = "box-img";
        imgDiv.id = "box-img";
        imgDiv.alt = "picture";
        imgDiv.src = item["ObjectUrl"]
        boxDiv.appendChild(imgDiv)
    })
})