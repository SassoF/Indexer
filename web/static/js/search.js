document.getElementById('searchForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const searchTerm = document.getElementById('searchInput').value.trim();
    const url = new URL(window.location.href);
    
    url.searchParams.delete('p');
    
    if (searchTerm) {
        url.searchParams.set('q', searchTerm);
    } else {
        url.searchParams.delete('q');
    }
    
    window.location.href = url.toString();
});

document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.getElementById('searchInput');
    if (searchInput) searchInput.focus();
});