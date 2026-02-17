// Previous/Next page navigation at bottom
(function() {
    'use strict';

    function createPageNavigation() {
        const content = document.querySelector('.content main');
        if (!content) return;

        // Remove existing nav if present
        const existingNav = content.querySelector('.page-navigation');
        if (existingNav) {
            existingNav.remove();
        }

        // Get current page info from sidebar
        const sidebarLinks = document.querySelectorAll('.sidebar .chapter li a');
        let currentIndex = -1;
        
        sidebarLinks.forEach((link, index) => {
            if (link.classList.contains('active')) {
                currentIndex = index;
            }
        });

        if (currentIndex === -1) return;

        // Create navigation container
        const navContainer = document.createElement('div');
        navContainer.className = 'page-navigation';

        // Previous link
        if (currentIndex > 0) {
            const prevLink = sidebarLinks[currentIndex - 1];
            const prevButton = document.createElement('a');
            prevButton.href = prevLink.href;
            prevButton.className = 'page-nav-button page-nav-prev';
            prevButton.innerHTML = `
                <span class="page-nav-label">Previous</span>
                <span class="page-nav-title">${prevLink.textContent}</span>
            `;
            navContainer.appendChild(prevButton);
        } else {
            // Empty spacer
            const spacer = document.createElement('div');
            navContainer.appendChild(spacer);
        }

        // Next link
        if (currentIndex < sidebarLinks.length - 1) {
            const nextLink = sidebarLinks[currentIndex + 1];
            const nextButton = document.createElement('a');
            nextButton.href = nextLink.href;
            nextButton.className = 'page-nav-button page-nav-next';
            nextButton.innerHTML = `
                <span class="page-nav-label">Next</span>
                <span class="page-nav-title">${nextLink.textContent}</span>
            `;
            navContainer.appendChild(nextButton);
        }

        // Add to content
        content.appendChild(navContainer);
    }

    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', createPageNavigation);
    } else {
        createPageNavigation();
    }

    // Re-initialize after page navigation
    document.addEventListener('mdbook-content-loaded', createPageNavigation);
})();
