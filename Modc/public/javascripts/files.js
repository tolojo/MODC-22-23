async function download() {
    try {

        let file_name = {
            file_name: document.getElementById("file_name").value,
          
        }
        let user = await $.ajax({
            url: '/api/files/download/',
            method: 'post',
            dataType: 'json',
            data: JSON.stringify(file_name),
            contentType: 'application/json'
        }).then(

        fetch('/api/files/download/', { method: 'get', mode: 'no-cors', referrerPolicy: 'no-referrer' })
            .then(res => res.blob())
            .then(res => {
                const aElement = document.createElement('a');
                aElement.setAttribute('download', "download.txt");
                const href = URL.createObjectURL(res);
                aElement.href = href;
                aElement.setAttribute('target', '_blank');
                aElement.click();
                URL.revokeObjectURL(href);
            }));
    } catch (err) {
        console.log(err);
    }
}