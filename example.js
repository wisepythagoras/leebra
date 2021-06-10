const what = 'world'

console.log('Hello', what, 1, true, null, navigator.userAgent);

localStorage.setItem('test', 'This is a test');
const val = localStorage.getItem('test');

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
