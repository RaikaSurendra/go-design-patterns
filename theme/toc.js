// Table of Contents - "On this page" component
(function() {
    'use strict';

    function createTOC() {
        const content = document.querySelector('.content main');
        if (!content) return;

        // Get all h2 and h3 headings
        const headings = content.querySelectorAll('h2, h3');
        if (headings.length === 0) return;

        // Create TOC container
        const tocContainer = document.createElement('div');
        tocContainer.className = 'toc-container';
        
        const tocTitle = document.createElement('div');
        tocTitle.className = 'toc-title';
        tocTitle.textContent = 'On this page';
        
        const tocList = document.createElement('ul');
        tocList.className = 'toc-list';

        headings.forEach((heading, index) => {
            // Add ID to heading if it doesn't have one
            if (!heading.id) {
                heading.id = `heading-${index}`;
            }

            const li = document.createElement('li');
            li.className = heading.tagName === 'H2' ? 'toc-h2' : 'toc-h3';
            
            const link = document.createElement('a');
            link.href = `#${heading.id}`;
            link.textContent = heading.textContent;
            link.className = 'toc-link';
            
            // Highlight active section on scroll
            link.addEventListener('click', (e) => {
                e.preventDefault();
                heading.scrollIntoView({ behavior: 'smooth', block: 'start' });
                window.history.pushState(null, null, `#${heading.id}`);
            });
            
            li.appendChild(link);
            tocList.appendChild(li);
        });

        tocContainer.appendChild(tocTitle);
        tocContainer.appendChild(tocList);

        // Insert TOC after the content
        const contentWrapper = document.querySelector('.content');
        if (contentWrapper) {
            contentWrapper.appendChild(tocContainer);
        }

        // Highlight current section on scroll
        let ticking = false;
        window.addEventListener('scroll', () => {
            if (!ticking) {
                window.requestAnimationFrame(() => {
                    updateActiveTOC(headings);
                    ticking = false;
                });
                ticking = true;
            }
        });
    }

    function updateActiveTOC(headings) {
        const scrollPos = window.scrollY + 100;
        let currentHeading = null;

        headings.forEach(heading => {
            if (heading.offsetTop <= scrollPos) {
                currentHeading = heading;
            }
        });

        document.querySelectorAll('.toc-link').forEach(link => {
            link.classList.remove('active');
        });

        if (currentHeading) {
            const activeLink = document.querySelector(`.toc-link[href="#${currentHeading.id}"]`);
            if (activeLink) {
                activeLink.classList.add('active');
            }
        }
    }

    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', createTOC);
    } else {
        createTOC();
    }

    // Re-initialize after page navigation
    document.addEventListener('mdbook-content-loaded', createTOC);
})();
