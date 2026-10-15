const root = document.documentElement;
const themeBtn = document.querySelector('#theme-toggle');
const saved = localStorage.getItem('jorge-portfolio-theme');
if (saved === 'dark') root.dataset.theme = 'dark';
themeBtn?.addEventListener('click', () => {
  const dark = root.dataset.theme !== 'dark';
  root.dataset.theme = dark ? 'dark' : 'light';
  localStorage.setItem('jorge-portfolio-theme', dark ? 'dark' : 'light');
  themeBtn.setAttribute('aria-label', dark ? 'Switch to light theme' : 'Switch to dark theme');
});
const navToggle = document.querySelector('#nav-toggle');
const mobileNav = document.querySelector('#mobile-nav');
function closeNav() {mobileNav.hidden = true; navToggle.setAttribute('aria-expanded', 'false'); navToggle.setAttribute('aria-label', 'Open navigation');}
navToggle?.addEventListener('click', () => {const open = mobileNav.hidden; mobileNav.hidden = !open; navToggle.setAttribute('aria-expanded', String(open)); navToggle.setAttribute('aria-label', open ? 'Close navigation' : 'Open navigation');});
mobileNav?.querySelectorAll('a').forEach(a => a.addEventListener('click', closeNav));
document.addEventListener('keydown', e => {if (e.key === 'Escape') closeNav();});
const filters = document.querySelectorAll('.filter');
filters.forEach(btn => btn.addEventListener('click', () => {
  filters.forEach(b => {b.classList.toggle('active', b === btn); b.setAttribute('aria-pressed', String(b === btn));});
  const category = btn.dataset.filter;
  document.querySelectorAll('.project').forEach(project => {project.hidden = category !== 'all' && project.dataset.category !== category;});
}));
const toast = document.querySelector('#toast');
let toastTimeout;
function showToast(message) {toast.textContent = message; toast.classList.add('show'); clearTimeout(toastTimeout); toastTimeout = setTimeout(() => toast.classList.remove('show'), 2600);}
document.querySelector('#copy-email')?.addEventListener('click', async () => {
  try {await navigator.clipboard.writeText('jorgecalero4667@gmail.com'); showToast('Email address copied to clipboard.');}
  catch {showToast('Email: jorgecalero4667@gmail.com');}
});
document.querySelector('#year').textContent = new Date().getFullYear();
