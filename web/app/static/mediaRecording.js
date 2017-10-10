navigator.getUserMedia(mediaConstraints, onMediaSuccess, onMediaError);
var mediaConstraints = {
    audio: true,
    video: true
};
function onMediaError(e) {
    console.error('media error', e);
}

console.log("Hello from a script");
function onMediaSuccess(stream) {
    var mediaRecorder = new MediaStreamRecorder(stream);
    mediaRecorder.mimeType = 'audio/wav'; // check this line for audio/wav
    mediaRecorder.ondataavailable = function (blob) {
        // POST/PUT "Blob" using FormData/XHR2
        var blobURL = URL.createObjectURL(blob);
        document.write('<a href="' + blobURL + '">' + blobURL + '</a>');
    };
    mediaRecorder.start(3000);
}

function upload(blob) {
    var file = new File([blob], 'msr-' + (new Date).toISOString().replace(/:|\./g, '-') + '.wav', {
        type: 'audio/wav'
    });

    // create FormData
    var formData = new FormData();
    formData.append('video-filename', file.name);
    formData.append('video-blob', file);

    makeXMLHttpRequest('https://path-to-your-server/save.php', formData, function() {
        var downloadURL = 'https://path-to-your-server/uploads/' + file.name;
        console.log('File uploaded to this path:', downloadURL);
    });
}
function onMediaError(e) {
    console.error('media error', e);
}
