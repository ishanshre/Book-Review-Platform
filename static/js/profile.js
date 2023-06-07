window.addEventListener('DOMContentLoaded', () => {
	const imageDataFront = document.querySelector('meta[name="citizen_front"]').content;
	const imageDataBack = document.querySelector('meta[name="citizen_back"]').content;
	const imgElement = document.getElementById('citizenship_front');
	const imgElement1 = document.getElementById('citizenship_back');
	imgElement.src = 'data:image/jpeg;base64,' + imageDataFront;
	imgElement1.src = 'data:image/jpeg;base64,' + imageDataBack;
    
});
