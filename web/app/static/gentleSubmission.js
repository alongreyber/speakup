var app = angular.module('app', []);
app.controller('gentleSubmission', function($scope) {
    $scope.mediaRequested = false;
    $scope.recording = false;
    $scope.startRecording = function() {
	if(!$scope.mediaRequested) {
	    var mediaConstraints = {
		audio: true
	    };
	    navigator.getUserMedia(mediaConstraints, onMediaSuccess, onMediaError);
	    $scope.mediaRequested = true;
	}
	mediaRecorder.start(5 * 1000);
	$scope.recording = true;
    }
    $scope.stopRecording = function() {

    }
});

function onMediaError(e) {
    console.error('media error', e);
}

function onMediaSuccess(stream) {
    var mediaRecorder = new MediaStreamRecorder(stream);
    mediaRecorder.mimeType = 'audio/wav'; 
    mediaRecorder.ondataavailable = function (blob) {
	var blobURL = URL.createObjectURL(blob);
    };
    mediaRecorder.start();
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
