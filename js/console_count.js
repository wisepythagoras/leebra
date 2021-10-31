/**
 * Example adapted from:
 * https://developer.mozilla.org/en-US/docs/Web/API/console/count
 * https://developer.mozilla.org/en-US/docs/Web/API/console/countReset
 */

let user = '';

function greet() {
    console.count(user);
    return 'hi ' + user;
}

user = 'bob';
greet();

user = 'alice';
greet();
greet();

console.countReset('bob');
console.count('alice');
