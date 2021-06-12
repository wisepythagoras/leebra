console.log("Create first record");

localStorage.setItem('test', 'This is a test');
localStorage.setItem('test', 'Override');

console.log("Get first record");

const val = localStorage.getItem('test');

console.log(`Contents of localStorage['test'] = '${val}'`)

console.log(` ${localStorage.length} items in localStorage`);

console.log("Create second record");

localStorage.setItem('test2', 'This is another test');

console.log(`2nd key is: "${localStorage.key(1)}"`)

console.log(` ${localStorage.length} items in localStorage`);

console.log("Removing records");

localStorage.removeItem('test');

console.log(` ${localStorage.length} items in localStorage`);

localStorage.clear();

console.log(` ${localStorage.length} items in localStorage`);
