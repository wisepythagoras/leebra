const label = 'myLabel';

console.time(label);

for (let i = 0; i < 50000000; i++) {}

console.timeLog(label);

for (let i = 0; i < 50000000; i++) {}

console.timeEnd(label);

// Test a label that doesn't exist.
console.timeLog('randomLabel');
