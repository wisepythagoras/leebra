const array = new Uint8Array(10);
crypto.getRandomValues(array);

const array1 = new Uint16Array(10);
crypto.getRandomValues(array1);

const array2 = new Uint32Array(10);
crypto.getRandomValues(array2);

const array3 = new Int8Array(10);
crypto.getRandomValues(array3);

const array4 = new Int16Array(10);
crypto.getRandomValues(array4);

const array5 = new Int32Array(10);
crypto.getRandomValues(array5);

const array6 = new Float32Array(10);
crypto.getRandomValues(array6);

const array7 = new Float64Array(10);
crypto.getRandomValues(array7);

console.log(array);
console.log(array1);
console.log(array2);
console.log(array3);
console.log(array4);
console.log(array5);
console.log(array6);
console.log(array7);
