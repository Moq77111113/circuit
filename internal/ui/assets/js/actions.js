// Confirmation dialog for actions
window.confirmAction = function(label, description) {
	const msg = description
		? `Are you sure you want to execute "${label}"?\n\n${description}`
		: `Are you sure you want to execute "${label}"?`;
	return confirm(msg);
};
