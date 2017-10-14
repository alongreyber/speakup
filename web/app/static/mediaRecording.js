var mediaConstraints = {
    audio: true,
    video: true
};
navigator.getUserMedia(mediaConstraints, onMediaSuccess, onMediaError);
function onMediaError(e) {
    console.error('media error', e);
}

function onMediaSuccess(stream) {
    var mediaRecorder = new MediaStreamRecorder(stream);
    mediaRecorder.mimeType = 'audio/wav'; // check this line for audio/wav
    mediaRecorder.ondataavailable = function (blob) {
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

    makeXMLHttpRequest('https://submit_gentle', formData, function() {
        console.log('File uploaded');
    });
}
function onMediaError(e) {
    console.error('media error', e);
}
