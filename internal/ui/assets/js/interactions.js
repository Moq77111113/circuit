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

	// Mobile sidebar toggle
	const menuToggle = document.querySelector('.mobile-menu-toggle');
	if (menuToggle) {
		menuToggle.addEventListener('click', toggleSidebar);
	}

	// Close sidebar overlay
	const overlay = document.querySelector('.mobile-sidebar-overlay');
	if (overlay) {
		overlay.addEventListener('click', closeSidebar);
	}
});

// Sidebar helper functions
function toggleSidebar() {
	const sidebar = document.querySelector('.app__sidebar');
	if (sidebar) {
		sidebar.classList.toggle('is-open');
	}
}

function closeSidebar() {
	const sidebar = document.querySelector('.app__sidebar');
	if (sidebar) {
		sidebar.classList.remove('is-open');
	}
}

// Export for inline usage (if needed during transition)
window.toggleCollapse = function(element) {
	const parent = element.closest('.collapsible, .slice, .container');
	if (parent) {
		parent.classList.toggle(parent.classList.contains('collapsible') ? 'collapsible--collapsed' : 'collapsed');
	}
};
