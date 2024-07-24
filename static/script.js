document.getElementById("form").addEventListener("submit", function (event) {
    event.preventDefault();

    const url = '/api/posts'; // 替换为你的API URL
    const formData = new FormData();

    let textInput = document.getElementById("input-text").value;
    formData.append('text', textInput);

    const fileInput = document.getElementById("file-upload");
    if (fileInput.files.length > 0) {
        formData.append('file', fileInput.files[0]);
    }

    fetch(url, {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        console.log('成功:', data);
        location.reload()

    })
    .catch((error) => {
        console.error('错误:', error);
    });
});

fetch("/api/posts", {method: "GET"}).then(function (response) {
    return response.json();
}).then(function (data) {
    if (data.data != null) {
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
    }
})