async function submitVote(event, form) {
    event.preventDefault();
    const formData = new FormData(form);
    const targetId = formData.get('target_id');
    const targetType = formData.get('target_type');
    
    try {
        const response = await fetch(form.action, {
            method: 'POST',
            body: formData,
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        });
        
        if (response.ok) {
            const data = await response.json();
            if (data.success) {
                const likeCount = document.getElementById(`likes-${targetType}-${targetId}`);
                const dislikeCount = document.getElementById(`dislikes-${targetType}-${targetId}`);
                if (likeCount) likeCount.textContent = data.likes;
                if (dislikeCount) dislikeCount.textContent = data.dislikes;
                
                // Optional: Add a little bump animation
                const btn = event.submitter;
                if (btn) {
                    btn.classList.add('scale-110');
                    setTimeout(() => btn.classList.remove('scale-110'), 200);
                }
            }
        }
    } catch (e) {
        console.error('Vote failed', e);
    }
}

async function submitComment(event, form) {
    event.preventDefault();
    const formData = new FormData(form);
    
    try {
        const response = await fetch(form.action, {
            method: 'POST',
            body: formData,
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        });
        
        if (response.ok) {
            const data = await response.json();
            if (data.success) {
                // For comments, a full refresh is often cleaner to ensure all attributes (Author badge, etc.) 
                // are calculated correctly by the server. But we do it instantly.
                window.location.reload();
            }
        }
    } catch (e) {
        console.error('Comment failed', e);
    }
}
