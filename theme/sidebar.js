// Collapsible sidebar categories for Go Design Patterns
(function() {
    'use strict';

    function initCollapsibleSidebar() {
        const categories = document.querySelectorAll('.sidebar .chapter li.part-title');
        
        categories.forEach((category) => {
            // Add collapse indicator
            const indicator = document.createElement('span');
            indicator.className = 'collapse-indicator';
            indicator.textContent = '−';
            category.appendChild(indicator);
            
            // Find the next chapter items (siblings) until the next part-title
            const items = [];
            let nextElement = category.parentElement.nextElementSibling;
            
            while (nextElement && !nextElement.querySelector('li.part-title')) {
                items.push(nextElement);
                nextElement = nextElement.nextElementSibling;
            }
            
            // Toggle collapse on click
            category.addEventListener('click', function() {
                const isCollapsed = category.classList.toggle('collapsed');
                indicator.textContent = isCollapsed ? '+' : '−';
                
                items.forEach(item => {
                    item.style.display = isCollapsed ? 'none' : 'block';
                });
            });
        });
    }

    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initCollapsibleSidebar);
    } else {
        initCollapsibleSidebar();
    }

    // Re-initialize after page navigation
    document.addEventListener('mdbook-content-loaded', initCollapsibleSidebar);
})();
