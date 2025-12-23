// Event Delegation System - Replaces inline onclick handlers

document.addEventListener('DOMContentLoaded', () => {
	// Collapsible toggle - handles both .collapsible and legacy .slice/.container
	document.addEventListener('click', (e) => {
		const header = e.target.closest('.collapsible__header');
		if (header) {
			const collapsible = header.closest('.collapsible');
			if (collapsible) {
				collapsible.classList.toggle('collapsible--collapsed');
			}
			return;
		}

		// Legacy support for .slice and .container classes (until migration complete)
		const legacyHeader = e.target.closest('.slice__header, .container__header');
		if (legacyHeader) {
			const parent = legacyHeader.closest('.slice, .container');
			if (parent) {
				parent.classList.toggle('collapsed');
			}
		}
	});
});

window.toggleCollapse = function(element) {
	const parent = element.closest('.collapsible, .slice, .container');
	if (parent) {
		parent.classList.toggle(parent.classList.contains('collapsible') ? 'collapsible--collapsed' : 'collapsed');
	}
};

window.toggleNavItem = function(element) {
	const navItem = element.closest('.nav__item--collapsible');
	if (navItem) {
		navItem.classList.toggle('collapsed');
	}
};
