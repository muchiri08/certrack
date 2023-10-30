const currentPath = window.location.pathname;
const domainsLink = document.getElementById('domains');
const newLink = document.getElementById('new');
const accountLink = document.getElementById('account');

if (currentPath === '/home/domains') {
    domainsLink.classList.add('bg-slate-100', 'rounded-full')
} else if (currentPath === '/home/new') {
    newLink.classList.add('bg-slate-100', 'rounded-full')
} else if (currentPath === '/home/account') {
    accountLink.classList.add('bg-slate-100', 'rounded-full')
}

console.log("Hello world");
