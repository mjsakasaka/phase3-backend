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

});