// Hamburger toggle
const burgerToggle = document.getElementById('burgerToggle');
const navLinks = document.getElementById('navLinks');
if (burgerToggle) {
    burgerToggle.addEventListener('click', () => {
        navLinks.classList.toggle('open');
    });
}

// Dark mode toggle
const themeToggle = document.getElementById('theme-toggle');
if (themeToggle) {
    themeToggle.addEventListener('click', () => {
        const isDark = document.documentElement.classList.toggle('dark');
        localStorage.setItem('theme', isDark ? 'dark' : 'light');
    });
}
