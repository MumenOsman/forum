document.addEventListener('DOMContentLoaded', () => {
    initTheme();
    initVoting();
    initSearch();
    initRichEditor();
});

/**
 * Theme Toggle logic
 */
function initTheme() {
    const themeToggle = document.getElementById('theme-toggle');
    const currentTheme = localStorage.getItem('theme') || 'light';

    document.documentElement.setAttribute('data-theme', currentTheme);

    themeToggle.addEventListener('click', () => {
        const newTheme = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);

        // Add a little animation class
        themeToggle.classList.add('rotate-animation');
        setTimeout(() => themeToggle.classList.remove('rotate-animation'), 500);
    });
}

/**
 * AJAX Voting logic
 */
function initVoting() {
    const voteForms = document.querySelectorAll('form[action="/vote"]');

    voteForms.forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const formData = new FormData(form);
            const submitter = e.submitter; // The button clicked
            const voteType = submitter.value;

            // Add loading state
            submitter.classList.add('loading');

            try {
                const response = await fetch(form.action, {
                    method: 'POST',
                    body: formData,
                });

                if (response.ok) {
                    // In a real app, we'd return the new counts in JSON
                    // For now, since the backend might just redirect, 
                    // we'll simulate the update or reload if needed.
                    // If the backend is set up for AJAX, it should return JSON.

                    // Optimization: Just reload for now if it's not an API yet
                    // but show a toast first.
                    showToast('Vote recorded!', 'success');

                    // Optional: Update UI optimistically if we had the data
                    // For this demo, let's just refresh to show the result
                    window.location.reload();
                } else {
                    showToast('Failed to record vote.', 'error');
                }
            } catch (err) {
                console.error('Voting error:', err);
                showToast('Network error.', 'error');
            } finally {
                submitter.classList.remove('loading');
            }
        });
    });
}

/**
 * Client-side Search
 */
function initSearch() {
    // Add search bar dynamically to the feed if it doesn't exist
    const feedHeader = document.querySelector('.feed h2');
    if (!feedHeader) return;

    const searchContainer = document.createElement('div');
    searchContainer.style.marginBottom = '1.5rem';
    searchContainer.innerHTML = `
        <input type="text" id="post-search" placeholder="Search discussions..." 
               class="glass" style="width: 100%; padding: 0.75rem 1rem; border-radius: var(--radius); outline: none;">
    `;
    feedHeader.after(searchContainer);

    const searchInput = document.getElementById('post-search');
    const posts = document.querySelectorAll('.post-card');

    searchInput.addEventListener('input', (e) => {
        const term = e.target.value.toLowerCase();

        posts.forEach(post => {
            const title = post.querySelector('a').textContent.toLowerCase();
            const visible = title.includes(term);
            post.style.display = visible ? 'block' : 'none';

            if (visible && term !== '') {
                post.style.animation = 'fadeIn 0.3s ease';
            }
        });
    });
}

/**
 * Simple Toast Notification system
 */
function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `glass toast toast-${type}`;
    toast.style.cssText = `
        position: fixed;
        bottom: 2rem;
        right: 2rem;
        padding: 1rem 1.5rem;
        border-radius: var(--radius);
        z-index: 1000;
        animation: slideIn 0.3s ease-out;
        box-shadow: var(--shadow);
    `;

    const colors = {
        success: '#2ecc71',
        error: '#e74c3c',
        info: '#3498db'
    };

    toast.style.borderLeft = `4px solid ${colors[type]}`;
    toast.textContent = message;

    document.body.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'slideOut 0.3s ease-in forwards';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

/**
 * Basic Rich Text Editor for post creation
 */
function initRichEditor() {
    const editor = document.getElementById('rich-editor');
    if (!editor) return;

    // Create toolbar
    const toolbar = document.createElement('div');
    toolbar.className = 'glass';
    toolbar.style.padding = '0.5rem';
    toolbar.style.marginBottom = '0.5rem';
    toolbar.style.borderRadius = 'var(--radius)';
    toolbar.style.display = 'flex';
    toolbar.style.gap = '0.5rem';

    const actions = [
        { label: '<b>B</b>', cmd: 'bold' },
        { label: '<i>I</i>', cmd: 'italic' },
        { label: '<u>U</u>', cmd: 'underline' },
        { label: '• List', cmd: 'insertUnorderedList' }
    ];

    actions.forEach(act => {
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'vote-btn';
        btn.innerHTML = act.label;
        btn.onclick = () => {
            document.execCommand(act.cmd, false, null);
            editor.focus();
        };
        toolbar.appendChild(btn);
    });

    editor.before(toolbar);

    // Sync editor content to hidden input for form submission
    const form = editor.closest('form');
    if (form) {
        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'content';
        form.appendChild(input);

        form.onsubmit = () => {
            input.value = editor.innerHTML;
        };
    }
}

// Add necessary animations to CSS via JS if not already there
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn { from { transform: translateX(100%); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
    @keyframes slideOut { from { transform: translateX(0); opacity: 1; } to { transform: translateX(100%); opacity: 0; } }
    @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
    @keyframes rotate-animation { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
    .rotate-animation { animation: rotate-animation 0.5s ease; }
    .loading { opacity: 0.5; pointer-events: none; }
`;
document.head.appendChild(style);
