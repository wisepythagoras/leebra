const label = 'myLabel';

console.time(label);

for (let i = 0; i < 100000000; i++) {}

console.timeEnd(label);
