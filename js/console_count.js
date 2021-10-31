/**
 * Example adapted from:
 * https://developer.mozilla.org/en-US/docs/Web/API/console/count
 */

let user = "";

function greet() {
    console.count(user);
    return "hi " + user;
}

user = "bob";
greet();
user = "alice";
greet();
greet();
console.count("alice");
