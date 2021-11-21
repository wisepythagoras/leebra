const what = 123;

console.log(JSON, localStorage, WebAssembly.instantiate)
console.log('Hello', what, 1, true, null, navigator.userAgent);
console.log(document.title, window.document.title);

const add = (a, b) => a + b;

console.log(add(3, 4))
console.log(window.navigator.userAgent);
console.log(window === window.window, window === window.window.window);

// This will work if you navigate to https://developer.mozilla.org/.
const el = document.getElementById('nav-footer');

if (el) {
    console.log("Class name:", el.className);
}
