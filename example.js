const what = 'world'

console.log(JSON, localStorage)
console.log('Hello', what, 1, true, null, navigator.userAgent);

localStorage.setItem('test', 'This is a test');
localStorage.setItem('test', 'Override');
const val = localStorage.getItem('test');
console.log(`${localStorage.length} items in localStorage`);
localStorage.setItem('test2', 'This is another test');
console.log(`${localStorage.length} items in localStorage`);

console.log(`Local storage content: "${val}"`);

const add = (a, b) => a + b;

const fetchTest = async () => {
    const response = await fetch('https://example.com/');

    console.log(response.text());
};

const clipboardTest = async () => {
    const initText = await navigator.clipboard.readText();
    console.log("Initial Clipboard", initText);
    const res = await navigator.clipboard.writeText('hello, world!');
    const text = await navigator.clipboard.readText();
    console.log("Clipboard", text);
};

clipboardTest();
