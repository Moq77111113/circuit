// Actions dropdown toggle
window.toggleActionsDropdown = function() {
	const menu = document.getElementById('actions-menu');
	if (!menu) return;

	const isVisible = menu.classList.contains('show');

	if (isVisible) {
		menu.classList.remove('show');
	} else {
		menu.classList.add('show');
	}
};

// Close dropdown when clicking outside
document.addEventListener('click', function(event) {
	const menu = document.getElementById('actions-menu');
	const button = document.querySelector('.actions-button');

	if (!menu || !button) return;

	// If click is outside menu and button, close the menu
	if (!menu.contains(event.target) && !button.contains(event.target)) {
		menu.classList.remove('show');
	}
});

// Confirmation dialog for actions
window.confirmAction = function(label, description) {
	const msg = description
		? `Are you sure you want to execute "${label}"?\n\n${description}`
		: `Are you sure you want to execute "${label}"?`;
	return confirm(msg);
};
