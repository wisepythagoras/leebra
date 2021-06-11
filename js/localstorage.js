localStorage.setItem('test', 'This is a test');
localStorage.setItem('test', 'Override');

const val = localStorage.getItem('test');

console.log(`Contents of localStorage['test'] = '${val}'`)

console.log(`${localStorage.length} items in localStorage`);

localStorage.setItem('test2', 'This is another test');

console.log(`${localStorage.length} items in localStorage`);
