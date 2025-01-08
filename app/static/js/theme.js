function toggleTheme() {
    const body = document.body;
    const icon = document.getElementById('theme-toggle-icon');

    if (body.classList.contains('light-mode')) {
        body.classList.remove('light-mode');
        body.classList.add('dark-mode');
        icon.textContent = '☀️';
        localStorage.setItem('theme', 'dark-mode');
    } else {
        body.classList.remove('dark-mode');
        body.classList.add('light-mode');
        icon.textContent = '🌙';
        localStorage.setItem('theme', 'light-mode');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
        document.body.classList.add(savedTheme);
    } else {
        document.body.classList.add('light-mode');
    }

    const icon = document.getElementById('theme-toggle-icon');
    icon.textContent = document.body.classList.contains('dark-mode') ? '☀️' : '🌙';
});