console.log('Hello, world!');

const fnWithError = (a, b) => {
    const c = a + b * 2;

    throw new Error('This is an error!');

    return c;
}

fnWithError();
